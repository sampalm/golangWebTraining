package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	tpl, err := template.ParseGlob("*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.ExecuteTemplate(os.Stdout, "one.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.ExecuteTemplate(os.Stdout, "two.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.ExecuteTemplate(os.Stdout, "vespa.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

}
