package main

import (
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handlers.WSEndpoint)

	log.Println("Starting requests queue listener")
	go handlers.ListenToRequestQueue()

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
