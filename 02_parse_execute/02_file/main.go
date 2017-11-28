package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}

	nf, err := os.Create("tpl.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	defer nf.Close()

	err = tpl.Execute(nf, nil)
	if err != nil {
		log.Fatal(err)
	}
}
