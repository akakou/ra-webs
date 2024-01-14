package main

import (
	"log"
	"net/http"

	"github.com/akakou/ra_webs/ta"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8081/redirect", http.StatusTemporaryRedirect)
}

func main() {
	config := ta.RAConfig{
		TTPDomain: "",
		Domain:    "",
		Email:     "",
	}

	ra := ta.NewRA(&config)
	_, err := ra.TLSConfig()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", RedirectHandler)
	http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
