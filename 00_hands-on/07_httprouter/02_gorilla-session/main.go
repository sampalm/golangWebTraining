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
	if session.Values["logged-in"] == nil || session.Values["logged-in"] == "" {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	tpl.ExecuteTemplate(w, "index.html", ud)
}

func profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check session status
	session, _ := store.Get(r, "session")
	if session.Values["logged-in"] == nil || session.Values["logged-in"] == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "user-page.html", nil)
	return
	// TODO: Check if session is valid and if users exists

}

func sign(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check session status
	session, _ := store.Get(r, "session")
	if session.Values["logged-in"] == nil || session.Values["logged-in"] == "" {
		tpl.ExecuteTemplate(w, "sign-up.html", nil)
		return
	}

	fmt.Fprintf(w, "Your email: %v", session.Values["email"])

	// TODO: CHECK if sessesion is valid and if users exists

}

// ACTIONS HANDLER
func createAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check session status
	session, _ := store.Get(r, "session")
	if session.Values["logged-in"] != nil || session.Values["logged-in"] != "" {
		log.Println("dbUsers: ", dbUsers)
		log.Printf("Login: %v,\t Email: %v", session.Values["logged-in"], session.Values["email"])
		// http.Redirect(w, r, "/profile/", http.StatusSeeOther)
		return
	}

	// Check formValues
	if r.FormValue("email") == "" || r.FormValue("password") == "" {
		http.Error(w, "Insuficient information.", http.StatusBadRequest)
		return
	}

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
	session.Values["logged-in"] = 1

	// Save session
	session.Save(r, w)
	return
}
func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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
	//POST METHODS
	router.POST("/login/", login)
	router.POST("/logout/", logout)
	router.POST("/create-account/", createAccount)

	http.Handle("/", router)
}
