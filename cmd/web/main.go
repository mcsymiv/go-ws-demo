package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mcsymiv/go-websocket/internal/handlers"
)

func main() {
	log.Println("Starting web app on port :8098")
	log.Println("Starting ws channel listener")
  go handlers.ListenWsChannel()

	mux := routes()

	srv := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("localhost:8098"),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Unable to start web app", err)
	}
}
