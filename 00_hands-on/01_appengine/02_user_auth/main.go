package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	ctx := appengine.NewContext(r)
	u := user.Current(ctx)

	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
		return
	}

	// Log user info
	log.Infof(ctx, "User logged: %s", u)
	log.Infof(ctx, "User admin logged: %s", u.Admin)

	url, _ := user.LogoutURL(ctx, "/")
	fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)<br>`, u, url)
	fmt.Fprint(w, "<a href='/admin'>Admin area</a>")
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	ctx := appengine.NewContext(r)
	u := user.Current(ctx)

	// If user didn't exists
	if u == nil {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	// If user exists but isn't admin
	if !u.Admin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Logs user info
	log.Infof(ctx, "User logged: %s", u)
	log.Infof(ctx, "User admin logged: %s", u.Admin)

	fmt.Fprintf(w, `Welcome, admin user %s!`, u)
	fmt.Fprint(w, "<br><a href='/'>Return to user area</a>")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/admin", adminHandler)
	appengine.Main()
}
