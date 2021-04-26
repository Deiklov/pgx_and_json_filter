package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

type Topic struct {
	Type    string `json:"type"`
	TopicID int    `json:"topic_id"`
}

const taskCnt = 40

func main() {
	opts := nats.GetDefaultOptions()
	fmt.Printf("hello, job started\n")
	opts.Timeout = 5 * time.Second
	opts.ReconnectWait = 10 * time.Second
	opts.Verbose = true
	nc, _ := opts.Connect()
	defer nc.Close()
	uniqueReplyTo := nats.NewInbox()

	var mCnt int
	wg := sync.WaitGroup{}
	wg.Add(1)

	if _, err := nc.Subscribe(uniqueReplyTo, func(m *nats.Msg) {
		//получили ответ от воркера
		mCnt += 1
		if mCnt == taskCnt {
			if err := nc.Publish("topic.final", []byte("{type:close}")); err != nil {
				log.Fatal(err)
			}
			//закрывем nats соединение
			if err := nc.Drain(); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}
	}); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < taskCnt; i++ {
		message := []byte(fmt.Sprintf("kekmda message %d", i))
		if err := nc.PublishRequest(fmt.Sprintf("kek.%d", i), uniqueReplyTo, message); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("kek.%d\n", i)
		time.Sleep(time.Millisecond * 100)
	}

	wg.Wait()
}
