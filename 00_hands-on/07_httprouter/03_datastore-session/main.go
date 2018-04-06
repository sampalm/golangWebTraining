package webapp

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type templateData struct {
	Name       string
	Email      string
	Subscribed time.Time
	Role       int
	errors     []error
}

var tpl *template.Template
var store = sessions.NewCookieStore([]byte("some-secret-string-here"))

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ud User

	userSession, _ := store.Get(r, "user-session")
	if userSession.Values["log-in"] == nil || userSession.Values["log-in"] == false {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	userSession.Save(r, w)
	tpl.ExecuteTemplate(w, "index.html", ud)
}

func sign(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userSession, _ := store.Get(r, "user-session")
	if userSession.Values["log-in"] == nil || userSession.Values["log-in"] == false {
		// Check if exists errors
		var messageErr string
		if ps.ByName("errors") != "" {
			messageErr = ps.ByName("errors")
		}

		tpl.ExecuteTemplate(w, "sign-up.html", messageErr)
		return
	}

	userSession.Save(r, w)
	http.Redirect(w, r, "/profile/", http.StatusSeeOther)
}

func profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userSession, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userSession.Values["log-in"] == nil || userSession.Values["log-in"] == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Define users role
	var userRole string
	switch userSession.Values["role"] {
	case 1:
		userRole = "User"
	case 2:
		userRole = "Administrator"
	case 3:
		userRole = "Owner"
	default:
		userRole = "User"
	}

	pData := map[string]string{
		"email": userSession.Values["email"].(string),
		"role":  userRole,
	}

	userSession.Save(r, w)
	tpl.ExecuteTemplate(w, "user-page.html", pData)
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userSession, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("email") == "" || r.FormValue("password") == "" {
		http.Redirect(w, r, "/sign/error=empty", http.StatusSeeOther)
		return
	}

	userEmail := r.FormValue("email")
	userPassword := r.FormValue("password")

	found, role, err := LoginAccount(r, userEmail, userPassword)
	if found != true {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userSession.Values["log-in"] = true
	userSession.Values["email"] = userEmail
	userSession.Values["role"] = role
	userSession.Save(r, w)
	http.Redirect(w, r, "/profile/", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userSession, _ := store.Get(r, "user-session")
	if userSession.Values["log-in"] == true {
		userSession.Options.MaxAge = -1
	}

	userSession.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userSession, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("email") == "" {
		http.Redirect(w, r, "/sign/error=email", http.StatusSeeOther)
		return
	} else if r.FormValue("password") != r.FormValue("password2") {
		http.Redirect(w, r, "/sign/error=password", http.StatusSeeOther)
		return
	} else if r.FormValue("password") == "" {
		http.Redirect(w, r, "/sign/error=password", http.StatusSeeOther)
		return
	}

	// Collect form data
	ud := User{
		Name:       r.FormValue("first"),
		Email:      r.FormValue("email"),
		Password:   r.FormValue("password"),
		Subscribed: time.Now(),
		Role:       1,
	}
	if errDatastore := SetAccount(r, ud); errDatastore != nil {
		http.Error(w, errDatastore.Error(), http.StatusInternalServerError)
		return
	}

	userSession.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userSession, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userSession.Values["log-in"] == nil || userSession.Values["log-in"] == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var td templateData
	usEmail := userSession.Values["email"].(string)
	user, err := GetAccount(r, usEmail)
	if err == nil {
		log.Println(user)
		// Get user info
		td.Name = user.Name
		td.Email = user.Email
		td.Subscribed = user.Subscribed
		td.Role = user.Role

		tpl.ExecuteTemplate(w, "update-account.html", td)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userSession.Save(r, w)
}

func updateAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userSession, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userSession.Values["log-in"] == nil || userSession.Values["log-in"] == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Retrieve user data from data store
	user, err := GetAccount(r, userSession.Values["email"].(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var newUser User
	// Check form values to see what will need to be updated
	// Verify if Firstname changed or is empty
	if r.FormValue("first") == "" || len(r.FormValue("first")) < 4 {
		http.Error(w, "Firstname too short or is empty.", http.StatusInternalServerError)
		return
	}
	// Set by default Username to the same from datastore
	newUser.Name = user.Name
	if r.FormValue("first") != user.Name {
		// Update Username from datastore
		newUser.Name = r.FormValue("first")
	}
	// Check if password was informed and if match
	if r.FormValue("old-password") == "" || HashPassword(r.FormValue("old-password"), nil) != user.Password {
		http.Error(w, "Password invalid or is empty.", http.StatusInternalServerError)
		return
	}
	// Set password to the same from datastore
	newUser.Password = user.Password
	// Check if password needs to be changed
	if r.FormValue("password") != "" {
		if r.FormValue("password") != r.FormValue("password2") {
			http.Error(w, "New Password doesn't match, if you don't want to change it just let this field empty.", http.StatusInternalServerError)
			return
		}
		// Change the password from datastore
		newUser.Password = HashPassword(r.FormValue("password"), nil)
	}
	// Fill in all other fields that will not change
	newUser.Email = user.Email
	newUser.Role = user.Role
	newUser.Subscribed = user.Subscribed
	// Put all the changes in datastore
	if status, err := UpdateAccount(r, newUser); status == false {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userSession.Save(r, w)
	http.Redirect(w, r, "/profile/", http.StatusSeeOther)
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/sign/", sign)
	router.GET("/sign/:errors", sign)
	router.GET("/profile/", profile)
	router.GET("/logout/", logout)
	router.GET("/update-account/", update)
	router.POST("/login/", login)
	router.POST("/create-account/", createAccount)
	router.POST("/update-account/", updateAccount)
	router.Handler("GET", "/favicon.ico", router.NotFound)

	http.Handle("/", router)
}
