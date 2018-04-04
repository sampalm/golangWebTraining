package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type UserData struct {
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Status    int
}

var dbUsers = map[string]UserData{}

var tpl *template.Template
var store = sessions.NewCookieStore([]byte("something-secret"))

// PAGE HANDLERS
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ud UserData
	log.Println(dbUsers)
	// Check session status
	session, _ := store.Get(r, "session")
	if session.Values["log-in"] == nil || session.Values["log-in"] == "" {

		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	session.Save(r, w)

	tpl.ExecuteTemplate(w, "index.html", ud)
}

func profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check session status
	session, _ := store.Get(r, "session")
	if session.Values["log-in"] == nil || session.Values["log-in"] == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var user UserData
	email := fmt.Sprintf("%s", session.Values["email"])

	if u, ok := dbUsers[email]; ok {
		user = u
	} else {
		http.Error(w, "Access Unauthorized", http.StatusUnauthorized)
		return
	}

	session.Save(r, w)

	tpl.ExecuteTemplate(w, "user-page.html", user)
	// TODO: Check if session is valid and if users exists

}

func sign(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check session status
	session, _ := store.Get(r, "session")
	if session.Values["log-in"] == nil || session.Values["log-in"] == "" {
		tpl.ExecuteTemplate(w, "sign-up.html", nil)
		return
	}

	session.Save(r, w)

	fmt.Fprintf(w, "Your email: %v", session.Values["email"])

	// TODO: CHECK if sessesion is valid and if users exists

}

// ACTIONS HANDLER
func createAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check session status
	session, _ := store.Get(r, "session")
	if r.FormValue("email") != "" {
		// Create account
		ud := UserData{
			Firstname: r.FormValue("first"),
			Lastname:  r.FormValue("last"),
			Email:     r.FormValue("email"),
			Password:  r.FormValue("password"),
		}
		dbUsers[ud.Email] = ud

		// Create session
		session.Values["email"] = ud.Email
		session.Values["log-in"] = true
	}

	// Save session
	session.Save(r, w)
	http.Redirect(w, r, "/profile/", http.StatusSeeOther)
}
func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, _ := store.Get(r, "session")

	// Get email if exists
	if email := r.FormValue("email"); email != "" {
		// Check if account already exists
		if u, ok := dbUsers[email]; ok {
			// Check if password match
			if u.Password == r.FormValue("password") {
				session.Values["email"] = r.FormValue("email")
				session.Values["log-in"] = true
			} else {
				http.Error(w, "Email invalid.", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Account doesn't exists.", http.StatusUnauthorized)
			return
		}
	}

	// Create session and make login
	session.Save(r, w)
	http.Redirect(w, r, "/profile/", http.StatusSeeOther)
}
func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, _ := store.Get(r, "session")
	if session.Values["log-in"] == true {
		session.Options.MaxAge = -1
	}

	session.Save(r, w)
	http.Redirect(w, r, "/sign/", http.StatusSeeOther)
}

func init() {
	// LOAD TEMPLATES
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	// HANDLERS FUNC
	router := httprouter.New()
	// GET METHODS
	router.GET("/", index)
	router.GET("/sign/", sign)
	router.GET("/profile/", profile)
	router.ServeFiles("/public/*filepath", http.Dir("public/"))
	router.GET("/logout/", logout)
	//POST METHODS
	router.POST("/login/", login)
	router.POST("/create-account/", createAccount)

	http.Handle("/", router)
}
