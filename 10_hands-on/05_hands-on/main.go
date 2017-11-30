package main

import (
	"log"
	"os"
	"text/template"
)

type menu struct {
	Breakfast []string
	Lunch     []string
	Dinner    []string
}

type restaurant struct {
	Name string
	Menu menu
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

func main() {
	restaurants := []restaurant{
		restaurant{
			Name: "Restaurant 01",
			Menu: menu{
				[]string{"Eggs", "Bacon", "Coffe"},
				[]string{"Rice", "Potatos", "Sandwich"},
				[]string{"Meat", "Paste", "Fries"},
			},
		},
		restaurant{
			Name: "Restaurant 02",
			Menu: menu{
				[]string{"Eggs", "Bacon", "Coffe"},
				[]string{"Rice", "Potatos", "Sandwich"},
				[]string{"Meat", "Paste", "Fries"},
			},
		},
		restaurant{
			Name: "Restaurant 03",
			Menu: menu{
				[]string{"Eggs", "Bacon", "Coffe"},
				[]string{"Rice", "Potatos", "Sandwich"},
				[]string{"Meat", "Paste", "Fries"},
			},
		},
		restaurant{
			Name: "Restaurant 04",
			Menu: menu{
				[]string{"Eggs", "Bacon", "Coffe"},
				[]string{"Rice", "Potatos", "Sandwich"},
				[]string{"Meat", "Paste", "Fries"},
			},
		},
	}

	err := tpl.Execute(os.Stdout, restaurants)
	if err != nil {
		log.Fatal(err)
	}
}
