package main

import (
	"fmt"
	"net/http"

	"github.com/akakou/ra_webs/ta"
)

const REDIRECT_PATH = "/app/redirect/"

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
		_, err := r.Cookie("isFirstAccess")

		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:  "isFirstAccess",
				Value: "true",
			})

			fmt.Fprintf(w, `
				We will redirect verifier after 3 second....
				<script>
					setTimeout(() => {
						location.href = '%v'
					}, 3000)
				</script>
				`, config.Verifier+REDIRECT_PATH)
		}

		fmt.Fprintln(w, "Hello from TA running on TEE :)")
	}

	tlsConfig, err := ta.TLSConfig()
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr:      "10.0.0.5:443:443",
		Handler:   nil,
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", handler)
	server.ListenAndServeTLS("", "")
}
