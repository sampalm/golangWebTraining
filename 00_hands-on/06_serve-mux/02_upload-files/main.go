package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", index)
	http.Handle("favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	if r.Method == "POST" {
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error file: %v", err)
		}
		defer file.Close()

		f, err := os.OpenFile("./uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("Error file: %v", err)
			return
		}
		defer f.Close()

		if _, err := io.Copy(f, file); err != nil {
			log.Printf("Error file: %v", err)
			return
		}

		io.WriteString(w, "File upload with success!")
		return

	}

	io.WriteString(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>DemoGCS</title>
		</head>
		<body>
			<form method="POST" enctype="multipart/form-data">
				<input type="file" name="file" accept=".jpg, .jpeg, .png">
				<button type="submit">Upload</button>
			</form>
		</body>
		</html>	
	`)

}
