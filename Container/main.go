package main

import (
	"fmt"
	"log"
	"net/http"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Handle the incoming request
	if r.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
		return
	}

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "only get requests are accepted")
		return
	}

	fmt.Fprintf(w, `I have decided not to release the web scraping code publicly. Web scraping operates in a legal gray area, 
				    and the primary purpose of this project was to utilize AWS Fargate services, not potentially engage in unlawful activities.`)

	w.WriteHeader(http.StatusOK)
}

func main() {

	fmt.Println("Starting the program")

	port := ":80"
	server := &http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(HandleRequest),
	}

	// Start the server
	fmt.Printf("Server listening on port %s\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
