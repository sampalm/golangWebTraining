package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-word"))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "session")
		if req.FormValue("email") != "" {
			session.Values["email"] = req.FormValue("email")
		}
		session.Values["id"] = "Testing."

		session.Save(req, w)

		io.WriteString(w, `<!DOCTYPE html>
			<html>
			<body>
				<form method="POST">
					`+fmt.Sprintln(session.Values["email"])+`
					<input type="email" name="email">
					<input type="password" name="password">
					<button type="submit">Submit</button>
				</form>
			</body>
			</html>	
		`)
	})

	http.ListenAndServe(":8080", nil)
}
