package main

import (
	"io/ioutil"
	"fmt"
	"io"
	"net/http"
)

func main(){
	http.HandleFunc("/", foo)
	http.Handle("favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, req *http.Request){
	var s string

	if req.Method == http.MethodPost {
		f, h, err := req.FormFile("q")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		fmt.Println("\n file: ", f, "\n header: ", h, "\n err: ", err)

		bs, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		s = string(bs)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf8")
	io.WriteString(w,`
		<form method="POST" enctype="multipart/form-data">
		 <input type="file" name="q">
		 <input type="submit">
		</form>	
	`+s)
}