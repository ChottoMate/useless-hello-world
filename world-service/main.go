package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
)

func main() {
	xray.Configure(xray.Config{
		LogLevel: "info",
	})

	http.Handle("/world", xray.Handler(xray.NewFixedSegmentNamer("world-service"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request: /world")
		fmt.Fprintf(w, "World")
	})))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("World-Service started on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
