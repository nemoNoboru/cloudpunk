package cloud

import (
	"cloudpunk/utils"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/yuin/gopher-lua"
)

func getFile(L *lua.LState) int {
	label := L.ToString(1)
	result := StorageGet(label)
	L.Push(lua.LString(result))
	return 1
}

func uploadFile(L *lua.LState) int {
	label := L.ToString(1)
	value := L.ToString(2)

	err := StorageLoad(label, []byte(value))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LTrue)
	return 1
}

func luaServer(path string, funcall string) {

	// TODO: It would be interesting to add another handler
	// to do a sort of a php thingy with go templates + lua code

	luaServe := path
	// this code here is a example of a LUA handler
	NatsConn.Subscribe(luaServe, func(msg *nats.Msg) {
		var L = lua.NewState()
		defer L.Close()

		// inject the cloudpunk object into lua
		var storageLuaNamespace = L.NewTable()
		L.SetTable(storageLuaNamespace, lua.LString("get"), L.NewFunction(getFile))
		L.SetTable(storageLuaNamespace, lua.LString("set"), L.NewFunction(uploadFile))

		var cloudpunkNamespace = L.NewTable()
		L.SetTable(cloudpunkNamespace, lua.LString("storage"), storageLuaNamespace)

		L.SetGlobal("cloudpunk", cloudpunkNamespace)

		path := utils.ExtractWildcardValues(luaServe, msg.Subject)[0]

		fn_source := string(StorageGet(path))

		if fn_source == "" {
			msg.Respond([]byte(""))
			return
		}

		if err := L.DoString(string(fn_source)); err != nil {
			msg.Respond([]byte(err.Error()))
		}

		err := L.CallByParam(lua.P{
			Fn:      L.GetGlobal(funcall),
			NRet:    1,
			Protect: true,
		}, lua.LString(msg.Data))

		if err != nil {
			log.Println(err.Error())
			msg.Respond([]byte(err.Error()))
		}

		result := L.Get(-1)
		msg.Respond([]byte(result.String()))
	})

}

func StartLuaServerless() {
	luaServer("cloudpunk.serverless.*", "api")
	luaServer("cloudpunk.pages.*", "render")
}

func LuaRun(path string, label string, data []byte) (string, error) {
	var fnpath = fmt.Sprintf(path, label)
	fmt.Println(fnpath)
	result, err := NatsConn.Request(fnpath, data, 30*time.Second)
	if err != nil {
		return "", err
	}

	return string(result.Data), nil
}
