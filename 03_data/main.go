package main

import (
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

func main() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
	err := tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", 42)
	if err != nil {
		log.Fatal(err)
	}

}
