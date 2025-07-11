package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"hackathon-2025/internal/database"
	"hackathon-2025/pkg/models"
	"hackathon-2025/pkg/services"

	"gorm.io/gorm"
)

type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

type PagesResponse struct {
	Pages     []models.PageInfo `json:"pages"`
	Count     int               `json:"count"`
	Timestamp time.Time         `json:"timestamp"`
	Status    string            `json:"status"`
}

type ProfileSummaryResponse struct {
	Data struct {
		Outputs interface{} `json:"outputs"`
		Status  string      `json:"status"`
	} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	Error     interface{} `json:"error,omitempty"`
}

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

func HealthHandler(w http.ResponseWriter, r *http.Request) {
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

func GetPagesByUserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("contributor")
	if email == "" {
		http.Error(w, "contributor parameter is required", http.StatusBadRequest)
		return
	}

	sync := r.URL.Query().Get("sync") == "true"

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	var userPages []models.UserPage
	var err error

	if sync {
		userPages, err = syncUserPagesFromConfluence(db, email)
	} else {
		err = db.Where("user_email = ? AND deleted_at IS NULL", email).Find(&userPages).Error
	}

	if err != nil {
		log.Printf("Error getting pages for user %s: %v", email, err)
		http.Error(w, "Failed to retrieve pages", http.StatusInternalServerError)
		return
	}

	var pages []models.PageInfo
	for _, up := range userPages {
		pages = append(pages, models.PageInfo{
			ID:      up.PageID,
			Type:    up.PageType,
			Title:   up.PageTitle,
			Content: up.PageContent,
			Link:    up.PageLink,
		})
	}

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

func GetProfileSummaryHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email parameter is required", http.StatusBadRequest)
		return
	}

	sync := r.URL.Query().Get("sync") == "true"

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	if !sync {
		// Check if user profile already exists in database
		var existingProfile models.UserProfile
		err := db.Where("user_email = ? AND deleted_at IS NULL", email).First(&existingProfile).Error

		if err == nil {
			// User profile exists, return from database
			var outputs interface{}
			if existingProfile.AISummary != "" {
				if err := json.Unmarshal([]byte(existingProfile.AISummary), &outputs); err != nil {
					log.Printf("Error unmarshaling stored AI summary for %s: %v", email, err)
					// If unmarshaling fails, continue to generate new summary
				} else {
					// Successfully retrieved from database
					response := ProfileSummaryResponse{
						Data: struct {
							Outputs interface{} `json:"outputs"`
							Status  string      `json:"status"`
						}{
							Outputs: outputs,
							Status:  "success",
						},
						Timestamp: time.Now(),
						Error:     nil,
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)

					if err := json.NewEncoder(w).Encode(response); err != nil {
						log.Printf("Error encoding response: %v", err)
						http.Error(w, "Internal server error", http.StatusInternalServerError)
					}
					return
				}
			}
		}
	}

	// If no existing profile or unmarshaling failed, generate new summary
	// 1. Get documents from GetPagesByUserHandler
	var userPages []models.UserPage
	err := db.Where("user_email = ? AND deleted_at IS NULL", email).Find(&userPages).Error
	if err != nil {
		log.Printf("Error getting pages for user %s: %v", email, err)
		http.Error(w, "Failed to retrieve pages", http.StatusInternalServerError)
		return
	}

	// Convert pages to documents string
	var documents string
	for i, up := range userPages {
		if i > 0 {
			documents += "\n\n"
		}
		// Truncate individual page content if too long
		content := up.PageContent
		if len(content) > 1000 {
			content = content[:1000] + "... [truncated]"
		}
		documents += fmt.Sprintf("Title: %s\nType: %s\nContent: %s", up.PageTitle, up.PageType, content)
	}

	// 2. Get template content by type
	templateService := services.NewTemplateService()
	template, err := templateService.GetTemplateContentByType("profile_summary")
	if err != nil {
		log.Printf("Error getting template for profile_summary: %v", err)
		http.Error(w, "Failed to retrieve template", http.StatusInternalServerError)
		return
	}

	// 3. Call ProfileSummaryAgent with payload
	apiKey := "W1fmYCsR1S6va1esyYSTwrFC14KELW4J"
	profileSummaryAgent := services.ProfileSummaryAgent(apiKey)

	profileResponse, err := profileSummaryAgent.RunProfileSummary(documents, template, email)
	if err != nil {
		log.Printf("Error calling ProfileSummaryAgent: %v", err)
		http.Error(w, "Failed to generate profile summary", http.StatusInternalServerError)
		return
	}

	response := ProfileSummaryResponse{
		Data: struct {
			Outputs interface{} `json:"outputs"`
			Status  string      `json:"status"`
		}{
			Outputs: profileResponse.Data.Outputs,
			Status:  profileResponse.Data.Status,
		},
		Timestamp: time.Now(),
		Error:     profileResponse.Error,
	}

	// Save AI summary to database asynchronously
	go func() {
		if profileResponse.Error == nil && profileResponse.Data.Outputs != nil {
			// Convert outputs to JSON string for storage
			outputsJSON, err := json.Marshal(profileResponse.Data.Outputs)
			if err != nil {
				log.Printf("Error marshaling outputs to JSON: %v", err)
				return
			}

			userProfile := models.UserProfile{
				UserEmail: email,
				AISummary: string(outputsJSON),
			}

			// Upsert the user profile
			result := db.Where("user_email = ? AND deleted_at IS NULL", email).
				Assign(userProfile).
				FirstOrCreate(&userProfile)

			if result.Error != nil {
				log.Printf("Error saving user profile for %s: %v", email, result.Error)
			} else {
				log.Printf("Successfully saved AI summary for user %s", email)
			}
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func syncUserPagesFromConfluence(db *gorm.DB, email string) ([]models.UserPage, error) {
	confluenceService := services.NewConfluenceService()

	pages, err := confluenceService.GetPagesByUserWithContent(email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Confluence: %w", err)
	}

	var userPages []models.UserPage

	for _, page := range pages {
		userPage := models.UserPage{
			UserEmail:   email,
			PageID:      page.ID,
			PageType:    page.Type,
			PageTitle:   page.Title,
			PageContent: page.Content,
			PageLink:    fmt.Sprintf("https://confluence.shopee.io/pages/viewpage.action?pageId=%s", page.ID),
		}

		result := db.Where("page_id = ?", page.ID).
			Assign(userPage).
			FirstOrCreate(&userPage)

		if result.Error != nil {
			log.Printf("Error upserting page %s: %v", page.ID, result.Error)
			continue
		}

		userPages = append(userPages, userPage)
	}

	return userPages, nil
}
