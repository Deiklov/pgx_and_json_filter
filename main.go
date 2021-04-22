package main

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	"io/ioutil"
	"log"
)

type req struct {
	OCONTROL struct {
		ALIAS string `json:"ALIAS" faker:"word"`
	} `json:"OCONTROL"`
	NewProductLine string `json:"NewProductLine" faker:"word"`
	External       string `json:"External" faker:"word"`
}

type resp struct {
	OCONTROL struct {
		ALIAS string `json:"ALIAS" faker:"word"`
	} `json:"OCONTROL"`
	NewProductLine string `json:"NewProductLine" faker:"word"`
	External       string `json:"External" faker:"word"`
	Result         string `json:"Result" faker:"word"`
	City           string `json:"City" `
}

func main() {
	connConfig, err := pgx.ParseConfig("postgres://test_user:test_password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal(err)
	}
	nc, _ := nats.Connect(nats.DefaultURL)
	nc.Publish("foo", []byte("Hello World"))

	var transaction_resp string

	data, err := ioutil.ReadFile("fat_json.json")

	if err != nil {
		log.Fatal(err)
	}
	err = conn.QueryRow(context.Background(), `insert into request_response values (default, default, $1, $2) returning response`,
		data, data).Scan(&transaction_resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(transaction_resp)

}
