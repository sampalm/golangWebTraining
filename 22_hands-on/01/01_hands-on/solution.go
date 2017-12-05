package main

import (
	"fmt"
	"io"
	"net/http"
)

func home(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Welcome to home page !")
}
func d(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Welcome to dog page !")
}
func me(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hello, I'm Samuel.")
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/dog/", d)
	http.HandleFunc("/me/", me)

	http.ListenAndServe(":8080", nil)
}
