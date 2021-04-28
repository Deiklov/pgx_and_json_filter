package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"sync"
	"time"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	wg := sync.WaitGroup{}
	wg.Add(1)
	conn, err := stan.Connect("cluster1", uuid.New().String(), stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}

	connConfig, err := pgx.ParseConfig("postgres://test_user:test_password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	connPgx, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(connPgx)

	sub, err := conn.QueueSubscribe("applicant", "workers", func(m *stan.Msg) {
		//ackwait время на обработку
		fmt.Printf("got task %s\n", m.Data)
		time.Sleep(1100 * time.Millisecond)
		//m.Ack()
	}, stan.AckWait(1*time.Second), stan.MaxInflight(100), stan.DeliverAllAvailable(), stan.SetManualAckMode())

	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	nc.Subscribe("applicant", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
	})
	wg.Wait()
	//закроет коннект и все подписки
	if err := nc.Drain(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("job ended")
}
