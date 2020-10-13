// NATS server subscriber 
package edgepipeline

import (
	"log"
	"runtime"
	"time"
	"github.com/nats-io/nats.go"
)



func subscribe(natsEndpoint, topic string, imChan chan<- []byte) {

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
	opts = setupConnOptions(opts)

	// Connect to NATS
	nc, err := nats.Connect(natsEndpoint, opts...)
	if err != nil {
		log.Fatal(err)
	}

	i := 0

	nc.Subscribe(topic, func(msg *nats.Msg) {
		i += 1
		printMsg(msg, i)
		im := make([]byte, len(msg.Data))
		copy(im, msg.Data)
		imChan <- im

	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", topic)
	runtime.Goexit()
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]", i, m.Subject)
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
