package main

import (
	"log"
	"net/http"

	"github.com/JyotinderSingh/trytok/pkg/handler"
)

func main() {
	http.HandleFunc("/execute", handler.ExecuteCode)
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
