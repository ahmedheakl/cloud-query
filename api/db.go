package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq" // import declerations/initializations only
)

const conn = "host=terraform-20221220092610558600000001.clgu6fwjxdcc.me-south-1.rds.amazonaws.com port=5432 user=postgres password=postgres database=testdb"

var db *sql.DB

type Record struct {
	Id   int
	Name string
	Test string
}

func columnField(name string, rec *Record) interface{} {
	switch name {
	case "id":
		return &rec.Id
	case "name":
		return &rec.Name
	case "test":
		return &rec.Test
	default:
		panic(fmt.Sprintf("Unknown name %s", name))
	}
}

func getOrCreate() (*sql.DB, error) {
	var err error

	if db == nil {
		db, err = sql.Open("postgres", conn)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func executeQuery(query string) ([]byte, error) {
	var records []Record

	db, err := getOrCreate()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	columns, _ := rows.Columns()

	for rows.Next() {
		record := Record{}
		cols := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			cols[i] = columnField(columns[i], &record)
		}
		if err := rows.Scan(cols...); err != nil {
			println(err.Error())
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		println(err.Error())
	}

	return json.Marshal(records)
}
