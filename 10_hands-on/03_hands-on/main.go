package main

import (
	"log"
	"os"
	"text/template"
)

type hotel struct {
	Name    string
	Address string
	City    string
	Zip     int
	region
}

type region struct {
	Region string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

func main() {
	hotels := []hotel{
		hotel{
			"Sea Blue Hotel",
			"Ocean Ave, Santa Monica",
			"CA",
			90401,
			region{
				"Southern",
			},
		},
		hotel{
			"Hyatt Residence Club Lake Tahoe",
			"Incline Way, Incline Village",
			"NV",
			89451,
			region{
				"Nothern",
			},
		},
		hotel{
			"Salinas Monterey Hotel",
			"Kern St, Salinas",
			"CA",
			93905,
			region{
				"Western",
			},
		},
		hotel{
			"Hotel Pacific",
			"Pacific St, Monterey",
			"CA",
			93940,
			region{
				"Central",
			},
		},
	}

	err := tpl.Execute(os.Stdout, hotels)
	if err != nil {
		log.Fatal(err)
	}
}
