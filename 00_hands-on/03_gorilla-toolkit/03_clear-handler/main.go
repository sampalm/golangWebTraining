package webapp

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Index page </h1>")
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Login page </h1>")
}

func init() {
	mux := http.DefaultServeMux

	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)
	mux.Handle("favicon.ico", http.NotFoundHandler())

	http.Handle("/", context.ClearHandler(mux))
}
