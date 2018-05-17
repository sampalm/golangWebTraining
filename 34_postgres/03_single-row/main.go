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

var db *sql.DB

func main() {
	q := os.Args[1]
	qConn()
	qSelectRow(q)
}

func qConn() {
	var err error
	db, err = sql.Open("postgres", "postgres://dev:devpass@localhost/company?sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func qSelectRow(query string) {
	var emp employes
	row := db.QueryRow("SELECT * FROM employess WHERE id = $1", query)
	err := row.Scan(&emp.ID, &emp.Name, &emp.Score, &emp.Salary)
	switch {
	case err == sql.ErrNoRows:
		log.Fatalln("NO ROWS")
		return
	case err != nil:
		log.Fatalln(err)
		return
	}
	log.Printf("\n-----------------------------------\nRECORD: %d, %s, %d, $%.2f\n-----------------------------------\n", emp.ID, emp.Name, emp.Score, emp.Salary)

}
