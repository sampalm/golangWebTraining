package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	res, err := http.Get(os.Args[1])
	if err != nil {
		fmt.Println("ERR: ", err)
		return
	}
	defer res.Body.Close()
	body := html.NewTokenizer(res.Body)
	// Loop over throught HTML tokens
	for {
		token := body.Next()
		if token == html.ErrorToken {
			//fmt.Println("ErrorToken")
			break // End of the document
		}
		if token != html.StartTagToken {
			//fmt.Println("WrongToken")
			continue // Ignore that token
		}
		// Check if token is an anchor
		node := body.Token()
		isAnchor := node.Data == "a"
		if isAnchor {
			// Prints out every link found
			for _, a := range node.Attr {
				if a.Key == "href" {
					url := a.Val
					if strings.Index(url, "http") == 0 {
						fmt.Println("URL FOUND: ", a.Val)
					}
				}
			}
		}
	}
}
