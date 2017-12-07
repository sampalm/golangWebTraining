package main

import (
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/", hand)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func hand(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.URL)
	fmt.Fprintln(w, "go look at your terminal")
}