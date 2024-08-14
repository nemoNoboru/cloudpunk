package cloud

import (
	"cloudpunk/utils"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/yuin/gopher-lua"
	"time"
)

func StartLuaServerless() {
	L := lua.NewState()
	defer L.Close()

	// TODO: It would be interesting to add another handler
	// to do a sort of a php thingy with go templates + lua code

	luaServe := "cloudpunk.serverless.*"
	// this code here is a example of a LUA handler
	NatsConn.Subscribe(luaServe, func(msg *nats.Msg) {
		path := utils.ExtractWildcardValues(luaServe, msg.Subject)[0]

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
}

func LuaRun(label string) (string, error) {
	var fnpath = fmt.Sprintf("cloudpunk.serverless.%s", label)
	fmt.Println(fnpath)
	result, err := NatsConn.Request(fnpath, []byte("hello"), 30*time.Second)
	if err != nil {
		return "", err
	}

	return string(result.Data), nil
}
