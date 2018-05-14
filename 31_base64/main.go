package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	s := "Douglas Eck is a scientist at Magenta, a Google AI project researching the use of machine learning to create music, video, images and text"
	std := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	s64 := base64.NewEncoding(std).EncodeToString([]byte(s))
	fmt.Println("DEFAULT: ", len(s))
	fmt.Println("ENCODED: ", len(s64))
	fmt.Println("DEFAULT: ", s)
	fmt.Println("ENCODED: ", s64)
}
