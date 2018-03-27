package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/satori/go.uuid"
)

type user struct {
	Username  string
	Firstname string
	Lastname  string
}

var tpl *template.Template
var dbUsers = map[string]user{}
var dbSessions = make(map[string]string)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", foo)
	http.Handle("favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {

	// Check if cookie exists
	cookie, err := r.Cookie("session")
	if err != nil {
		// Set a new unique hash and create "session" cookie
		sID, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

		http.SetCookie(w, cookie)
	}

	// If cookie session user = username.FormValue, get user data
	var u user
	if un, ok := dbSessions[cookie.Value]; ok {
		u = dbUsers[un]
	}

	// Process form submission
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		f := r.FormValue("firstname")
		l := r.FormValue("lastname")
		u := user{un, f, l}

		// ex: dbSessions[021203530654] = user01
		// dbUsers[user01] = {un, f, l}
		dbSessions[cookie.Value] = un
		dbUsers[un] = u
	}

	err = tpl.ExecuteTemplate(w, "index.gohtml", u)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}

}
