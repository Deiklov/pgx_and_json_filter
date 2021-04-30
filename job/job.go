package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"sync"
	"time"
)

type PrimaryKey struct {
	GuidTransaction string `json:"guid_transaction"`
	GuidStrategy    string `json:"guid_strategy"`
}

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)

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

	var wg sync.WaitGroup

	wg.Add(1)
	//в реале ставить таймаут на 10 минут
	time.AfterFunc(3*time.Second, wg.Done)

	var pKey PrimaryKey

	sub, err := conn.QueueSubscribe("applicant", "workers", func(m *stan.Msg) {
		wg.Add(1)
		err := json.Unmarshal(m.Data, &pKey)
		if err != nil {
			log.Fatal(err)
		}
		go messageHandler(pKey, connPgx)
		m.Ack()
		time.AfterFunc(3*time.Second, wg.Done)

	}, stan.AckWait(100*time.Second), stan.SetManualAckMode())
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	defer sub.Unsubscribe()

	if err := nc.Drain(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("job ended")
}

func messageHandler(pKey PrimaryKey, conn *pgx.Conn) {
	//идет по grpc в api ОКР
	fmt.Printf("analyze task, %v\n", pKey)
	_, err := conn.Exec(context.Background(), "update jobs_result set worker_status=true where guid_transaction=$1 and guid_strategy =$2",
		pKey.GuidTransaction, pKey.GuidStrategy)
	if err != nil {
		log.Fatal(err)
	}

}
