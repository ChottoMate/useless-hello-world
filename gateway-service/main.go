package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
)

func main() {
	xray.Configure(xray.Config{
		LogLevel: "info",
	})

	http.Handle("/", xray.Handler(xray.NewFixedSegmentNamer("gateway-service"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Gateway received a request")

		client := xray.Client(nil)

		helloURL := os.Getenv("HELLO_URL")
		if helloURL == "" {
			helloURL = "http://localhost:8080/hello"
		}
		
		req, err := http.NewRequest(http.MethodGet, helloURL, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating Hello request: %v", err), 500)
			return
		}
		req = req.WithContext(r.Context())

		helloResp, err := client.Do(req)
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

		req, err = http.NewRequest(http.MethodGet, worldURL, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating World request: %v", err), 500)
			return
		}
		req = req.WithContext(r.Context())

		worldResp, err := client.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error calling World: %v", err), 500)
			return
		}
		defer worldResp.Body.Close()
		worldBody, _ := io.ReadAll(worldResp.Body)

		result := fmt.Sprintf("%s %s", string(helloBody), string(worldBody))
		fmt.Fprintln(w, result)
	})))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway started on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
