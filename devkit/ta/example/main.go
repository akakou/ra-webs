package main

import (
	"fmt"
	"net/http"

	// "github.com/akakou/ra-webs/devkit/core/attest/debug"
	"github.com/akakou/ra-webs/ta"
)

func main() {
	// debug.EnableDebug()

	config, err := ta.DefaultConfig()
	if err != nil {
		panic(err)
	}

	ta, err := ta.Init(config)
	if err != nil {
		panic(err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		fmt.Fprintln(w, "Hello from TA running on TEE :)<br/>")

		// for _, v := range config.CTR {
		// fmt.Fprintf(w, `<button onclick="window.open('%s');">Monitor Page (%s)</button><br/>`, v+VERIFIER_PATH, v)
		// }
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
