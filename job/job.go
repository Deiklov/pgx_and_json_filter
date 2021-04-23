package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	ch := make(chan *nats.Msg, 64)
	ch2 := make(chan *nats.Msg, 64)
	fmt.Printf("hello, job started")
	_, err := nc.ChanSubscribe("kek", ch)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		_, err = nc.ChanSubscribe("kek", ch2)
		if err != nil {
			log.Fatal(err)
		}
		msg := <-ch
		fmt.Printf("FROM goroutine %s\n", string(msg.Data))
	}()
	msg := <-ch
	fmt.Println(string(msg.Data))
	msg = <-ch
	fmt.Println(string(msg.Data))

}
