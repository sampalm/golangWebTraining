package entrychecker

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

var tpl *template.Template

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/api/check", check)
	http.Handle("favicon.ico", http.NotFoundHandler())

	// index/public/public -> /public = index/public
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	tpl = template.Must(template.ParseGlob("*.html"))
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func check(w http.ResponseWriter, req *http.Request) {
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "Server received: "+string(bs))
}
