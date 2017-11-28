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

	err = tpl.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatal(err)
	}
}
