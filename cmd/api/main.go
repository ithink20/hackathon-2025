package main

import (
	"fmt"
	"log"
	"net/http"

	"hackathon-2025/pkg/handlers"
)

func main() {
	// Define port
	port := ":8080"

	// Set up routes
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/hello", handlers.HelloHandler)
	http.HandleFunc("/health", handlers.HealthHandler)

	// Start server
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET /         - Welcome message")
	fmt.Println("  GET /hello    - Hello World message")
	fmt.Println("  GET /health   - Health check")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop the server")

	log.Fatal(http.ListenAndServe(port, nil))
}
