package main

import (
	"log"
	"os"
	"text/template"
	"time"
)

var tpl *template.Template

func monthDayYear(t time.Time) string {
	return t.Format("02/01/2006")
}

func monthDayYearHour(t time.Time) string {
	return t.Format("15:04:05 - 02/01/2006")
}

func hourMinuteSecond(t time.Time) string {
	return t.Format("15:04:05s")
}

var fm = template.FuncMap{
	"fdateMDY":  monthDayYear,
	"fdateMDYH": monthDayYearHour,
	"fdateHMS":  hourMinuteSecond,
}

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseFiles("tpl.gohtml"))
}

func main() {
	err := tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", time.Now())
	if err != nil {
		log.Fatal(err)
	}
}
