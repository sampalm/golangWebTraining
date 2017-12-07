package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/dog/", dog)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.ListenAndServe(":8080", nil)
}

func foo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "text/html; charset=utf-8")
	io.WriteString(res, "foo ran")
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.gohtml"))
}

func dog(res http.ResponseWriter, req *http.Request) {
	h1 := `This is from dog`
	res.Header().Set("Content-type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "index.gohtml", h1)
	if err != nil {
		log.Println(err)
	}
}
