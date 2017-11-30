package main

import (
	"log"
	"os"
	"text/template"
	"strconv"
)

type person struct {
	Name string
	Age  int
}

func (p person) Code() string {
	c := p.Age * 42
	return "00" + strconv.Itoa(c)
}

func (p person) Even(n int) bool {
	c := n % 2
	if c == 0 {
		return true
	}

	return false
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

func main() {
	p := person{
		"Ian Fleming",
		32,
	}

	err := tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", p)
	if err != nil {
		log.Fatal(err)
	}
}
