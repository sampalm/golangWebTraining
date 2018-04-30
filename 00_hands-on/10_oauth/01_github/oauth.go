package oauthweb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"google.golang.org/appengine/urlfetch"
)

// oAuth is the struct body to make any oauth request.
// Path must be the exacly callback URL used in your API configuration.
type oAuth struct {
	// App config
	ClientID   string
	SecretID   string
	RequestURI string
	TokenURI   string
	AuthURI    string

	// Api config
	Token string
	Path  string
	Code  string
	State string
}

// NewGitOAuth returns a default config to github api.
func NewGitOAuth(path string, uID string) *oAuth {
	var gitOAuth oAuth
	// Set app config
	gitOAuth.ClientID = "88bccd52adcc738aed04"
	gitOAuth.SecretID = "7888bc468d37fcc61bbae40ba51647aece853f60"
	gitOAuth.RequestURI = "https://api.github.com"
	gitOAuth.TokenURI = "https://github.com/login/oauth/access_token"
	gitOAuth.AuthURI = "https://github.com/login/oauth/authorize"
	// Set api config
	gitOAuth.Path = path
	gitOAuth.State = uID
	return &gitOAuth
}

type user struct {
	ID        int    `json:"id"`
	Avatar    string `json:"avatar_url"`
	Profile   string `json:"html_url"`
	Email     string `json:"email"`
	Username  string `json:"login"`
	Name      string `json:"name"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}
type email struct {
	Email      string `json:"email"`
	Verified   bool   `json:"verified"`
	Primary    bool   `json:"primary"`
	Visibility string `json:"visibility"`
}

// GetAuth gets the authorization of user to use the api.
func (auth *oAuth) GetAuthURI() (string, error) {
	switch {
	case auth.ClientID == "":
		return "", fmt.Errorf("GetAuth Error: oAuth ClientID undefined, you need to define it before use oAuth requests")
	case auth.Path == "":
		return "", fmt.Errorf("GetAuth Error: oAuth STATE undefined, you need to define it before use oAuth requests")
	case auth.State == "":
		return "", fmt.Errorf("GetAuth Error: oAuth STATE undefined, you need to define it before use oAuth requests")
	}

	values := make(url.Values)
	values.Add("client_id", auth.ClientID)
	values.Add("redirect_uri", auth.Path)
	values.Add("scope", "user")
	values.Add("state", auth.State)

	return fmt.Sprintf("%s?%s", auth.AuthURI, values.Encode()), nil
}

// GetAccessToken gets the access token to make requests using apis.
func (auth *oAuth) GetAccessToken(ctx context.Context) error {
	switch {
	case auth.ClientID == "":
		return fmt.Errorf("GetAccessToken Error: oAuth ClientID undefined, you need to define it before use oAuth requests")
	case auth.SecretID == "":
		return fmt.Errorf("GetAccessToken Error: oAuth SecretID undefined, you need to define it before use oAuth requests")
	case auth.Code == "":
		return fmt.Errorf("GetAccessToken Error: oAuth CODE undefined, you need to define it before use oAuth requests")
	case auth.State == "":
		return fmt.Errorf("GetAccessToken Error: oAuth STATE undefined, you need to define it before use oAuth requests")
	}

	client := urlfetch.Client(ctx)

	// Set URL Values
	values := make(url.Values)
	values.Add("client_id", auth.ClientID)
	values.Add("client_secret", auth.SecretID)
	values.Add("code", auth.Code)
	values.Add("state", auth.State)

	res, err := client.PostForm(auth.TokenURI, values)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	response, err := url.ParseQuery(string(bs))
	if err != nil {
		return err
	}
	// Set token access
	auth.Token = response.Get("access_token")
	return nil
}

func (auth *oAuth) GetEmails(ctx context.Context) ([]email, error) {
	switch {
	case auth.Token == "":
		return nil, fmt.Errorf("GetEmails Error: oAuth TOKEN undefined, you need to define it before use oAuth requests")
	case auth.RequestURI == "":
		return nil, fmt.Errorf("GetEmails Error: oAuth RequestURI undefined, you need to define it before use oAuth requests")
	}

	var data []email
	requestURL := fmt.Sprintf("%s/user/emails?access_token=%s", auth.RequestURI, auth.Token)
	client := urlfetch.Client(ctx)
	res, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("GetEmails GetEmail Error: %v", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("GetEmails DecodeEmail Error: %v", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("GetEmails DataEmail Error: user emails not found")
	}
	return data, nil
}

func (auth *oAuth) GetUser(ctx context.Context) (user, error) {
	var u user
	switch {
	case auth.Token == "":
		return u, fmt.Errorf("GetUser Error: oAuth TOKEN undefined, you need to define it before use oAuth requests")
	case auth.RequestURI == "":
		return u, fmt.Errorf("GetUser Error: oAuth RequestURI undefined, you need to define it before use oAuth requests")
	}

	requestURL := fmt.Sprintf("%s/user?access_token=%s", auth.RequestURI, auth.Token)
	client := urlfetch.Client(ctx)
	res, err := client.Get(requestURL)
	if err != nil {
		return u, fmt.Errorf("GetUser GetUser Error: %v", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return u, fmt.Errorf("GetUser DecodeUser Error: %v", err)
	}
	return u, nil
}
