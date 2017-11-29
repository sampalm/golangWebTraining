package main

import (
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

func main() {
	slc := []string{"Jesus", "Bhudda", "MLK", "Gandhi"}

	err := tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", slc)
	if err != nil {
		log.Fatal(err)
	}
}
