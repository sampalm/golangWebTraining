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
	qInsert()
}

func qConn() {
	var err error
	db, err = sql.Open("postgres", "postgres://dev:devpass@localhost/company?sslmode=disable")
	if err != nil {
		log.SetPrefix("Log qConn output: ")
		log.SetFlags(0)
		log.Fatalln(err)
	}
}

func qInsert() {
	res, err := db.Exec("INSERT INTO employess (name, score, salary) VALUES ($1, $2, $3)", os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		log.SetPrefix("Log qInsert output: ")
		log.SetFlags(0)
		log.Fatalln(err)
	}
	if n, _ := res.RowsAffected(); n != 0 {
		log.Println("Employee inserted into database!")
	}
}
