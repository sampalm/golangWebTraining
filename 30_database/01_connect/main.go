package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// connect to local database
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/test02?charset=utf8")
	check("Connection func", err)
	defer db.Close()
	// check connection
	err = db.Ping()
	check("Ping func", err)
	// handle serve
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	err = http.ListenAndServe(":8080", nil)
	check("Serve func", err)
}

func index(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Connection completed.")
	check("Index func", err)
}

func check(caller string, err error) {
	if err != nil {
		fmt.Printf("%s: %s\n", caller, err.Error())
	}
}
