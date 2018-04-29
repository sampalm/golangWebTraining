package oauthweb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"

	"github.com/satori/go.uuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

var root = "http://localhost:8080"
var store = sessions.NewCookieStore([]byte("something-really-secret"))

func init() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login-github", loginHandler)
	http.HandleFunc("/oauth-github", oauthHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `<!DOCTYPE html>
	<html>
		<head></head>
		<body>
			<a href="/login-github">Login with Github.</a>
		</body>
	</hmtl>	
	`)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Generates a new unique ID
	id, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Makes a new cookie_session if not exits
	session, err := r.Cookie("oauth_session")
	if err != nil {
		session = &http.Cookie{
			Name:     "oauth_session",
			Value:    id.String(),
			HttpOnly: true,
		}
		http.SetCookie(w, session)
	}
	// Get oAuth for Github API
	auth := NewGitOAuth(fmt.Sprintf("%s/oauth-github", root), id.String())
	// Save oAuth into a memcache and use session id as a key
	ctx := appengine.NewContext(r)
	_, err = memcache.Get(ctx, session.Value)
	if err != nil {
		if err := memcache.Gob.Set(ctx, &memcache.Item{
			Key:    session.Value,
			Object: auth,
		}); err != nil {
			log.Printf("**********************************************************\nloginHandler Error: %v\n", err)
			return
		}
	}

	// Get user authentication url
	authURL, err := auth.GetAuthURI()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Redirect user to Github authentication page
	http.Redirect(w, r, authURL, 303)
}
func oauthHandler(w http.ResponseWriter, r *http.Request) {
	// Check if session exists
	session, err := r.Cookie("oauth_session")
	if err != nil {
		log.Println("oauthHandler Cookie error: " + err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	// Gen a variable to receive oauth pointer from memcache
	var auth *oAuth
	ctx := appengine.NewContext(r)
	// Check if memcache has been set
	_, err = memcache.Gob.Get(ctx, session.Value, &auth)
	if err != nil {
		log.Println("oauthHandler MemCache error: " + err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	// Parse URL values
	response, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("oauthHandler ParseQuery error: " + err.Error())
		return
	}
	// Check if authentication url has been changed.
	if auth.State != response.Get("state") {
		log.Println(fmt.Errorf("oauthHandler Error: oauth STATE undefined"))
		return
	}
	// Now try to generate a token access
	auth.Code = response.Get("code")
	if err := auth.GetAccessToken(ctx); err != nil {
		log.Println("oauthHandler GetAccessToken error: " + err.Error())
		return
	}
	// Now we have an api token and we can make requests to the api.
	// To test we'll get users email
	emails, err := auth.GetEmails(ctx)
	if err != nil {
		log.Println("oauthHandler GetEmails error: " + err.Error())
		return
	}

	if len(emails) == 0 {
		io.WriteString(w, "email not found")
		return
	}
	io.WriteString(w, "Your email: "+emails[0].Email)
}
