package main

import (
	"fmt"
	"log"
	"net/http"

	"hackathon-2025/pkg/handlers"
)

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

func main() {
	// Define port
	port := ":8080"

	// Set up routes with CORS middleware
	http.HandleFunc("/", corsMiddleware(handlers.RootHandler))
	http.HandleFunc("/hello", corsMiddleware(handlers.HelloHandler))
	http.HandleFunc("/health", corsMiddleware(handlers.HealthHandler))

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
