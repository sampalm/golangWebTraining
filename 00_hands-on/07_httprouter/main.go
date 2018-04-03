package webapp

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := r.URL.Query()
	if usname := q.Get("name"); usname != "" {
		tpl.Execute(w, usname)
		return
	}
	tpl.Execute(w, "Anonymous")
}

func fooWt(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Foo page without params.\n")
}

func foo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "Foo page, Hello %v", p.ByName("param"))
}

func gitCat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Octocat Git!</h1><br> <img src='../static/Octocat.png' alt='Octocat'> ")
}

func query(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryValues := r.URL.Query()
	fmt.Fprintf(w, "Params found in URL: \n Name: %s\n Age: %s", queryValues.Get("name"), queryValues.Get("age"))
}

func MyNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("My own Not Found Handler.\n"))
	w.Write([]byte("The page you requested could not be found."))
}

func MyNotFoundHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/notfound.html")
}

func init() {
	// Instantiate new router
	router := httprouter.New()

	// Parse template
	tpl = template.Must(template.ParseFiles("templates/index.html"))

	// Some Handlers e.g
	router.GET("/", index)
	router.GET("/foo/", fooWt)
	router.GET("/foo/:param", foo)
	router.GET("/octocat/", gitCat)

	// Query example: query/?name=Name&age=20
	router.GET("/query/", query)
	router.ServeFiles("/static/*filepath", http.Dir("static/"))

	// Custom NotFound Handler
	// #1 Without any template router.NotFound = http.HandlerFunc(MyNotFound)
	router.NotFound = http.HandlerFunc(MyNotFoundHTML)

	http.Handle("/", router)
}
