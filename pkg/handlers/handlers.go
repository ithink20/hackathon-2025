package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
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

type Contribution struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Documents   []string `json:"documents"`
}

type ProcessedProfileSummary struct {
	Role                string         `json:"role"`
	Team                string         `json:"team"`
	Tags                []string       `json:"tags"`
	Summary             string         `json:"summary"`
	RecentContributions []Contribution `json:"recentContributions"`
}

type ProcessedProfileSummaryResponse struct {
	Data      ProcessedProfileSummary `json:"data"`
	Timestamp time.Time               `json:"timestamp"`
	Error     interface{}             `json:"error,omitempty"`
}

type ProfileSummaryResponse struct {
	Data struct {
		Outputs interface{} `json:"outputs"`
		Status  string      `json:"status"`
	} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	Error     interface{} `json:"error,omitempty"`
}

// UserProfileResponse represents the response structure for user profile operations
type UserProfileResponse struct {
	Data      models.UserProfile `json:"data"`
	Message   string             `json:"message"`
	Timestamp time.Time          `json:"timestamp"`
	Status    string             `json:"status"`
}

// UserProfileListResponse represents the response structure for listing user profiles
type UserProfileListResponse struct {
	Data      []models.UserProfile `json:"data"`
	Count     int                  `json:"count"`
	Message   string               `json:"message"`
	Timestamp time.Time            `json:"timestamp"`
	Status    string               `json:"status"`
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
		err = db.Where("user_email = ? AND deleted_at IS NULL", email).Order("page_id DESC").Find(&userPages).Error
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
					// Process the stored data using the same parsing logic
					processedData := processProfileResponse(outputs)

					// Successfully retrieved from database
					response := ProcessedProfileSummaryResponse{
						Data:      processedData,
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
	err := db.Where("user_email = ? AND deleted_at IS NULL", email).Order("page_id DESC").Find(&userPages).Error
	if err != nil {
		log.Printf("Error getting pages for user %s: %v", email, err)
		http.Error(w, "Failed to retrieve pages", http.StatusInternalServerError)
		return
	}

	var limitPages []models.UserPage

	if len(userPages) == 0 {
		userPages, err = syncUserPagesFromConfluence(db, email)
	}

	// Ensure we don't exceed array bounds
	if len(userPages) > 5 {
		limitPages = userPages[:5]
	} else {
		limitPages = userPages
	}

	// Create comma-separated string of page links
	var pageIds []string
	for _, page := range limitPages {
		if page.PageID != "" {
			pageIds = append(pageIds, page.PageID)
		}
	}
	pageIdsStr := strings.Join(pageIds, ",")

	agentAns, agentErr := services.SmartAgentInvoke(pageIdsStr, services.SmartAgentRequest{
		EndpointDeploymentHashID: "2i3rtmg8ssjttba1724ynqde",
		EndpointDeploymentKey:    "yx4q46lnxo6kpyals96vm43c",
		UserID:                   "2i3rtmg8ssjttba1724ynqde",
	})

	if agentErr != nil {
		log.Printf("Error calling SmartAgentInvoke: %v, %v", agentErr, pageIdsStr)
		http.Error(w, "Failed to get AI response", http.StatusInternalServerError)
		return
	}

	log.Printf("data : %v", agentAns.Data.Response.ResponseStr)

	// Clean the response string by removing markdown code block markers and trimming whitespace
	cleanResponseStr := strings.TrimSpace(agentAns.Data.Response.ResponseStr)

	// Remove markdown code block markers
	cleanResponseStr = strings.TrimPrefix(cleanResponseStr, "```json")
	cleanResponseStr = strings.TrimPrefix(cleanResponseStr, "```")
	cleanResponseStr = strings.TrimSuffix(cleanResponseStr, "```")

	// Remove any remaining quotes and whitespace
	cleanResponseStr = strings.TrimSpace(cleanResponseStr)
	cleanResponseStr = strings.Trim(cleanResponseStr, "\"'")

	// Parse the AI response string into a map first
	var responseData map[string]interface{}
	if err := json.Unmarshal([]byte(cleanResponseStr), &responseData); err != nil {
		log.Printf("Error unmarshaling AI response: %v", err)
		log.Printf("Raw response: %s", agentAns.Data.Response.ResponseStr)
		log.Printf("Cleaned response: %s", cleanResponseStr)
		http.Error(w, "Failed to parse AI response", http.StatusInternalServerError)
		return
	}

	// Parse the AI response using the existing processProfileResponse function
	processedData := processProfileResponse(responseData)

	// Save the processed data to the UserProfile table
	// First, check if user profile exists
	var userProfile models.UserProfile
	err = db.Where("user_email = ? AND deleted_at IS NULL", email).First(&userProfile).Error

	if err != nil {
		// User profile doesn't exist, create a new one
		userProfile = models.UserProfile{
			UserEmail:  email,
			UserName:   email,                                                                                         // Default to email if no name provided
			ProfileImg: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face", // Default profile image
		}
	}

	// Convert the processed data back to JSON string for storage
	aiSummaryJSON, err := json.Marshal(responseData)
	if err != nil {
		log.Printf("Error marshaling AI summary for storage: %v", err)
		http.Error(w, "Failed to process AI summary", http.StatusInternalServerError)
		return
	}

	// Update the AISummary field
	userProfile.AISummary = string(aiSummaryJSON)

	// Save or update the user profile
	if userProfile.ID == 0 {
		// Create new profile
		if err := db.Create(&userProfile).Error; err != nil {
			log.Printf("Error creating user profile: %v", err)
			http.Error(w, "Failed to save user profile", http.StatusInternalServerError)
			return
		}
	} else {
		// Update existing profile
		if err := db.Model(&userProfile).Update("ai_summary", userProfile.AISummary).Error; err != nil {
			log.Printf("Error updating user profile: %v", err)
			http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
			return
		}
	}

	// Create the proper response structure
	response := ProcessedProfileSummaryResponse{
		Data:      processedData,
		Timestamp: time.Now(),
		Error:     nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func parseRecentContributions(contributionsStr string) []Contribution {
	var contributions []Contribution

	// Split by "|" delimiter
	parts := strings.Split(contributionsStr, "|")

	// Process in groups of 2 (key, value)
	var currentContribution Contribution
	var contributionCount int

	for i := 0; i < len(parts)-1; i += 2 {
		if i+1 >= len(parts) {
			break
		}

		key := strings.TrimSpace(parts[i])
		value := strings.TrimSpace(parts[i+1])

		// Skip if key or value is empty
		if key == "" || value == "" {
			continue
		}

		// Check if this is a new contribution (title1, title2, etc.)
		if strings.HasPrefix(key, "title") {
			// If we have a previous contribution, save it
			if currentContribution.Title != "" {
				contributions = append(contributions, currentContribution)
			}

			// Start new contribution
			currentContribution = Contribution{
				Title:       value,
				Description: "",
				Tags:        []string{},
				Documents:   []string{},
			}
			contributionCount++
		} else if strings.HasPrefix(key, "description") {
			currentContribution.Description = value
		} else if strings.HasPrefix(key, "tags") {
			// Parse tags (assuming they're comma-separated)
			if value != "" && value != "Unknown" {
				tagParts := strings.Split(value, ",")
				for _, tag := range tagParts {
					trimmedTag := strings.TrimSpace(tag)
					if trimmedTag != "" && trimmedTag != "Unknown" {
						currentContribution.Tags = append(currentContribution.Tags, trimmedTag)
					}
				}
			}

			// If no valid tags, provide default
			if len(currentContribution.Tags) == 0 {
				currentContribution.Tags = []string{"No tags"}
			}
		} else if strings.HasPrefix(key, "documents") {
			// Parse documents (assuming they're comma-separated or single value)
			if value != "" && value != "Unknown" {
				docParts := strings.Split(value, ",")
				for _, doc := range docParts {
					trimmedDoc := strings.TrimSpace(doc)
					if trimmedDoc != "" && trimmedDoc != "Unknown" {
						currentContribution.Documents = append(currentContribution.Documents, trimmedDoc)
					}
				}
			}
		}
	}

	// Append the last contribution if it exists
	if currentContribution.Title != "" {
		contributions = append(contributions, currentContribution)
	}

	return contributions
}

func processProfileResponse(rawOutputs interface{}) ProcessedProfileSummary {
	// Convert raw outputs to map for easier processing
	outputsMap, ok := rawOutputs.(map[string]interface{})
	if !ok {
		log.Printf("Error: rawOutputs is not a map")
		return ProcessedProfileSummary{}
	}

	processed := ProcessedProfileSummary{}

	// Extract basic fields with better validation
	if role, ok := outputsMap["role"].(string); ok && role != "" && role != "Unknown" {
		processed.Role = role
	} else {
		processed.Role = "Not Specified"
	}

	if team, ok := outputsMap["team"].(string); ok && team != "" {
		processed.Team = team
	} else {
		processed.Team = "Not Specified"
	}

	if summary, ok := outputsMap["summary"].(string); ok && summary != "" {
		processed.Summary = summary
	} else {
		processed.Summary = "No summary available"
	}

	// Parse tags with better validation
	if tagsStr, ok := outputsMap["tags"].(string); ok && tagsStr != "" {
		tagParts := strings.Split(tagsStr, ",")
		for _, tag := range tagParts {
			trimmedTag := strings.TrimSpace(tag)
			if trimmedTag != "" && trimmedTag != "Unknown" {
				processed.Tags = append(processed.Tags, trimmedTag)
			}
		}
	}

	// If no valid tags found, provide default
	if len(processed.Tags) == 0 {
		processed.Tags = []string{""}
	}

	// Parse recent contributions with better validation
	if contributionsStr, ok := outputsMap["recentContributions"].(string); ok && contributionsStr != "" {
		contributions := parseRecentContributions(contributionsStr)
		// Only add contributions that have valid titles
		for _, contribution := range contributions {
			if contribution.Title != "" && contribution.Title != "Unknown" {
				processed.RecentContributions = append(processed.RecentContributions, contribution)
			}
		}
	}

	// If no valid contributions found, provide a default message
	if len(processed.RecentContributions) == 0 {
		processed.RecentContributions = []Contribution{
			{
				Title:       "No Recent Contributions",
				Description: "No recent contributions found for this user",
				Tags:        []string{"No Data"},
				Documents:   []string{},
			},
		}
	}

	return processed
}

func syncUserPagesFromConfluence(db *gorm.DB, email string) ([]models.UserPage, error) {
	confluenceService := services.NewConfluenceService()

	pages, err := confluenceService.GetPagesByUser(email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Confluence: %w", err)
	}

	var userPages []models.UserPage

	for _, page := range pages {
		userPage := models.UserPage{
			UserEmail:     email,
			PageID:        page.ID,
			PageType:      page.Type,
			PageTitle:     page.Title,
			PageContent:   page.Content,
			PageLink:      page.Link,
			PageTimestamp: page.Timestamp,
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

	// Sort pages by timestamp in descending order
	sort.Slice(userPages, func(i, j int) bool {
		return userPages[i].PageTimestamp > userPages[j].PageTimestamp
	})

	return userPages, nil
}
