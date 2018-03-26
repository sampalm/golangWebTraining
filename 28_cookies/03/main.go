package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var count int

func main() {
	http.HandleFunc("/", set)
	http.HandleFunc("/read", read)
	http.Handle("favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func set(w http.ResponseWriter, req *http.Request) {
	count++
	http.SetCookie(w, &http.Cookie{
		Name:  "count-cookie",
		Value: strconv.Itoa(count),
	})

	fmt.Fprintln(w, "COOKIE WRITTE - CHECK YOUR BROWSER")
}

func read(w http.ResponseWriter, req *http.Request) {
	ck, err := req.Cookie("count-cookie")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Fprintln(w, "YOUR COOKIE: ", ck)
	}
}
