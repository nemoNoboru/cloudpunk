package main

import (
	"cloudpunk/cloud"
	"cloudpunk/handlers"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func serve(port string) {
	http.HandleFunc("/api/*", handlers.HandleAPI)
	http.HandleFunc("/*", handlers.HandleStatic)

	cloud.StartLuaServerless()
	cloud.StartStorageServer()

	http.ListenAndServe(port, nil)
}

// ServeCommand handles the 'serve' command
func ServeCommand() {
	// Define flags specific to the 'serve' command
	port := flag.String("port", ":1332", "Port to run the server on")

	// Parse the flags
	flag.Parse()

	serve(*port)
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

	// Your upload logic here
	if filePath == "" || destination == "" {
		fmt.Println("Error: -file and -dest flags are required for the upload command")
		return
	}

	source, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	cloud.StorageLoad(destination, source)

	log.Printf("uploaded file at %s \n", destination)
}

func main() {
	defer cloud.NatsConn.Drain()
	defer cloud.NatsConn.Close()

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
