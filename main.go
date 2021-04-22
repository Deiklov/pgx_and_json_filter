package main

import (
	"context"
	"encoding/json"
	"github.com/bxcodec/faker"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4"
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
	//var guid string
	//rows, err := conn.Query(context.Background(), `select guid_transaction
	//from request_response
	//where (request -> 'Applicant') @> '{"Cur_Flt": ["78"]}';`)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for rows.Next() {
	//	err = rows.Scan(&guid)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(guid)
	//}
	//defer rows.Close()

	var reqv req
	var resp resp
	var transaction_d string
	for i := 1; i < 100000; i++ {
		faker.FakeData(&reqv)
		faker.FakeData(&resp)
		reqvJson, _ := json.Marshal(reqv)
		respJson, _ := json.Marshal(resp)
		err := conn.QueryRow(context.Background(), `insert into request_response values (default, default, default, $1, $2) returning guid_transaction`,
			reqvJson, respJson).Scan(&transaction_d)
		if err != nil {
			log.Fatal(err)
		}
	}
	//fmt.Println(conn)
}
