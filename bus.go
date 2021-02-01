package xlib

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nats-io/stan.go"
)

type IBus interface {
}

type natsClient struct {
	con *stan.Conn
}

func (nc *natsClient) Publish(subj string, data []byte) error {
	return nil
}

func (nc *natsClient) QueueSubscribe(subj string, queue string) {

}

func NewBus() (IBus, func(), error) {

	dsn := os.Getenv("NATS_DSN")
	if len(dsn) == 0 {
		return nil, nil, errors.New("NATS_URL must be defined")
	}

	clientID := os.Getenv("NATS_CLIENT_ID")
	if len(dsn) == 0 {
		return nil, nil, errors.New("NATS_CLIENT_ID must be defined")
	}

	clusterID := os.Getenv("NATS_CLUSTER_ID")
	if len(dsn) == 0 {
		return nil, nil, errors.New("NATS_CLUSTER_ID must be defined")
	}
	// Connect to NATS
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(dsn))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, dsn)
	}
	defer sc.Close()

	subj := "test subject"
	msg := []byte("test message")

	ch := make(chan bool)
	var glock sync.Mutex
	var guid string
	acb := func(lguid string, err error) {
		glock.Lock()
		log.Printf("Received ACK for guid %s\n", lguid)
		defer glock.Unlock()
		if err != nil {
			log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
		}
		if lguid != guid {
			log.Fatalf("Expected a matching guid in ack callback, got %s vs %s\n", lguid, guid)
		}
		ch <- true
	}
	async := false
	if !async {
		err = sc.Publish(subj, msg)
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	} else {
		glock.Lock()
		guid, err = sc.PublishAsync(subj, msg, acb)
		if err != nil {
			log.Fatalf("Error during async publish: %v\n", err)
		}
		glock.Unlock()
		if guid == "" {
			log.Fatal("Expected non-empty guid to be returned.")
		}
		log.Printf("Published [%s] : '%s' [guid: %s]\n", subj, msg, guid)

		select {
		case <-ch:
			break
		case <-time.After(5 * time.Second):
			log.Fatal("timeout")
		}

	}
	return &natsClient{con: &sc}, func() { sc.Close() }, nil
}
