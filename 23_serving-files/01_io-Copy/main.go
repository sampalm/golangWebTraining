package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/dog.jpg", dog)
	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "text/html; charset=utf-8")

	io.WriteString(res, `
	<img src="/dog.jpg">
	`)
}

func dog(res http.ResponseWriter, req *http.Request) {
	file, err := os.Open("dog.jpg")
	if err != nil {
		http.Error(res, "file not found", 404)
		return
	}
	defer file.Close()

	io.Copy(res, file)
}
