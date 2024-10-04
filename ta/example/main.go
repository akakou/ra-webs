package main

import (
	"fmt"
	"net/http"

	"github.com/akakou/ra_webs/ta"
)

const VERIFIER_PATH = "/app/verification-status/"

func main() {
	config, err := ta.DefaultConfig()
	if err != nil {
		panic(err)
	}

	ta, err := ta.InitTA(config)
	if err != nil {
		panic(err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		_, err := r.Cookie("isFirstAccess")

		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:  "isFirstAccess",
				Value: "true",
			})

			fmt.Fprintln(w, `We open verifier....<br/>`)
			for _, v := range config.Verifiers {
				fmt.Fprintf(w, `<button onclick="window.open('%s');">Verifier Page (%s)</button>`, v+VERIFIER_PATH, v)
			}
		}

		fmt.Fprintln(w, "Hello from TA running on TEE :)")
	}

	tlsConfig, err := ta.TLSConfig()
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr:      ":443",
		Handler:   nil,
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", handler)

	err = server.ListenAndServeTLS("", "")
	panic(err)
}
