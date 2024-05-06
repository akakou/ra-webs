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

			fmt.Fprintf(w, "<script>location.href = '%v'</script>", config.TTP+REDIRECT_PATH)
		}
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
	server.ListenAndServeTLS("", "")
}
