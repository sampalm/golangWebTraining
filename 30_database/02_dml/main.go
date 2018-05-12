package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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
	router.GET("/post-user/:name", postUser)
	router.GET("/update-user/:id", updateUser)
	router.GET("/delete-user/:id", deleteUser)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Connection completed. Welcome!")
}

// SELECT FROM DATABASE
func users(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var query string
	query = `SELECT * FROM amigos` // Default query
	if pName := p.ByName("name"); pName != "" {
		query = fmt.Sprintf("SELECT * FROM amigos WHERE aNomes=\"%s\"", pName)
	}

	// user data
	var (
		id     int64
		name   string
		cidade string
	)
	rw, err := db.Query(query)
	check("querySelect func", err)

	for rw.Next() {
		err = rw.Scan(&id, &name, &cidade)
		check("querySelect Next func", err)
		fmt.Fprintf(w, "ID: %d - Name: %s - Cidade: %s\n", id, name, cidade)
	}
}

// INSERT INTO DATABASE
func postUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t := time.Now().Unix()
	var name string
	if name = p.ByName("name"); name == "" {
		http.Error(w, "Undefined name, must define one", 500)
		return
	}
	query := fmt.Sprintf("INSERT INTO amigos (aID, aNomes) VALUES (%d,\"%s\");", t, name)
	smt, err := db.Prepare(query)
	check("postUser Query", err)

	res, err := smt.Exec()
	check("postUser Execute", err)

	n, err := res.RowsAffected()
	check("postUser CountRows", err)
	fmt.Fprintln(w, "Data inserted successfully, affected rows: ", n)
}

// UPDATE INTO DATABASE
func updateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var id, cidade string
	if id = p.ByName("id"); id == "" {
		http.Error(w, "Undefined id, must define one", 500)
		return
	}
	values := r.URL.Query()
	if cidade = values.Get("cidade"); cidade == "" {
		http.Error(w, "Undefined cidade, must define one", 500)
		return
	}
	q, err := db.Prepare(fmt.Sprintf("UPDATE amigos SET aCidade=\"%s\" WHERE aID = %s;", cidade, id))
	check("postUser PrepareQuery", err)
	res, err := q.Exec()
	check("postUser Execute", err)
	count, err := res.RowsAffected()
	check("postUser CountRows", err)
	fmt.Fprintln(w, "Data updated successfully, affected rows: ", count)
}

// DELETE FROM DATABASE
func deleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var id string
	if id = p.ByName("id"); id == "" {
		http.Error(w, "Undefined id, must define one", 500)
		return
	}
	smt, err := db.Prepare(fmt.Sprintf("DELETE FROM amigos WHERE aID = %s;", id))
	check("postUser PrepareQuery", err)
	res, err := smt.Exec()
	check("postUser Execute", err)
	n, err := res.RowsAffected()
	check("postUser CountRows", err)
	fmt.Fprintln(w, "Data deleted successfully, affected rows: ", n)
}

// CHECK ERRORS AND PRINT THEM OUT
func check(from string, err error) {
	if err != nil {
		log.Printf("%s: %s", from, err.Error())
	}
}
