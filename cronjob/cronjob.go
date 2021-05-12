package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

const taskCnt = 10

type PrimaryKey struct {
	GuidTransaction pgtype.UUID `json:"guid_transaction"`
	GuidStrategy    pgtype.UUID `json:"guid_strategy"`
}

func main() {
	fmt.Printf("Hello, cronjob started\n")

	nc, err := nats.Connect(nats.DefaultURL)

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

	conn, err := stan.Connect("cluster1", uuid.New().String(), stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}

	var transact pgtype.UUID
	var strategy pgtype.UUID
	primaryKeys := make([]PrimaryKey, taskCnt)
	//вместо этого будет запрос в api ОКР
	rows, err := connPgx.Query(context.Background(), "select guid_transaction, guid_strategy from request_response order by random() limit $1", taskCnt)
	defer rows.Close()
	if rows == nil || err != nil {
		log.Fatal(err)
	}

	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&transact, &strategy)
		primaryKeys[i].GuidTransaction = transact
		primaryKeys[i].GuidStrategy = strategy
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i < taskCnt; i++ {
		message, err := json.Marshal(primaryKeys[i])
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(500 * time.Millisecond)
		// можно публиковать синхронно, тк публикация намного быстрее чем обработка
		fmt.Println(string(message))
		if err := conn.Publish("applicant", message); err != nil {
			log.Fatal(err)
		}

	}
	fmt.Println("CronJob work ended")
}
