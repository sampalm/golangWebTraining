package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
)

var html string
var store = sessions.NewCookieStore([]byte("something-secret"))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// default cookie
		session, _ := store.Get(req, "logged-in")
		session.Values["status"] = 0

		// check log in
		if req.Method == "POST" {
			username := req.FormValue("username")
			password := req.FormValue("password")
			if password == "secret" {
				session.Values["status"] = 1
				session.Values["user"] = username
			}
		}

		// check log out
		if req.URL.Path == "/logout" {
			session.Values["status"] = 0
			http.Redirect(w, req, "/", http.StatusSeeOther)
		}

		// not logged in template
		if session.Values["status"] == 0 {
			html = `
				<!DOCTYPE html>
				<html>
					<head>
						<meta charset="UTF-8">
						<title>GOlang Web</title>
					</head>
					<body>
						<h1> Welcome </h1>
						<form method="POST">
							Your username: <input type="text" name="username">
							Your password: <input type="password" name="password">
							<br>
							<button type="submit">Submit</button>
						</form>
					</body>
				</html>
			`
		}

		// log in template
		if session.Values["status"] == 1 {
			html = `
				<!DOCTYPE html>
				<html>
					<head>
						<meta charset="UTF-8">
						<title>GOlang Web</title>
					</head>
					<body>
						<h1> Welcome, ` + fmt.Sprint(session.Values["user"]) + ` </h1>
						<h2><a href="/logout">Logout</a></h2>
					</body>
				</html>
			`
		}

		session.Save(req, w)
		io.WriteString(w, html)

	})

	http.ListenAndServe(":8080", nil)
}
