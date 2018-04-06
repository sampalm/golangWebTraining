package main

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
)

var tpl *template.Template

type User struct {
	Email    string
	UserName string `datastore:"-"`
	Password string
}

func init() {
	tpl = template.Must(template.ParseGlob("*.html"))
	http.HandleFunc("/", index)
	http.HandleFunc("/about/", about)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	msg := "Data inserted into database!"
	ctx := appengine.NewContext(r)

	NewUser := User{
		Email:    "samuelpalmeira@outlook.com",
		UserName: "sampalm",
		Password: "123456",
	}
	key := datastore.NewKey(ctx, "Users", NewUser.UserName, 0, nil)
	key, err := datastore.Put(ctx, key, &NewUser)
	if err != nil {
		msg = "Error adding to database."
	}

	tpl.ExecuteTemplate(w, "about.html", msg)
}
