package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("./starting-files/templates/index.gohtml"))
}

func main() {
	http.HandleFunc("/templates/", indexFunc)
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./starting-files/public/"))))
	http.ListenAndServe(":8080", nil)
}

func indexFunc(res http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(res, "index.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}
