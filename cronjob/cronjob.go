package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
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

	connConfig, err := pgx.ParseConfig("postgres://test_user:test_password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	connPgx, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(connPgx)


	wg := sync.WaitGroup{}
	wg.Add(1)

	conn, err := stan.Connect("cluster1", uuid.New().String(), stan.NatsConn(nc))
	fmt.Println(conn)
	if err != nil {
		log.Fatal(err)
	}
	ackHandler := func(ackedNuid string, err error) {
		fmt.Printf("err %v\n", err)
		fmt.Printf("acked nuid %s\n", ackedNuid)
	}
	for i := 0; i < taskCnt; i++ {
		message := []byte(fmt.Sprintf("kekmda message %d", i))
		nc.Publish("applicant", []byte("some data from nats connect"))
		if _, err := conn.PublishAsync("applicant", message, ackHandler); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("kek.%d\n", i)
		time.Sleep(time.Millisecond * 1000)
	}

	wg.Wait()
}
