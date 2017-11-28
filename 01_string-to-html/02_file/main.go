package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	name := "User Name"

	str := fmt.Sprint(`
			<!DOCTYPE html>
			<html lang="pt-br>
			<head>
			<meta charset="utf-8">
			<title>Hello World</title>
			</head>
			<body>
			<h1>` + name + `</h1>
			</body>
			</html>
		`)

	nf, err := os.Create("index.html")
	if err != nil {
		log.Fatal("error creating file")
	}
	defer nf.Close()

	io.Copy(nf, strings.NewReader(str))
}
