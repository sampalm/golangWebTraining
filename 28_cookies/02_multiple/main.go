package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", set)
	http.HandleFunc("/read", read)
	http.HandleFunc("/abundance", abundance)
	http.Handle("favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func set(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "another-cookie",
		Value: "Something here.",
	})

	fmt.Fprintln(w, "COOKIE WRITTEN - CHECK YOUR BROWSER")
}

func read(w http.ResponseWriter, req *http.Request) {
	ck, err := req.Cookie("another-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		fmt.Fprintln(w, "YOUR COOKIE: ", ck)
	}

	ck2, err := req.Cookie("general")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Fprintln(w, "YOUR COOKIE: ", ck2)
	}

	ck3, err := req.Cookie("default")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Fprintln(w, "YOUR COOKIE: ", ck3)
	}
}

func abundance(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "general",
		Value: "Another cookie",
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "default",
		Value: "Default cookie",
	})

	fmt.Fprintln(w, "COOKIES WRITTEN - CHECK YOUR BROWSER")
}
