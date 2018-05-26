package spotifyapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

var client *API
var ctx context.Context

func init() {
	r := httprouter.New()
	r.GET("/tracks/:ids", tracks)
	http.Handle("/", r)
}

func tracks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; chaset=utf8")
	tracksID := p.ByName("ids")
	ctx = appengine.NewContext(r)
	client = NewClient()
	if err := client.GetAuthorization(ctx); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Infof(ctx, "CLIENT\nclient_token:%s\naccess_token:%s\ntoken:%s\n", client.ClientToken, client.AccessToken, client.Token)
	tracks, err := client.GetTracks(ctx, tracksID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(tracks)
}
