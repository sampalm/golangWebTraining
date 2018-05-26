package oauthweb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"

	"github.com/gorilla/sessions"

	"github.com/satori/go.uuid"
)

var root = "http://localhost:8080"
var store = sessions.NewCookieStore([]byte("something-really-secret"))

func init() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/oauth", oauthHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `<!DOCTYPE html>
	<html>
		<head></head>
		<body>
			<a href="/login">Login with Twitter.</a>
		</body>
	</hmtl>	
	`)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// Generates a new unique ID
	id, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Create a session
	session, err := r.Cookie("aouth")
	if err != nil {
		session = &http.Cookie{
			Name:  "oauth",
			Value: id.String(),
		}
	}
	http.SetCookie(w, session)
	// Create a twitter aouth
	redirectURI := fmt.Sprintf("%s/oauth", root)
	oAuth := NewTwitterOAuth(redirectURI, id.String()) // Change to NewGitOAuth() to use Github api
	// Generate authorization URL
	authURL, err := oAuth.GetAuthURI(ctx)
	if err != nil {
		log.Printf("loginHandler Error: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	// Create a memcache to save this session
	_, err = memcache.Get(ctx, session.Value)
	if err != nil {
		if err := memcache.Gob.Set(ctx, &memcache.Item{
			Key:    session.Value,
			Object: oAuth,
		}); err != nil {
			log.Printf("loginHandler memCache Error: %v", err)
			http.Error(w, "Something went wront, try again later.", 500)
			return
		}
	}
	// Redirect the user to the dropbox authenticate page
	http.Redirect(w, r, authURL, 302)
}
func oauthHandler(w http.ResponseWriter, r *http.Request) {
	// Check session state
	session, err := r.Cookie("oauth")
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	// Checks for authentication
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(fmt.Errorf("oauthHandler Error: failed to get HTTP response"))
		session.MaxAge = -1
		http.SetCookie(w, session)
		http.Redirect(w, r, "/", 302)
		return
	}
	// Get oauth from memcache
	var oauth *oAuthV4 // Change to *oAuth to use github api
	ctx := appengine.NewContext(r)
	_, err = memcache.Gob.Get(ctx, session.Value, &oauth)
	if err != nil {
		log.Printf("oauthHandler memCache Error: %v", err)
		http.Error(w, "Something went wrong, try again later.", 500)
		return
	}
	// Set oauth_verifier to get an access_token
	oauth.Verifier = values.Get("oauth_verifier")
	log.Printf("***** LOG REDIRECTING THE USER WITH SUCCESS: OTK=%s, NTK=%s, VERIFIER=%s\n", oauth.OTokenID, values.Get("oauth_token"), oauth.Verifier)
	if err := oauth.GetAccessToken(ctx); err != nil {
		log.Printf("oauthHandler Error: %v", err)
		http.Error(w, "Something went wrong, try again later.", 500)
		return
	}
	log.Println("READY TO REQUEST USER IDENTITY.")
	// Get user identity
	user, err := oauth.GetUser(ctx)
	if err != nil {
		log.Printf("oauthHandler Error: %v", err)
		http.Error(w, "Something went wrong, try again later.", 500)
		return
	}

	io.WriteString(w, `<!DOCTYPE html>
		<html>
		<head><title>Github oAuth</title></head>
		<body>
		<img src="`+user.Avatar+`" width="150px">
		<p>`+user.Name+`</p>
		<p>Email: `+user.Email+`</p>
		<ol>
			<li>ID: `+user.ID+`</li>
		</ol>
		</body>
		</html>
	`)
}
