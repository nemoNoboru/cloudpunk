package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/yuin/gopher-lua"
)

func ExtractWildcardValues(pattern, filled string) []string {
	// Split both pattern and filled strings into segments
	patternSegments := strings.Split(pattern, ".")
	filledSegments := strings.Split(filled, ".")

	// Check if the number of segments matches
	if len(filledSegments) != len(patternSegments) {
		return nil
	}

	// Prepare a slice to hold the wildcard values
	wildcardValues := make([]string, 0, len(patternSegments))

	// Iterate through the segments
	for i := range patternSegments {
		switch patternSegments[i] {
		case "*":
			// If the segment is a wildcard, add the corresponding filled segment
			wildcardValues = append(wildcardValues, filledSegments[i])
		default:
			// If the segment is not a wildcard, check if it matches the filled segment
			if patternSegments[i] != filledSegments[i] {
				return nil
			}
		}
	}

	return wildcardValues
}

// Connect to main nats server
var nc, err = nats.Connect("ws://punk:cloud@connect.cloudpunk.org")

func StorageGet(label string) []byte {
	storePath := fmt.Sprintf("cloudpunk.storage.get.%s", label)

	result, err := nc.Request(storePath, nil, 30*time.Second)
	if err != nil {
		return nil
	}

	return result.Data
}

func handleAPI(w http.ResponseWriter, req *http.Request) {
	var fnpath = fmt.Sprintf("cloudpunk.serverless.%s", strings.Split(req.URL.Path, "/")[2])
	fmt.Println(fnpath)
	result, err := nc.Request(fnpath, []byte("hello"), 30*time.Second)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	fmt.Fprint(w, string(result.Data))
}

func handleStatic(w http.ResponseWriter, req *http.Request) {
	label := strings.Split(req.URL.Path, "/")[2]

	result := StorageGet(label)

	w.Write(result)
}

// Roadmap:
// Adding and running functions dynamically is DONE AND WORKING.
// Design Decisions:
// all information should be in nats. HTML, Images, Code
// Next:
// - [] Implement a bridge nats-lua to be able to call remote functions
// - [x] Implement a HTTP endpoint to call functions
// - - [] adapt common HTTP params into a map to be sent to the lua function
// - [] Research: Check how we could do templating here. Go templates vs client directly

// node bootstraping:
// Plan: we have three diferent type of data.
// 1. DB: mutable, runtime data that MUST be preserved
// 2. BLob: inmutable, on-disk data
// 3. Functions: runnable on-disk data.

// the plan is to run cloudpunk from a folder that has three subfolders:
// /functions - lua files containing functions. They will callable from nats or http://namespace.cloudpunk.org/api
// /static - various files containing static data. They will be callable also from nats or http://namespace.cloudpunk.org/*

// the idea is to have a "LOAD" message with certains metadata
// - data ([]byte the actual data we want to load)
// - handler (string, luavm, just a http response etc. )
// - topic
// - optional fields...

// then, if the handler accepts the LOAD
// it suscribes to the topic in order to RUN the HANDLER

// handlers could be:
// DB
// LUA
// WASM
// FILE - for serving static files

func serve(port string) {

	L := lua.NewState()
	defer L.Close()

	// NOTE: this lives in memory right now,
	// maybe is a good idea to use badger to move this to disc.
	store := make(map[string]([]byte))

	storageLoad := "cloudpunk.storage.load.*"
	nc.Subscribe(storageLoad, func(msg *nats.Msg) {
		path := ExtractWildcardValues(storageLoad, msg.Subject)[0]
		store[path] = msg.Data
	})

	// TODO: don't really listen to the whole wildcard.
	// Listen only to the paths that we have in store
	storageServe := "cloudpunk.storage.get.*"
	nc.Subscribe(storageServe, func(msg *nats.Msg) {
		path := ExtractWildcardValues(storageServe, msg.Subject)[0]
		msg.Respond(store[path])
	})

	// TODO: It would be interesting to add another handler
	// to do a sort of a php thingy with go templates + lua code

	luaServe := "cloudpunk.serverless.*"
	// this code here is a example of a LUA handler
	nc.Subscribe(luaServe, func(msg *nats.Msg) {
		path := ExtractWildcardValues(luaServe, msg.Subject)[0]

		fn_source := StorageGet(path)

		if err := L.DoString(string(fn_source)); err != nil {
			msg.Respond([]byte(err.Error()))
		}

		err := L.CallByParam(lua.P{
			Fn:      L.GetGlobal("serverless"),
			NRet:    1,
			Protect: true,
		}, lua.LString(msg.Data))

		if err != nil {
			msg.Respond([]byte(err.Error()))
		}

		result := L.Get(-1)
		msg.Respond([]byte(result.String()))
	})

	http.HandleFunc("/api/*", handleAPI)
	// http.HandleFunc("/*", handleStatic)

	http.ListenAndServe(port, nil)

}

// ServeCommand handles the 'serve' command
func ServeCommand() {
	// Define flags specific to the 'serve' command
	port := flag.String("port", ":1332", "Port to run the server on")

	// Parse the flags
	flag.Parse()

	// Your serve logic here
	serve(*port)
	// For example: start a server or perform actions related to serving
}

// NOTE: Would be nice to be able to load a entire folder
// more convinient etc. The only thing to do is to be able
// to do more neat routes as /pages/users?id=123
// that could be a label like pages-users

// UploadCommand handles the 'upload' command
func UploadCommand() {
	// Define flags specific to the 'upload' command
	// filePath := flag.String("file", "", "Path to the file to upload")
	// destination := flag.String("label", "", "Upload destination")

	filePath := os.Args[2]
	destination := os.Args[3]

	// Parse the flags
	flag.Parse()

	fmt.Println(filePath)
	fmt.Println(destination)

	// Your upload logic here
	if filePath == "" || destination == "" {
		fmt.Println("Error: -file and -dest flags are required for the upload command")
		return
	}

	source, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	storagePath := fmt.Sprintf("cloudpunk.storage.load.%s", destination)

	nc.Publish(storagePath, source)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("uploaded file at %s \n", storagePath)
}

func main() {
	// check status of nats
	if err != nil {
		log.Fatal(err.Error())
	}

	defer nc.Drain()
	defer nc.Close()

	if len(os.Args) < 2 {
		fmt.Println("Expected 'serve' or 'upload' command")
		os.Exit(1)
	}

	// Switch based on the second argument (first argument is the program name)
	switch os.Args[1] {
	case "serve":
		// Re-parse the arguments to isolate flags from other commands
		flag.CommandLine = flag.NewFlagSet("serve", flag.ExitOnError)
		ServeCommand()
	case "upload":
		// Re-parse the arguments to isolate flags from other commands
		flag.CommandLine = flag.NewFlagSet("upload", flag.ExitOnError)
		UploadCommand()
	default:
		fmt.Println("Unknown command. Available commands are: 'serve' and 'upload'")
		os.Exit(1)
	}
}

