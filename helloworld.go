package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("helloworld: received a request")

	// Log all of the headers in the request
	for name, values := range r.Header {
		for _, value := range values {
			log.Println(name, value)
		}
	}

	// Respond w/ JSON for convienence
	json.NewEncoder(w).Encode(r.Header)
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
