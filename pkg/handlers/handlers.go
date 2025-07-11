package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"hackathon-2025/internal/database"
	"hackathon-2025/pkg/models"
	"hackathon-2025/pkg/services"
)

// Response represents the API response structure
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// PagesResponse represents the response for pages endpoint
type PagesResponse struct {
	Pages     []models.PageInfo `json:"pages"`
	Count     int               `json:"count"`
	Timestamp time.Time         `json:"timestamp"`
	Status    string            `json:"status"`
}

// HelloHandler handles the /hello endpoint
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:   "Hello, World!",
		Timestamp: time.Now(),
		Status:    "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HealthHandler handles the /health endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	dbStatus := "healthy"
	if !database.IsConnected() {
		dbStatus = "unhealthy"
	}

	response := Response{
		Message:   "Service is " + dbStatus,
		Timestamp: time.Now(),
		Status:    dbStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// RootHandler handles the root endpoint
func RootHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:   "Welcome to Go HTTP API",
		Timestamp: time.Now(),
		Status:    "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetPagesByUserHandler handles the /pages/user endpoint
func GetPagesByUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get email from query parameter
	email := r.URL.Query().Get("contributor")
	if email == "" {
		http.Error(w, "email parameter is required", http.StatusBadRequest)
		return
	}

	// Create Confluence service
	confluenceService := services.NewConfluenceService()

	// Get pages by user
	pages, err := confluenceService.GetPagesByUserWithContent(email)
	if err != nil {
		log.Printf("Error getting pages for user %s: %v", email, err)
		http.Error(w, "Failed to retrieve pages", http.StatusInternalServerError)
		return
	}

	// Create response
	response := PagesResponse{
		Pages:     pages,
		Count:     len(pages),
		Timestamp: time.Now(),
		Status:    "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
