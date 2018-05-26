package spotifyapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine/urlfetch"
)

type API struct {
	ClientID     string
	ClientSecret string
	ClientToken  string
	Token
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires_in"`
}

func NewClient() *API {
	var api API
	api.ClientID = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	api.ClientSecret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	api.ClientToken = base64.StdEncoding.EncodeToString([]byte(api.ClientID + ":" + api.ClientSecret))
	return &api
}

func (api *API) GetAuthorization(ctx context.Context) error {
	client := urlfetch.Client(ctx)
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token?grant_type=client_credentials", nil)
	if err != nil {
		return fmt.Errorf("getAuthorization: request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+api.ClientToken)
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("getAuthorization: client: %s", err.Error())
	}
	defer res.Body.Close()
	var token Token
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return fmt.Errorf("getAuthorization: decoder: %s", err.Error())
	}
	api.Token = token
	return nil
}
func (api *API) GetTracks(ctx context.Context, trackID string) (TrackList, error) {
	client := urlfetch.Client(ctx)
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/tracks/?ids="+trackID, nil)
	if err != nil {
		return TrackList{}, fmt.Errorf("getAuthorization: request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+api.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		return TrackList{}, fmt.Errorf("getAuthorization: client: %s", err.Error())
	}
	defer res.Body.Close()
	var tracks TrackList
	err = json.NewDecoder(res.Body).Decode(&tracks)
	if err != nil {
		return TrackList{}, fmt.Errorf("getAuthorization: decoder: %s", err.Error())
	}
	return tracks, nil
}
