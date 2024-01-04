package pubsub

import (
	"os"

	nats "github.com/nats-io/nats.go"
)

var NATS *nats.Conn

func ConnectPubSub() (*nats.Conn, error) {
	natsHost := os.Getenv("NATS_HOST")
	natsPort := os.Getenv("NATS_PORT")
	// Connect to a server
	nc, err := nats.Connect("nats://" + natsHost + ":" + natsPort)
	NATS = nc
	// if err != nil {
	// 	fmt.Println("Nats error :: ", err.Error())
	// } else {
	//
	// 	// Simple Async Subscriber
	// 	nc.Subscribe("foo", func(m *nats.Msg) {
	// 		fmt.Printf("Received a message: %s\n", string(m.Data))
	// 	})
	// 	go func() {
	// 		timer := time.After(5 * time.Second)
	// 		<-timer
	// 		// Simple Publisher
	// 		nc.Publish("foo", []byte("Hello World"))
	// 	}()
	// }
	return nc, err
}
