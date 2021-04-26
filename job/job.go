package main

import (
	"bytes"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Подписываемся на все каналалы вида kek.* ждем прихода сообщения о закрытии
	if _, err := nc.Subscribe("topic.final", func(msg *nats.Msg) {
		if bytes.Compare(msg.Data, []byte("{type:close}")) == 0 {
			wg.Done()
		}
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := nc.QueueSubscribe("kek.>", "workers", func(m *nats.Msg) {
		fmt.Printf("reply data %s\n", m.Reply)
		if err := m.Respond([]byte("reply from job")); err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(m.Data))
	}); err != nil {
		log.Fatal(err)

	}

	wg.Wait()
	//закроет коннект и все подписки
	if err := nc.Drain(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("job ended")
}
