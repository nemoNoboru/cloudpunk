package cloud

import (
	"cloudpunk/utils"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/yuin/gopher-lua"
)

func StartLuaServerless() {

	// TODO: It would be interesting to add another handler
	// to do a sort of a php thingy with go templates + lua code

	luaServe := "cloudpunk.serverless.*"
	// this code here is a example of a LUA handler
	NatsConn.Subscribe(luaServe, func(msg *nats.Msg) {
		var L = lua.NewState()
		defer L.Close()

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
			Fn:      L.GetGlobal("serverless"),
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

func LuaRun(label string, data []byte) (string, error) {
	var fnpath = fmt.Sprintf("cloudpunk.serverless.%s", label)
	fmt.Println(fnpath)
	result, err := NatsConn.Request(fnpath, data, 30*time.Second)
	if err != nil {
		return "", err
	}

	return string(result.Data), nil
}
