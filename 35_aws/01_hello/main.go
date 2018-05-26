package main

import (
	"net/http"
	"io"
)

func main(){
	http.HandleFunc("/", hello)
	http.Handle("favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":80", nil)
}

func hello(w http.ResponseWriter, r *http.Request){
	io.WriteString(w, "Hello AWS")
}

