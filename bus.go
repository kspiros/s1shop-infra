package xlib

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"time"

	"github.com/nats-io/stan.go"
)

const ACKWAIT = 30 * time.Second

type IBus interface {
	Publish(subj string, data []byte) error
	QueueSubscribe(listener func(msg *stan.Msg), subj string, queue string) error
}

type stanClient struct {
	con *stan.Conn
}

func (sc *stanClient) Publish(subj string, data []byte) error {
	return (*sc.con).Publish(subj, data)

}

func (sc *stanClient) QueueSubscribe(listener func(msg *stan.Msg), subj string, queue string) error {
	durable := queue
	_, err := (*sc.con).QueueSubscribe(subj, queue, listener, stan.DeliverAllAvailable(), stan.SetManualAckMode(), stan.AckWait(ACKWAIT), stan.DurableName(durable))
	if err != nil {
		(*sc.con).Close()
		return err
	}

	log.Printf("Listening on [%s],  qgroup=[%s] durable=[%s]\n", subj, queue, queue)
	return nil
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewBus() (IBus, func(), error) {

	dsn := os.Getenv("NATS_DSN")
	if len(dsn) == 0 {
		return nil, nil, errors.New("NATS_URL must be defined")
	}

	clusterID := os.Getenv("NATS_CLUSTER_ID")
	if len(dsn) == 0 {
		return nil, nil, errors.New("NATS_CLUSTER_ID must be defined")
	}
	// Connect to NATS
	clientID, _ := randomHex(20)

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(dsn),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))

	if err != nil {
		return nil, nil, errors.New("Can't connect: Make sure a NATS Streaming Server is running")
	}
	return &stanClient{con: &sc}, func() { sc.Close() }, nil
}
