package main

import (
	"fmt"
	"net/http"

	"github.com/akakou/ra-webs/ta"
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

		fmt.Fprintln(w, "Hello from TA running on TEE :)<br/>")

		for _, v := range config.Monitors {
			fmt.Fprintf(w, `<button onclick="window.open('%s');">Monitor Page (%s)</button><br/>`, v+VERIFIER_PATH, v)
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

	err = server.ListenAndServeTLS("", "")
	panic(err)
}
