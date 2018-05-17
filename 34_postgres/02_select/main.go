package main

import (
	"database/sql"
	"fmt"
	"log"

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
	qSelect()
	fmt.Printf("FOUND %d RECORDS!\n", len(emps))
}

func qConn() {
	var err error
	db, err = sql.Open("postgres", "postgres://dev:devpass@localhost/company?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
}

func qSelect() {
	rows, err := db.Query("SELECT * FROM employess")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var emp employes
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Score, &emp.Salary); err != nil {
			log.Println(err)
		}
		emps = append(emps, emp)
	}
}
