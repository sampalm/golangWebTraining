package main

import (
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

type sage struct {
	Name  string
	Motto string
}

type car struct {
	Manufacturer string
	Model       string
	Doors        int
}

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

func main() {
	b := sage{
		Name:  "Buddha",
		Motto: "The belief of no beliefs",
	}
	g := sage{
		Name:  "Ghandi",
		Motto: "Be the change",
	}
	j := sage{
		Name:  "Jesus",
		Motto: "Love all",
	}
	f := car{
		Manufacturer: "Ford",
		Model:       "F150",
		Doors:        2,
	}
	c := car{
		Manufacturer: "Toyota",
		Model:       "Corolla",
		Doors:        4,
	}

	sages := []sage{b, g, j}
	cars := []car{f, c}

	data := struct {
		Wisdom    []sage
		Transport []car
	}{
		Wisdom:    sages,
		Transport: cars,
	}

	err := tpl.Execute(os.Stdout, data)
	if err != nil {
		log.Fatal(err)
	}

}
