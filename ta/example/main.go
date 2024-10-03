package main

import (
	"fmt"
	"net/http"

	"github.com/akakou/ra_webs/ta"
)

const REDIRECT_PATH = "/app/verification-status/"

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
			fmt.Fprintln(w, `<script>`)
			for _, v := range config.Verifiers {
				fmt.Fprintf(w, `open('%v')`, v+REDIRECT_PATH)
			}
			fmt.Fprintln(w, `</script>`)
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
