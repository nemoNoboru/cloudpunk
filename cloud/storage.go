package cloud

import (
	"cloudpunk/utils"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

// Connect to main nats server
var NatsConn, err = nats.Connect("ws://punk:cloud@connect.cloudpunk.org")

// NOTE: this lives in memory right now,
// maybe is a good idea to use badger to move this to disc.
var store = make(map[string]([]byte))

func StorageGet(label string) []byte {
	storePath := fmt.Sprintf("cloudpunk.storage.get.%s", label)

	result, err := NatsConn.Request(storePath, nil, 30*time.Second)
	if err != nil {
		return nil
	}

	return result.Data
}

func StorageLoad(label string, data []byte) error {
	storePath := fmt.Sprintf("cloudpunk.storage.load.%s", label)
	return NatsConn.Publish(storePath, data)
}

func StartStorageServer() {
	storageLoad := "cloudpunk.storage.load.*"
	NatsConn.Subscribe(storageLoad, func(msg *nats.Msg) {
		path := utils.ExtractWildcardValues(storageLoad, msg.Subject)[0]
		store[path] = msg.Data

		// start serving this topic, now that we have it.
		storageServe := fmt.Sprintf("cloudpunk.storage.get.%s", path)
		NatsConn.Subscribe(storageServe, func(msg *nats.Msg) {
			path := utils.ExtractWildcardValues("cloudpunk.storage.get.*", msg.Subject)[0]
			msg.Respond(store[path])
		})
	})
}
