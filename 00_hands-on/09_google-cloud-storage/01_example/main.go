package webgcloud

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/up", upSomething)
	http.Handle("favicon.ico", http.NotFoundHandler())
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	io.WriteString(w, "Index Page!")
}

func upSomething(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// GET DEFAULT BUCKET
	bck, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default bucket: %v", err)
	}

	// initialize a new client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to create client: %v", err)
		return
	}
	// close client at the end
	defer client.Close()

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Demo GCS Application running from Version: %v\n", appengine.VersionID(ctx))
	fmt.Fprintf(w, "Using bucket name: %v\n", bck)

	// Creates a file in Google Cloud Storage
	clientBck := client.Bucket(bck)
	filename := "demo-test-file"
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "Creating file /%v/%v\n", bck, filename)

	clientW := clientBck.Object(filename).NewWriter(ctx)
	clientW.ContentType = "text/plain"
	clientW.Metadata = map[string]string{
		"x-goog-meta-foo": "foo",
		"x-goog-meta-bar": "bar",
	}
	defer clientW.Close()

	if _, err := clientW.Write([]byte("Some text here just to test.\n")); err != nil {
		log.Errorf(ctx, "createFile: unable to write data to bucket %q, file %q: %v", bck, filename, err)
		return
	}
	if err := clientW.Close(); err != nil {
		log.Errorf(ctx, "createFile: unable to close bucket %q, file %q: %v", bck, filename, err)
		return
	}
}
