package main

import (
	"fmt"
	"os"
	"text/template"
)

var tpl = template.New("template.gohtml")

func main() {
	tpl = tpl.Funcs(template.FuncMap{
		"sumNums": func(nums []int) int {
			var tot int
			lg := len(nums)
			for _, num := range nums {
				tot += num
			}
			return tot / lg
		},
		"subNums": func(nums []int) int {
			var tot int
			lg := len(nums)
			for _, num := range nums {
				tot -= num
			}

			return tot / lg
		},
	})

	tpl, err := tpl.ParseFiles("template.gohtml")
	if err != nil {
		fmt.Printf("An error occours: %v", err)
		return
	}

	err = tpl.Execute(os.Stdout, []int{5, 6, 8, 9, 2})
	if err != nil {
		fmt.Printf("An error occours: %v", err)
		return
	}
}
