package oauthweb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"google.golang.org/appengine"

	"github.com/gorilla/sessions"
	"google.golang.org/appengine/urlfetch"

	"github.com/satori/go.uuid"
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
	session, err := store.Get(r, "oauth")
	if err != nil {
		log.Println(err)
		return
	}
	id, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	session.Values["state"] = id.String()

	values := make(url.Values)
	values.Add("client_id", "88bccd52adcc738aed04")
	values.Add("redirect_uri", fmt.Sprintf("%s/oauth-github", root))
	values.Add("scope", "user:email")
	values.Add("state", id.String())

	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("https://github.com/login/oauth/authorize?%s", values.Encode()), 303)

}
func oauthHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "oauth")
	if err != nil {
		log.Println(err)
		return
	}
	// Parse URL values
	response, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
		return
	}
	// Check code authenticity
	if response.Get("state"); session.Values["state"] != response.Get("state") {
		io.WriteString(w, "Cross site attack detected! \n")
		log.Printf("Your code: %v, Your sCode: %v", response.Get("state"), session.Values["state"])
		return
	}
	io.WriteString(w, "Making the authorization... ")

	values := make(url.Values)
	values.Add("client_id", "88bccd52adcc738aed04")
	values.Add("client_secret", "7888bc468d37fcc61bbae40ba51647aece853f60")
	values.Add("code", response.Get("code"))
	values.Add("state", response.Get("state"))

	token, err := accessToken(r, values)
	if err != nil {
		http.Error(w, err.Error(), 501)
	}
	email, err := getEmail(r, token)
	if err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
	session.Values["email"] = email
	session.Save(r, w)
	io.WriteString(w, "Your email: "+email)
	// TODO: Make redirect to users page
}
func accessToken(r *http.Request, values url.Values) (string, error) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	res, err := client.PostForm("https://github.com/login/oauth/access_token", values)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	response, err := url.ParseQuery(string(bs))
	if err != nil {
		return "", err
	}
	return response.Get("access_token"), nil
}
func getEmail(r *http.Request, token string) (string, error) {
	var data []struct {
		Email      string
		Verified   bool
		Primary    bool
		Visibility string
	}

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	res, err := client.Get("https://api.github.com/user/emails?access_token=" + token)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", fmt.Errorf("no data found")
	}
	return data[0].Email, nil
}
