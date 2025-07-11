package main

import (
	"fmt"
	"log"
	"net/http"

	"hackathon-2025/internal/database"
	"hackathon-2025/pkg/handlers"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	port := ":8080"

	if err := database.Init(); err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
		log.Println("Server will start without database connection")
	}

	http.HandleFunc("/", corsMiddleware(handlers.RootHandler))
	http.HandleFunc("/hello", corsMiddleware(handlers.HelloHandler))
	http.HandleFunc("/health", corsMiddleware(handlers.HealthHandler))
	http.HandleFunc("/pages/user", corsMiddleware(handlers.GetPagesByUserHandler))
	http.HandleFunc("/get_profile_summary", corsMiddleware(handlers.GetProfileSummaryHandler))

	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET /                    - Welcome message")
	fmt.Println("  GET /hello               - Hello World message")
	fmt.Println("  GET /health              - Health check")
	fmt.Println("  GET /pages/user          - Get pages by user (requires email parameter)")
	fmt.Println("  GET /get_profile_summary - Get profile summary (requires email parameter)")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop the server")

	log.Fatal(http.ListenAndServe(port, nil))
}
