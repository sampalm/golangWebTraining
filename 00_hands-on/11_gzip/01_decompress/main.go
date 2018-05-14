package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
)

var file = "https://dl.opensubtitles.org/en/download/src-api/vrf-19de0c5c/filead/1955648481.gz"

func main() {
	res, err := http.Get(file)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	rd, err := gzip.NewReader(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer rd.Close()
	io.Copy(os.Stdout, rd)
}
