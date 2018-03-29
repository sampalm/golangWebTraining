package main

import (
	"io"
	"net/http"
	"strings"
)

type foxHandler int

func (f foxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	sl := strings.Split(r.URL.Path, "/")
	animal := sl[1]

	io.WriteString(w, "<p><strong>"+animal+"</strong> is a clever animal.</p>\n <img width='500px' src='http://www.pbs.org/wnet/nature/files/2017/09/x1WLcZn-asset-mezzanine-16x9-6kkb4dA.jpg'>")
}

type hawkHandler int

func (h hawkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	sl := strings.Split(r.RequestURI, "/")
	animal := sl[1]

	io.WriteString(w, "<p><strong>"+animal+"</strong> can fly very high.</p>\n <img width='500px' src='https://upload.wikimedia.org/wikipedia/commons/8/8e/Parabuteo_unicinctus_-_01.jpg'>")
}

func owl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	sl := strings.Split(r.RequestURI, "/")
	animal := sl[1]

	io.WriteString(w, `<p><strong>`+animal+`</strong> is a beautiful animal.</p><img width="500px" src="https://steemitimages.com/DQmek2PfaZ5aWwQgscMMvS9nwnKTxgsNZ63s6hyfDneRtHK/20170905_233853_rmscr.jpg">`)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	io.WriteString(w, `
		<h1>Meet an animal:</h1>
		<ul>
			<li><h2><a href="/fox">Fox</a></h2></li>
			<li><h2><a href="/hawk">Hawk</a></h2></li>
			<li><h2><a href="/owl">Owl</a></h2></li>						
		</ul>
	`)
}

func main() {
	var fox foxHandler
	var hawk hawkHandler

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.Handle("/fox/", fox)
	mux.Handle("/hawk/", hawk)
	mux.HandleFunc("/owl/", owl)

	http.ListenAndServe(":8080", mux)
}
