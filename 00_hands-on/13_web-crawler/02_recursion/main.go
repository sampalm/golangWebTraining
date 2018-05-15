package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var urls []string

func main() {
	res, err := http.Get(os.Args[1])
	if err != nil {
		fmt.Println("ERR: ", err)
		return
	}
	defer res.Body.Close()
	node, err := html.Parse(res.Body)
	if err != nil {
		fmt.Println("ERR: ", err)
		return
	}
	visit(node)

	for _, url := range urls {
		fmt.Println("URL FOUND: ", url)
	}
}

func visit(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				url := a.Val
				if strings.Index(url, "http") == 0 {
					urls = append(urls, url)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visit(c)
	}
}
