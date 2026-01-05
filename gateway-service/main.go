package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
)

func main() {
	xray.Configure(xray.Config{
		LogLevel: "info",
	})

	http.Handle("/", xray.Handler(xray.NewFixedSegmentNamer("gateway-service"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Gateway received a request")

		// Configure HttpClient with timeout
		httpClient := &http.Client{
			Timeout: 5 * time.Second,
		}
		client := xray.Client(httpClient)

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

		helloResp, err := doWithRetry(client, req)
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

		worldResp, err := doWithRetry(client, req)
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

func doWithRetry(client *http.Client, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < 3; i++ {
		if i > 0 {
			// Use select to allow context cancellation during backoff
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(200 * time.Millisecond):
				log.Printf("Retrying request to %s (attempt %d)", req.URL, i+1)
			}
		}
		
		// Create a new request based on the original one if needed?
		// But for GET requests without body, reuse is mostly fine.
		// However, context cancellation should be respected. 
		// req is already bound to ctx from r.Context() in main.
		
		resp, err = client.Do(req)
		
		// If success and status is OK (less than 500)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// If we got a response but it was 5xx, ensure we close body before retrying
		if resp != nil {
			resp.Body.Close()
		}
	}
	// Return the last error or response (which might be 5xx)
	return resp, err
}
