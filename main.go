package main

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4"
	"log"
)

func main() {
	connConfig, err := pgx.ParseConfig("postgres://test_user:test_password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal(err)
	}
	var guid string
	rows, err := conn.Query(context.Background(), `select guid_transaction
	from request_response
	where (request -> 'Applicant') @> '{"Cur_Flt": ["78"]}';`)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&guid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(guid)
	}
	defer rows.Close()
	//fmt.Println(conn)
}
