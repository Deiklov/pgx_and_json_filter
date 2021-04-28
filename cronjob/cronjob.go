package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

const taskCnt = 20

func main() {
	opts := nats.GetDefaultOptions()
	fmt.Printf("hello, job started\n")
	opts.Timeout = 5 * time.Second
	opts.ReconnectWait = 10 * time.Second
	opts.Verbose = true
	nc, err := opts.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	uniqueReplyTo := nats.NewInbox()

	var mCnt int
	wg := sync.WaitGroup{}
	wg.Add(1)

	//conn, err := stan.Connect("cluster1", "clientID1", stan.NatsConn(nc))
	//fmt.Println(conn)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//stan.NatsConn(nc)
	if _, err := nc.Subscribe(uniqueReplyTo, func(m *nats.Msg) {
		fmt.Println(string(m.Data))
		//получили ответ от воркера
		//todo может быть гонка(нету тк таски идут последовательно)
		//todo ответ если не смогли обработать сообщение
		//todo проверить когда топик удаляется
		//todo  добавить запись в postgresql и проверить
		//todo схема плоха тем что будет бесконечный цикл при ошибке
		//todo таски вырубаются по таймаутам, крон берет из базы только id
		mCnt += 1
		if mCnt == taskCnt {
			if err := nc.Publish("notifications", []byte("{type:close}")); err != nil {
				log.Fatal(err)
			}
			//закрывем nats соединение
			if err := nc.Drain(); err != nil {
				log.Fatal(err)
			}
			fmt.Println("send message to close topic")

			wg.Done()
		}
	}); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < taskCnt; i++ {
		message := []byte(fmt.Sprintf("kekmda message %d", i))
		if err := nc.PublishRequest("applicant", uniqueReplyTo, message); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("kek.%d\n", i)
		time.Sleep(time.Millisecond * 1000)
	}

	wg.Wait()
}
