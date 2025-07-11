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
	http.HandleFunc("/user/post", corsMiddleware(handlers.UserPostHandler))
	
	// User Profile CRUD endpoints
	http.HandleFunc("/user/profile", corsMiddleware(handlers.CreateUserProfileHandler))
	http.HandleFunc("/user/profile/get", corsMiddleware(handlers.GetUserProfileHandler))
	http.HandleFunc("/user/profile/update", corsMiddleware(handlers.UpdateUserProfileHandler))
	http.HandleFunc("/user/profile/delete", corsMiddleware(handlers.DeleteUserProfileHandler))
	http.HandleFunc("/user/profile/list", corsMiddleware(handlers.ListUserProfilesHandler))

	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET /                    - Welcome message")
	fmt.Println("  GET /hello               - Hello World message")
	fmt.Println("  GET /health              - Health check")
	fmt.Println("  GET /pages/user          - Get pages by user (requires email parameter)")
	fmt.Println("  GET /get_profile_summary - Get profile summary (requires email parameter)")
	fmt.Println("  POST/PUT/GET/DELETE /user/post - CRUD operations for user posts")
	fmt.Println("")
	fmt.Println("User Profile API Usage:")
	fmt.Println("  POST /user/profile       - Create a new user profile")
	fmt.Println("  GET /user/profile/get?email=<email> - Get a specific user profile")
	fmt.Println("  PUT /user/profile/update - Update a user profile")
	fmt.Println("  DELETE /user/profile/delete?email=<email> - Delete a user profile")
	fmt.Println("  GET /user/profile/list   - List all user profiles (supports limit, offset)")
	fmt.Println("")
	fmt.Println("User Post API Usage:")
	fmt.Println("  POST /user/post?op_type=create - Create a new post")
	fmt.Println("  GET /user/post?op_type=read&post_id=<id> - Read a specific post")
	fmt.Println("  PUT /user/post?op_type=update&post_id=<id> - Update a post")
	fmt.Println("  DELETE /user/post?op_type=delete&post_id=<id> - Delete a post")
	fmt.Println("  GET /user/post?op_type=list - List all posts (supports limit, offset, author_id)")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop the server")

	log.Fatal(http.ListenAndServe(port, nil))
}
