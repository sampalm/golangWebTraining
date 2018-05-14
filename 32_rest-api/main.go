package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string
		Type string
	}
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		//io.Copy(w, r.Body)
		json.NewDecoder(r.Body).Decode(&body)
		json.NewEncoder(w).Encode(body)
	}
}
