package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func init() {
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/test02?charset=utf8")
	check("Connection func", err)
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/users/", users)
	router.GET("/user/:name", users)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Connection completed. Welcome!")
}

func users(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var query string
	query = `SELECT * FROM amigos` // Default query
	if pName := p.ByName("name"); pName != "" {
		query = fmt.Sprintf("SELECT * FROM amigos WHERE aNomes=\"%s\"", pName)
	}

	// user data
	var (
		id   int64
		name string
	)
	rw, err := db.Query(query)
	check("querySelect func", err)

	for rw.Next() {
		err = rw.Scan(&id, &name)
		check("querySelect Next func", err)
		fmt.Fprintf(w, "ID: %d - Name: %s\n", id, name)
	}
}

func check(from string, err error) {
	if err != nil {
		log.Printf("%s: %s", from, err.Error())
	}
}
