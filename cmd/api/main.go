package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nullzone-test/api-gateway/internal/gateway"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	gw := gateway.New()
	gw.Use(gateway.Logger())
	gw.Use(gateway.RateLimit(100))

	gw.Handle("GET", "/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	log.Printf("api-gateway listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, gw))
}
