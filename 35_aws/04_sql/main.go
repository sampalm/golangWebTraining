package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	// connect to local database
	db, err = sql.Open("mysql", "awsuser:mypassword@tcp(mydbinstance.aws:3306)/db_test?charset=utf8")
	check("Connection func", err)
	defer db.Close()
	// check connection
	err = db.Ping()
	check("Ping func", err)
	// handle serve
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/ping/", ping)
	http.HandleFunc("/instance/", instance)
	http.HandleFunc("/amigos/", amigos)
	err = http.ListenAndServe(":80", nil)
	check("Serve func", err)
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from AWS...")
}

func ping(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "OK")
}

func instance(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, getInstance())
}

func amigos(w http.ResponseWriter, r *http.Request) {
	var query string
	query = `SELECT * FROM amigos` // Default query

	// user data
	var (
		id   int64
		name string
	)
	ins := getInstance()
	rw, err := db.Query(query)
	check("querySelect func", err)

	for rw.Next() {
		err = rw.Scan(&id, &name)
		check("querySelect Next func", err)
		fmt.Fprintf(w, "ID: %d - Name: %s\nInstance: %s\n", id, name, ins)
	}
}

func getInstance() string {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	bs := make([]byte, resp.ContentLength)
	resp.Body.Read(bs)
	resp.Body.Close()
	return string(bs)
}

func check(caller string, err error) {
	if err != nil {
		fmt.Printf("%s: %s\n", caller, err.Error())
	}
}
