package oauthweb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"google.golang.org/appengine/urlfetch"
)

// oAuth is to use Github API
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

// oAuthV2 is tp use Dropbox API
type oAuthV2 struct {
	// App config
	Type       string // token or code
	ClientID   string // the app key
	SecretID   string // secret app key
	AuthURI    string
	RequestURI string // api request url
	TokenURI   string // access_token request url

	// Api config
	Path    string // redirect_uri
	State   string
	Code    string
	Token   string
	Account string
}

// oAuthV3 is to use Facebooks API
type oAuthV3 struct {
	// App config
	ClientID   string // the app key
	SecretID   string // secret app key
	AuthURI    string
	RequestURI string // api request url
	TokenURI   string // access_token request url

	// Api config
	Path  string // redirect_uri
	State string
	Code  string
	Token string
}

// oAuthV4 is to use Twitter API
type oAuthV4 struct {
	// App config
	ClientID        string // the app key
	SecretID        string // secret app key
	OTokenID        string // app access token
	OTokenSecretID  string // app secret access token
	AuthURI         string
	RequestURI      string // api request url
	RequestTokenURI string // request oauth_token
	TokenURI        string // access_token request url

	// Api config
	Path     string // redirect_uri
	ANonce   string // oauth_nonce ID
	Verifier string // oauth_verifier ID

	// User config
	UserID     string // twitter user_id
	ScreenName string // twitter screen_name
}

// NewGitOAuth returns a default config to github api.
func NewGitOAuth(path string, uID string) *oAuth {
	var gitOAuth oAuth
	// Set app config
	gitOAuth.ClientID = "XXXXXXXXXX"
	gitOAuth.SecretID = "XXXXXXXXXXXXXXXXXX"
	gitOAuth.RequestURI = "https://api.github.com"
	gitOAuth.TokenURI = "https://github.com/login/oauth/access_token"
	gitOAuth.AuthURI = "https://github.com/login/oauth/authorize"
	// Set api config
	gitOAuth.Path = path
	gitOAuth.State = uID
	return &gitOAuth
}

// NewBoxOAuth returns a default config to dropbox api.
func NewBoxOAuth(path string, uID string) *oAuthV2 {
	var boxOAuth oAuthV2
	// Set app config
	boxOAuth.Type = "code"
	boxOAuth.ClientID = "XXXXXXXXXX"
	boxOAuth.SecretID = "XXXXXXXXXXXXXXXXXXXX"
	boxOAuth.AuthURI = "https://www.dropbox.com/oauth2/authorize"
	boxOAuth.RequestURI = "https://api.dropboxapi.com/2"
	boxOAuth.TokenURI = "https://api.dropboxapi.com/oauth2/token"
	// Set api config
	boxOAuth.Path = path
	boxOAuth.State = uID
	return &boxOAuth
}

// NewFaceOAuth returns a default config to facebook api.
func NewFaceOAuth(path string, uID string) *oAuthV3 {
	var faceOAuth oAuthV3
	// Set app config
	faceOAuth.ClientID = "XXXXXXXXXX"
	faceOAuth.SecretID = "XXXXXXXXXXXXXXXXXXXX"
	faceOAuth.AuthURI = "https://www.facebook.com/v2.12/dialog/oauth"
	faceOAuth.TokenURI = "https://graph.facebook.com/v2.12/oauth/access_token"
	// Set api config
	faceOAuth.Path = path
	faceOAuth.State = uID
	return &faceOAuth
}

// NewTwitterOAuth returns a default config to twitter api
func NewTwitterOAuth(path string, uID string) *oAuthV4 {
	var twitterOAuth oAuthV4
	// Set app config
	twitterOAuth.ClientID = "6zE9MU9bBTM2ACVowu0m6cWTz"
	twitterOAuth.SecretID = "R3WGD7KKxvTWTEKvg5sijLWrB5GQfVRWNCA3hV6NLYOHstYv7C"
	twitterOAuth.OTokenID = "979769837894426624-sKd8jnpGnIopcY9Y8NGQwPolV939ecv"
	twitterOAuth.OTokenSecretID = "mAWVDS5Xf0sNje75xzd3nhlOnOlh5csnrnofhpD6fk4Aj"
	twitterOAuth.AuthURI = "https://api.twitter.com/oauth/authenticate"
	twitterOAuth.RequestTokenURI = "https://api.twitter.com/oauth/request_token"
	twitterOAuth.TokenURI = "https://api.twitter.com/oauth/access_token"
	twitterOAuth.RequestURI = "https://api.twitter.com/1.1"
	// Set api config
	twitterOAuth.ANonce = uID
	twitterOAuth.Path = path
	return &twitterOAuth
}

type user struct {
	ID       string
	Name     string
	Username string
	Email    string
	Avatar   string
}

type email struct {
	Email      string `json:"email"`
	Verified   bool   `json:"verified"`
	Primary    bool   `json:"primary"`
	Visibility string `json:"visibility"`
}

// GetAuth gets the authorization of user and returns an access token which is used to make calls to the api.
func (auth *oAuth) GetAuthURI() (string, error) { // Github Autorizate
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
func (auth *oAuthV2) GetAuthURI() (string, error) { // Dropbox Authorizate
	switch {
	case auth.Type == "":
		return "", fmt.Errorf("GetAuth Error: oAuthV2 TYPE undefined, you need to define it before use oAuthV2 requests")
	case auth.ClientID == "":
		return "", fmt.Errorf("GetAuth Error: oAuthV2 ClientID undefined, you need to define it before use oAuthV2 requests")
	case auth.Path == "":
		return "", fmt.Errorf("GetAuth Error: oAuthV2 STATE undefined, you need to define it before use oAuthV2 requests")
	case auth.State == "":
		return "", fmt.Errorf("GetAuth Error: oAuthV2 STATE undefined, you need to define it before use oAuthV2 requests")
	}

	values := make(url.Values)
	values.Add("response_type", auth.Type)
	values.Add("client_id", auth.ClientID)
	values.Add("state", auth.State)
	values.Add("redirect_uri", auth.Path)
	return fmt.Sprintf("%s?%s", auth.AuthURI, values.Encode()), nil
}
func (auth *oAuthV3) GetAuthURI() (string, error) { // Facebook Authorizate
	switch {
	case auth.ClientID == "":
		return "", fmt.Errorf("GetAuth Error: oAuthV3 ClientID undefined, you need to define it before use oAuthV3 requests")
	case auth.Path == "":
		return "", fmt.Errorf("GetAuth Error: oAuthV3 PATH undefined, you need to define it before use oAuthV3 requests")
	case auth.Code == "":
		return "", fmt.Errorf("GetAuth Error: oAuthV3 CODE undefined, you need to define it before use oAuthV3 requests")
	}

	values := make(url.Values)
	values.Add("client_id", auth.ClientID)
	values.Add("redirect_uri", auth.Path)
	values.Add("client_secret", auth.SecretID)
	values.Add("code", auth.Code)
	return fmt.Sprintf("%s?%s", auth.AuthURI, values.Encode()), nil
}
func (auth *oAuthV4) GetAuthURI(ctx context.Context) (string, error) { // Twitter Authorizate
	switch {
	case auth.ClientID == "":
		return "", fmt.Errorf("GetAuthURI Error: oAuthV4 ClientID undefined, you need to define it before use oAuthV4 requests")
	case auth.Path == "":
		return "", fmt.Errorf("GetAuthURI Error: oAuthV4 PATH undefined, you need to define it before use oAuthV4 requests")
	case auth.ANonce == "":
		return "", fmt.Errorf("GetAuthURI Error: oAuthV4 ANonce undefined, you need to define it before use oAuthV4 requests")
	}
	// Get oauth_signature and header values
	// Set header values
	header := EncodeNewHeader(EncodeSignature("POST", auth.RequestTokenURI, auth, nil))
	// Makes client
	client := urlfetch.Client(ctx)
	// Makes the http request
	req, err := http.NewRequest("POST", auth.RequestTokenURI, nil)
	req.Header.Set("Authorization", header)
	if err != nil {
		return "", fmt.Errorf("GetAuthURI Error: %s", err.Error())
	}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("GetAuthURI Error: %s", err.Error())
	}
	defer res.Body.Close()

	// RESPONSE BODY
	if res.StatusCode != 200 {
		return "", fmt.Errorf("GetAuthURI Error: failed to get response body")
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("GetAuthURI Error: %s", err.Error())
	}
	response, err := url.ParseQuery(string(bs))
	if err != nil {
		return "", fmt.Errorf("GetAuthURI ParseQuery Error: %s", err.Error())
	}
	// SET TOKEN VALUES
	if response.Get("oauth_callback_confirmed") != "true" {
		return "", fmt.Errorf("GetAuthURI ParseQuery Error: request not authorized")
	}
	auth.OTokenID = response.Get("oauth_token")
	auth.OTokenSecretID = response.Get("oauth_token_secret")
	log.Printf("***** LOG OBTAIN A REQUEST TOKEN WITH SUCCESS: TK=%s, STK=%s\n", auth.OTokenID, auth.OTokenSecretID)
	// RETURN AUTHENTICATE URL TO THE USER
	return fmt.Sprintf("%s?oauth_token=%s", auth.AuthURI, auth.OTokenID), nil
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
		return fmt.Errorf("GetAccessToken POST Error: %v", err)
	}
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("GetAccessToken ReadBODY Error: %v", err)
	}
	response, err := url.ParseQuery(string(bs))
	if err != nil {
		return fmt.Errorf("GetAccessToken ParseQUERY Error: %v", err)
	}
	// Set token access
	auth.Token = response.Get("access_token")
	return nil
}
func (auth *oAuthV2) GetAccessToken(ctx context.Context) error {
	switch {
	case auth.ClientID == "":
		return fmt.Errorf("GetAccessTokenV2 Error: oAuthV2 ClientID undefined, you need to define it before use oAuth requests")
	case auth.SecretID == "":
		return fmt.Errorf("GetAccessTokenV2 Error: oAuthV2 SecretID undefined, you need to define it before use oAuth requests")
	case auth.Code == "":
		return fmt.Errorf("GetAccessTokenV2 Error: oAuthV2 CODE undefined, you need to define it before use oAuth requests")
	case auth.State == "":
		return fmt.Errorf("GetAccessTokenV2 Error: oAuthV2 STATE undefined, you need to define it before use oAuth requests")
	}

	// Set URL params
	client := urlfetch.Client(ctx)
	values := make(url.Values)
	values.Add("code", auth.Code)
	values.Add("grant_type", "authorization_code")
	values.Add("client_id", auth.ClientID)
	values.Add("client_secret", auth.SecretID)
	values.Add("redirect_uri", auth.Path)

	res, err := client.PostForm(auth.TokenURI, values)
	if err != nil {
		log.Printf("GetAccessTokenV2 POST Error: %v", err)
		return fmt.Errorf("GetAccessTokenV2 POST Error: %v", err)
	}
	defer res.Body.Close()
	// Get the response content
	var response struct {
		Token   string `json:"access_token"`
		Account string `json:"account_id"`
		ID      string `json:"uid"`
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		log.Printf("GetAccessTokenV2 DECODE Error: %v", err)
		return fmt.Errorf("GetAccessTokenV2 DECODE Error: %v", err)
	}
	if response.Token == "" {
		log.Printf("GetAccessTokenV2 DECODE Error: given access token is invalid")
		return fmt.Errorf("GetAccessTokenV2 DECODE Error: given access token is invalid")
	}
	auth.Token = response.Token
	auth.ClientID = response.Account
	return nil
}
func (auth *oAuthV3) GetAccessToken(ctx context.Context) error {
	switch {
	case auth.ClientID == "":
		return fmt.Errorf("GetAccessTokenV3 Error: oAuthV3 ClientID undefined, you need to define it before use oAuth requests")
	case auth.SecretID == "":
		return fmt.Errorf("GetAccessTokenV3 Error: oAuthV3 SecretID undefined, you need to define it before use oAuth requests")
	case auth.Code == "":
		return fmt.Errorf("GetAccessTokenV3 Error: oAuthV3 CODE undefined, you need to define it before use oAuth requests")
	case auth.State == "":
		return fmt.Errorf("GetAccessTokenV3 Error: oAuthV3 STATE undefined, you need to define it before use oAuth requests")
	}

	// Set URL params
	client := urlfetch.Client(ctx)
	values := make(url.Values)
	values.Add("code", auth.Code)
	values.Add("grant_type", "authorization_code")
	values.Add("client_id", auth.ClientID)
	values.Add("client_secret", auth.SecretID)
	values.Add("redirect_uri", auth.Path)

	res, err := client.PostForm(auth.TokenURI, values)
	if err != nil {
		log.Printf("GetAccessTokenV3 POST Error: %v", err)
		return fmt.Errorf("GetAccessTokenV3 POST Error: %v", err)
	}
	defer res.Body.Close()

	// Get the response content
	var response struct {
		Token    string `json:"access_token"`
		Type     string `json:"token_type"`
		Duration string `json:"expires_in"`
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		log.Printf("GetAccessTokenV3 DECODE Error: %v", err)
		return fmt.Errorf("GetAccessTokenV3 DECODE Error: %v", err)
	}
	if response.Token == "" {
		log.Printf("GetAccessTokenV3 DECODE Error: given access token is invalid")
		return fmt.Errorf("GetAccessTokenV3 DECODE Error: given access token is invalid")
	}

	auth.Token = response.Token
	return nil
}

func (auth *oAuthV4) GetAccessToken(ctx context.Context) error {
	switch {
	case auth.ClientID == "":
		return fmt.Errorf("GetAccessTokenV3 Error: oAuthV4 ClientID undefined, you need to define it before use oAuthV4 requests")
	case auth.SecretID == "":
		return fmt.Errorf("GetAccessTokenV3 Error: oAuthV4 SecretID undefined, you need to define it before use oAuthV4 requests")
	case auth.ANonce == "":
		return fmt.Errorf("GetAccessTokenV3 Error: oAuthV4 AuthNonce undefined, you need to define it before use oAuthV4 requests")
	case auth.OTokenID == "":
		return fmt.Errorf("GetAccessTokenV3 Error: oAuthV4 AccessToken undefined, you need to define it before use oAuthV4 requests")
	}
	log.Printf("***** LOG COVERTING THE REQUEST TOKEN [BODY]: OAUTH_TOKEN=%s, OAUTH_SECRET_TOKEN=%s\n", auth.OTokenID, auth.OTokenSecretID)
	client := urlfetch.Client(ctx)
	header := EncodeNewHeader(EncodeSignature("POST", auth.RequestTokenURI, auth, nil))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?oauth_verifier=%s", auth.TokenURI, auth.Verifier), nil)
	if err != nil {
		log.Printf("GetAccessTokenV3 REQUEST Error: %v", err)
		return fmt.Errorf("GetAccessTokenV3 REQUEST Error: %v", err)
	}
	req.Header.Set("Authorization", header)
	res, err := client.Do(req)
	if err != nil {
		log.Printf("GetAccessTokenV3 CLIENT Error: %v", err)
		return fmt.Errorf("GetAccessTokenV3 CLIENT Error: %v", err)
	}
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("GetAccessTokenV3 READBODY Error: %v", err)
		return fmt.Errorf("GetAccessTokenV3 READBODY Error: %v", err)
	}
	log.Println("BODY RESPONSE: ", string(bs))
	response, err := url.ParseQuery(string(bs))
	if err != nil {
		return fmt.Errorf("GetAccessTokenV3 ParseQuery Error: %s", err.Error())
	}
	log.Println("QUERY BODY: ", response)
	// SET TOKEN VALUES
	auth.OTokenID = response.Get("oauth_token")
	auth.OTokenSecretID = response.Get("oauth_token_secret")
	auth.UserID = response.Get("user_id")
	auth.ScreenName = response.Get("screen_name")
	log.Printf("***** LOG COVERTING THE REQUEST TOKEN WITH SUCCESS: OTK=%s, STK=%s\n", response.Get("oauth_token"), response.Get("oauth_token_secret"))
	return nil
	// TODO: GET QUE USER IDENTITY USING "GET account/verify_credentials"

}

func (auth *oAuth) GetEmails(ctx context.Context) ([]email, error) {
	switch {
	case auth.Token == "":
		return nil, fmt.Errorf("GetEmails Error: oAuth TOKEN undefined, you need to define it before use oAuthV2 requests")
	case auth.RequestURI == "":
		return nil, fmt.Errorf("GetEmails Error: oAuth RequestURI undefined, you need to define it before use oAuthV2 requests")
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
	var u struct {
		Account     int    `json:"id"`
		GivenName   string `json:"login"`
		DisplayName string `json:"name"`
		Email       string `json:"email"`
		Avatar      string `json:"avatar_url"`
		Country     string `json:"location"`
	}
	switch {
	case auth.Token == "":
		return user{}, fmt.Errorf("GetUser Error: oAuth TOKEN undefined, you need to define it before use oAuth requests")
	case auth.RequestURI == "":
		return user{}, fmt.Errorf("GetUser Error: oAuth RequestURI undefined, you need to define it before use oAuth requests")
	}

	requestURL := fmt.Sprintf("%s/user?access_token=%s", auth.RequestURI, auth.Token)
	client := urlfetch.Client(ctx)
	res, err := client.Get(requestURL)
	if err != nil {
		return user{}, fmt.Errorf("GetUser GetUser Error: %v", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return user{}, fmt.Errorf("GetUser DecodeUser Error: %v", err)
	}
	// Return user data
	data := user{
		ID:       strconv.Itoa(u.Account),
		Name:     u.DisplayName,
		Username: u.DisplayName,
		Email:    u.Email,
		Avatar:   u.Avatar,
	}
	return data, nil
}
func (auth *oAuthV2) GetUser(ctx context.Context) (user, error) {
	var u struct {
		Account string `json:"account_id"`
		Name    struct {
			GivenName   string `json:"given_name"`
			DisplayName string `json:"display_name"`
		}
		Email   string
		Avatar  string `json:"profile_photo_url"`
		Country string `json:"country"`
	}
	switch {
	case auth.Token == "":
		return user{}, fmt.Errorf("GetUserV2 Error: oAuthV2 TOKEN undefined, you need to define it before use oAuthV2 requests")
	case auth.RequestURI == "":
		return user{}, fmt.Errorf("GetUser Error: oAuth RequestURI undefined, you need to define it before use oAuth requests")
	}

	client := urlfetch.Client(ctx)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/get_current_account", auth.RequestURI), nil)
	if err != nil {
		return user{}, fmt.Errorf("GetUserV2 RequestGetUser Error: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+auth.Token)

	res, err := client.Do(req)
	if err != nil {
		return user{}, fmt.Errorf("GetUserV2 ClientGetUser Error: %v", err)
	}

	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return user{}, fmt.Errorf("GetUserV2 DecodeUser Error: %v", err)
	}
	// Return user data
	data := user{
		ID:       u.Account,
		Name:     u.Name.DisplayName,
		Username: u.Name.GivenName,
		Email:    u.Email,
		Avatar:   u.Avatar,
	}
	return data, nil
}

func (auth *oAuthV4) GetUser(ctx context.Context) (user, error) {
	var u struct {
		Account     string `json:"id_str"`
		GivenName   string `json:"screen_name"`
		DisplayName string `json:"name"`
		Email       string
		Avatar      string `json:"profile_image_url"`
		Country     string `json:"location"`
	}

	switch {
	case auth.RequestURI == "":
		return user{}, fmt.Errorf("GetUser Error: oAuthV4 RequestURI undefined, you need to define it before use oAuthV4 requests")
	case auth.UserID == "":
		return user{}, fmt.Errorf("GetUser Error: oAuthV4 UserID undefined, you need to define it before use oAuthV4 requests")
	case auth.ScreenName == "":
		return user{}, fmt.Errorf("GetUser Error: oAuthV4 ScreenName undefined, you need to define it before use oAuthV4 requests")
	}
	client := urlfetch.Client(ctx)
	params := map[string]string{"include_email": "true"}
	sign, values := EncodeSignature("GET", fmt.Sprintf("%s/account/verify_credentials.json", auth.RequestURI), auth, params)
	res, err := client.Get(EncodeNewHeaderHTTP(fmt.Sprintf("%s/account/verify_credentials.json", auth.RequestURI), sign, values))
	if err != nil {
		return user{}, fmt.Errorf("GetUser CLIENT Error: %v", err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		return user{}, fmt.Errorf("GetUser DECODE Error: %v", err)
	}
	// Return user data
	data := user{
		ID:       u.Account,
		Name:     u.DisplayName,
		Username: u.GivenName,
		Email:    u.Email,
		Avatar:   u.Avatar,
	}
	return data, nil
}
