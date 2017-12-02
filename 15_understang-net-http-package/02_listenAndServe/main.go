package main

import (
	"fmt"
	"net/http"
)

type hand int

func (h hand) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "printing some code...")
}

func main() {
	var h hand
	http.ListenAndServe(":8080", h)
}
