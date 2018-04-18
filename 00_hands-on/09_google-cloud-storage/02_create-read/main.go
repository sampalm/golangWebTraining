package webgcloud

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
)

type File struct {
	Name        string
	ContentType string
	Bucket      string
	Link        string
	Size        int64
}

func init() {
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/", upload)

	http.Handle("/", router)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	io.WriteString(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>DemoGCS</title>
		</head>
		<body>
			<form method="POST" enctype="multipart/form-data">
				<input type="file" name="file" accept=".jpg, .jpeg, .png">
				<button type="submit">Upload</button>
			</form>
			<br><br>
	`)

	if files := getAllFromBucket(r); files != nil {
		bd := "<div>"
		for _, file := range files {
			bd += fmt.Sprintf(`
				<span>
					<p>Name: %s</p>
					<p>ContentType: %s</p>
					<p>Bucket: %s</p>
					<img src="%v" width="200px">
					<p>Size: %v</p>
				</span>	
			`, file.Name, file.ContentType, file.Bucket, file.Link, file.Size)
		}
		bd += "</div>"
		io.WriteString(w, bd)
	}

	io.WriteString(w, `
		</body>
		</html>	
	`)
}

func upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)
	buffer := bytes.NewBuffer(nil)

	// Get file from form
	f, header, err := r.FormFile("file")
	if err != nil {
		log.Errorf(ctx, "failed to get file: %v", err)
		http.Error(w, err.Error(), 500) // Internal Server Error
		return
	}
	filename := header.Filename
	filetype := header.Header.Get("Content-type")

	// Copy file to buffer
	if _, err := io.Copy(buffer, f); err != nil {
		log.Errorf(ctx, "failed to copy file: %v", err)
		http.Error(w, err.Error(), 500) // Internal Server Error
		return
	}

	// Get default bucket from GCS
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default bucket: %v", err)
		http.Error(w, err.Error(), 500) // Internal Server Error
		return
	}

	// Get new client and close it
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get new client: %v", err)
		http.Error(w, err.Error(), 500) // Internal Server Error
		return
	}
	defer client.Close()

	// Creates a file in Google Cloud Storage
	cb := client.Bucket(bucket)

	// Creates a NewWrite from  cb
	nw := cb.Object(filename).NewWriter(ctx)
	nw.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	nw.ContentType = filetype
	defer nw.Close()

	// Entries are immutable, be aggressive about caching (1 day).
	nw.CacheControl = "public, max-age=86400"

	if _, err := nw.Write(buffer.Bytes()); err != nil {
		log.Errorf(ctx, "createFile: unable to write data to bucket %q, file %q: %v", bucket, filename, err)
		http.Error(w, err.Error(), 500) // Internal Server Error
		return
	}

	if err := nw.Close(); err != nil {
		log.Errorf(ctx, "createFile: unable to close bucket %q, file %q: %v", bucket, filename, err)
		http.Error(w, err.Error(), 500) // Internal Server Error
		return
	}
}

func getAllFromBucket(r *http.Request) []File {
	var files []File
	ctx := appengine.NewContext(r)

	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get new client: %v", err)
		//http.Error(w, err.Error(), 500) // Internal Server Error
		return nil
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get new client: %v", err)
		//http.Error(w, err.Error(), 500) // Internal Server Error
		return nil
	}
	defer client.Close()

	cb := client.Bucket(bucket)
	query := &storage.Query{}

	itn := cb.Objects(ctx, query)
	for {
		obj, err := itn.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "listBucket: unable to list bucked: %v", err)
			return nil
		}

		newFile := File{
			Name:        obj.Name,
			ContentType: obj.ContentType,
			Bucket:      obj.Bucket,
			Link:        obj.MediaLink,
			Size:        obj.Size,
		}
		files = append(files, newFile)
	}

	return files
}
