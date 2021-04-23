package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Topic struct {
	Type    string `json:"type"`
	TopicID int    `json:"topic_id"`
}

func main() {
	opts := nats.GetDefaultOptions()
	fmt.Printf("hello, job started")
	opts.Timeout = 5 * time.Second
	opts.ReconnectWait = 10 * time.Second
	opts.Verbose = true
	nc, _ := opts.Connect()
	defer nc.Close()

	err := nc.Publish("kek", []byte("kekmda"))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		err := nc.Publish("kek", []byte(fmt.Sprintf("kekmda message %d", i)))
		time.Sleep(time.Second * 1)
		if err != nil {
			log.Fatal(err)
		}
	}
}
