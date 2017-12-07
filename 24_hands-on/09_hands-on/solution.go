package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./starting-files/")))
	http.ListenAndServe(":8080", nil)
}
