package main

import (
	"io"
	"net/http"
)

func main(){
	http.HandleFunc("/", foo)
	http.Handle("favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, req *http.Request){
	fname := req.FormValue("fname")
	lname := req.FormValue("lname")
	check := req.FormValue("checkbox") == "on"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `
		<form method=post>
			<input type="text" name="fname" placeholder="Insert your first name"><br>
			<input type="text" name="lname" placeholder="Insert your last name"><br>
			<label for="check">Accept box: </label>
			<input type="checkbox" id="check" name="checkbox"><br>
			<input type="submit">		
		</form>
		<br>	
	` + fname +`<br>`+ lname+`<br>`)

	 if check {
		 io.WriteString(w, "Checkbox: true")
	 } else {
		io.WriteString(w, "Checkbox: false")
	 }
}