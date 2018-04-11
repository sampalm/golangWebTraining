package webpack

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("*.html"))

	http.HandleFunc("/", index)
	http.Handle("favicon.icon", http.NotFoundHandler())
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		serveTemplate(w, r, "index.html")
	} else {
		http.Error(w, "Unavailable service", http.StatusInternalServerError)
		return
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request, temp string) {
	c := appengine.NewContext(r)
	// Check if item already exists
	item, err := memcache.Get(c, temp)
	// If not create item
	if err != nil {
		// Create a buffer
		bf := new(bytes.Buffer)
		// Write template to both memcache and RW at same time
		mWriter := io.MultiWriter(w, bf)
		tpl.ExecuteTemplate(mWriter, temp, nil)
		// Set new item with the template
		memcache.Set(c, &memcache.Item{
			Key:   temp,
			Value: bf.Bytes(),
		})
		return
	}
	// If exists serve item
	// Serve template from memcache
	io.WriteString(w, string(item.Value))
}
