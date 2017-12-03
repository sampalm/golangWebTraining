package main

import (
	"fmt"
	"net/http"
)

type hand int

func (h hand) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Key", "0007894321")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<h1>Write function in action ...</h1>")
}

func main() {
	var h hand
	http.ListenAndServe(":8080", h)
}
