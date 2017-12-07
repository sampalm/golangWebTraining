package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./starting-files/public/")))
	http.HandleFunc("/templates/", index)
	http.ListenAndServe(":8080", nil)
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("./starting-files/templates/index.gohtml"))
}

func index(res http.ResponseWriter, req *http.Request) {

	err := tpl.ExecuteTemplate(res, "index.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}
