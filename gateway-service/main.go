package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Gateway received a request")

		helloURL := os.Getenv("HELLO_URL")
		if helloURL == "" {
			helloURL = "http://localhost:8080/hello"
		}
		helloResp, err := http.Get(helloURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error calling Hello: %v", err), 500)
			return
		}
		defer helloResp.Body.Close()
		helloBody, _ := io.ReadAll(helloResp.Body)

		worldURL := os.Getenv("WORLD_URL")
		if worldURL == "" {
			worldURL = "http://localhost:8081/world"
		}
		worldResp, err := http.Get(worldURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error calling World: %v", err), 500)
			return
		}
		defer worldResp.Body.Close()
		worldBody, _ := io.ReadAll(worldResp.Body)

		result := fmt.Sprintf("%s %s", string(helloBody), string(worldBody))
		fmt.Fprintln(w, result)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway started on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
