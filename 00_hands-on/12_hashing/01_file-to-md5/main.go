package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./Supernatural.S13E21.720p.HDTV.x264-SVA[eztv].mkv")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	st, _ := file.Stat()
	hash := md5.New()
	io.Copy(hash, file)
	fmt.Printf("Hash found: %x\n", hash.Sum(nil))
	fmt.Println("Size file: ", st.Size())
}
