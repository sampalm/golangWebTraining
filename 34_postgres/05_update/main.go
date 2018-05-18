package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type employes struct {
	ID     int64
	Name   string
	Score  int
	Salary float32
}

var emps []employes
var db *sql.DB

func main() {
	qConn()
	qUpdated()
}

func qConn() {
	var err error
	db, err = sql.Open("postgres", "postgres://dev:devpass@localhost/company?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
}

func qUpdated() {
	_, err := db.Exec("UPDATE employess SET name=$2, score=$3, salary=$4 WHERE id=$1;", os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Data updated!")
}
