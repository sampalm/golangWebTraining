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
	qDelete()
}

func qConn() {
	var err error
	db, err = sql.Open("postgres", "postgres://dev:devpass@localhost/company?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
}

func qDelete() {
	_, err := db.Exec("DELETE FROM employess WHERE id=$1;", os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Data deleted!")
}
