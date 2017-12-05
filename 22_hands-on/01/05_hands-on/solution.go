package main

import (
	"html/template"
	"log"
	"net/http"
)

func home(res http.ResponseWriter, req *http.Request) {
	data := struct {
		Title   string
		Message string
	}{
		"HomePage",
		"Welcome to Homepage!",
	}

	err := tpl.ExecuteTemplate(res, "index.html", data)
	if err != nil {
		log.Println(err)
	}
}
func d(res http.ResponseWriter, req *http.Request) {
	data := struct {
		Title   string
		Message string
	}{
		"Dog Page",
		"Welcome to Dog page!",
	}

	err := tpl.ExecuteTemplate(res, "index.html", data)
	if err != nil {
		log.Println(err)
	}
}
func me(res http.ResponseWriter, req *http.Request) {
	data := struct {
		Title   string
		Message string
	}{
		"My Page",
		"Hello, my name is Samuel.",
	}

	err := tpl.ExecuteTemplate(res, "index.html", data)
	if err != nil {
		log.Println(err)
	}
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func main() {
	http.Handle("/", http.HandlerFunc(home))
	http.Handle("/dog/", http.HandlerFunc(d))
	http.Handle("/me/", http.HandlerFunc(me))

	http.ListenAndServe(":8080", nil)
}
