package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func serveDefault() *http.Server {
	fmt.Println(`
    .___      _____             .__   __
  __| _/_____/ ____\____   __ __|  |_/  |_
 / __ |/ __ \   __\\__  \ |  |  \  |\   __\
/ /_/ \  ___/|  |   / __ \|  |  /  |_|  |
\____ |\___  >__|  (____  /____/|____/__|
     \/    \/           \/
Go standard package.
https://golang.org/pkg/net/http/
____________________________________O/_______
																		O\       `)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		mockProcessing()
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.Write(message)
	})

	s := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	if isHTTPS() {
		go s.ListenAndServeTLS(tls.cert, tls.key)
	} else {
		go s.ListenAndServe()
	}
	return s
}

func shutdownDefault(s *http.Server) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	s.Shutdown(ctx)
}
