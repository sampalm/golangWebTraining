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
	gitOAuth.RequestURI = "https://api.github.com/user/emails?access_token="
	gitOAuth.TokenURI = "https://github.com/login/oauth/access_token"
	gitOAuth.AuthURI = "https://github.com/login/oauth/authorize"
	// Set api config
	gitOAuth.Path = path
	gitOAuth.State = uID
	return &gitOAuth
}

// AuthResponse returns the body of a request result.
// Data is a slice of any struct that implements data interface.
type AuthResponse struct {
	Body  data
	Error error
}
type data interface {
	getData() ([]string, error)
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
	values.Add("scope", "user:email")
	values.Add("state", auth.State)

	return fmt.Sprintf("%s?%s", auth.AuthURI, values.Encode()), nil
}

// GetAccessToken gets the access token to make requests using apis.
func (auth *oAuth) GetAccessToken(ctx context.Context) error {
	switch {
	case auth.ClientID == "":
		return fmt.Errorf("GetAuth Error: oAuth ClientID undefined, you need to define it before use oAuth requests")
	case auth.SecretID == "":
		return fmt.Errorf("GetAuth Error: oAuth SecretID undefined, you need to define it before use oAuth requests")
	case auth.Code == "":
		return fmt.Errorf("GetAuth Error: oAuth CODE undefined, you need to define it before use oAuth requests")
	case auth.State == "":
		return fmt.Errorf("GetAuth Error: oAuth STATE undefined, you need to define it before use oAuth requests")
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
		return nil, fmt.Errorf("GetAuth Error: oAuth TOKEN undefined, you need to define it before use oAuth requests")
	case auth.RequestURI == "":
		return nil, fmt.Errorf("GetAuth Error: oAuth RequestURI undefined, you need to define it before use oAuth requests")
	}

	var data []email
	client := urlfetch.Client(ctx)

	res, err := client.Get(auth.RequestURI + auth.Token)
	if err != nil {
		return nil, fmt.Errorf("GetAuth GetEmail Error: %v", err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("GetAuth DecodeEmail Error: %v", err)
	}
	return data, nil
}
