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

	http.Handle("/hello", xray.Handler(xray.NewFixedSegmentNamer("hello-service"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request: /hello")
		fmt.Fprintf(w, "Hello")
	})))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Hello-Service started on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
