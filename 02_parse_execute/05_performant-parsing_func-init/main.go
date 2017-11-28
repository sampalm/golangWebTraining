package main

import (
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("*.gohtml"))
}

func main() {

	nf, err := os.Create("AllinOne.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	defer nf.Close()

	err = tpl.Execute(nf, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.ExecuteTemplate(nf, "one.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.ExecuteTemplate(nf, "two.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.ExecuteTemplate(nf, "vespa.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}
}
