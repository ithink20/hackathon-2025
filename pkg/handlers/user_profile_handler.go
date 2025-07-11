package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hackathon-2025/internal/database"
	"hackathon-2025/pkg/models"
	"hackathon-2025/pkg/services"

	"gorm.io/gorm"
)

// generateAndSaveAISummary is a helper function to generate AI summary asynchronously
func generateAndSaveAISummary(userEmail string) {
	db := database.GetDB()
	if db == nil {
		log.Printf("Database connection not available for AI summary generation for %s", userEmail)
		return
	}

	// Get user pages
	var userPages []models.UserPage
	err := db.Where("user_email = ? AND deleted_at IS NULL", userEmail).Order("page_id DESC").Find(&userPages).Error
	if err != nil {
		log.Printf("Error getting pages for user %s during AI summary generation: %v", userEmail, err)
		return
	}

	// If no pages found, try to sync from Confluence
	if len(userPages) == 0 {
		userPages, err = syncUserPagesFromConfluence(db, userEmail)
		if err != nil {
			log.Printf("Error syncing pages from Confluence for %s: %v", userEmail, err)
			return
		}
	}

	// Limit to 5 pages for AI processing
	var limitPages []models.UserPage
	if len(userPages) > 5 {
		limitPages = userPages[:5]
	} else {
		limitPages = userPages
	}

	// Create comma-separated string of page IDs
	var pageIds []string
	for _, page := range limitPages {
		if page.PageID != "" {
			pageIds = append(pageIds, page.PageID)
		}
	}
	pageIdsStr := strings.Join(pageIds, ",")

	// Call AI service
	agentAns, agentErr := services.SmartAgentInvoke(pageIdsStr, services.SmartAgentRequest{
		EndpointDeploymentHashID: "2i3rtmg8ssjttba1724ynqde",
		EndpointDeploymentKey:    "yx4q46lnxo6kpyals96vm43c",
		UserID:                   "2i3rtmg8ssjttba1724ynqde",
	})

	if agentErr != nil {
		log.Printf("Error calling SmartAgentInvoke for %s: %v", userEmail, agentErr)
		return
	}

	// Clean the response string
	cleanResponseStr := strings.TrimSpace(agentAns.Data.Response.ResponseStr)
	cleanResponseStr = strings.TrimPrefix(cleanResponseStr, "```json")
	cleanResponseStr = strings.TrimPrefix(cleanResponseStr, "```")
	cleanResponseStr = strings.TrimSuffix(cleanResponseStr, "```")
	cleanResponseStr = strings.TrimSpace(cleanResponseStr)
	cleanResponseStr = strings.Trim(cleanResponseStr, "\"'")

	// Parse the AI response
	var responseData map[string]interface{}
	if err := json.Unmarshal([]byte(cleanResponseStr), &responseData); err != nil {
		log.Printf("Error unmarshaling AI response for %s: %v", userEmail, err)
		return
	}

	processedData := processProfileResponse(responseData)

	// Convert to JSON string for storage
	aiSummaryJSON, err := json.Marshal(processedData)
	if err != nil {
		log.Printf("Error marshaling AI summary for %s: %v", userEmail, err)
		return
	}

	// Update the user profile with AI summary
	if err := db.Model(&models.UserProfile{}).Where("user_email = ? AND deleted_at IS NULL", userEmail).Update("ai_summary", string(aiSummaryJSON)).Error; err != nil {
		log.Printf("Error updating AI summary for %s: %v", userEmail, err)
		return
	}

	log.Printf("Successfully generated and saved AI summary for %s", userEmail)
}

// CreateUserProfileHandler handles POST requests to create a new user profile
func CreateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userProfile models.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&userProfile); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if userProfile.UserEmail == "" {
		http.Error(w, "user_email is required", http.StatusBadRequest)
		return
	}

	if userProfile.UserName == "" {
		http.Error(w, "user_name is required", http.StatusBadRequest)
		return
	}

	if userProfile.ProfileImg == "" {
		http.Error(w, "profile_img is required", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	// Check if user profile already exists
	var existingProfile models.UserProfile
	if err := db.Where("user_email = ? AND deleted_at IS NULL", userProfile.UserEmail).First(&existingProfile).Error; err == nil {
		// User profile exists, update it but preserve AISummary
		updateData := map[string]interface{}{
			"user_name":   userProfile.UserName,
			"profile_img": userProfile.ProfileImg,
		}

		if err := db.Model(&existingProfile).Updates(updateData).Error; err != nil {
			log.Printf("Error updating existing user profile: %v", err)
			http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
			return
		}

		// Get the updated profile
		if err := db.Where("user_email = ? AND deleted_at IS NULL", userProfile.UserEmail).First(&userProfile).Error; err != nil {
			log.Printf("Error retrieving updated user profile: %v", err)
			http.Error(w, "Failed to retrieve updated user profile", http.StatusInternalServerError)
			return
		}

		response := UserProfileResponse{
			Data:      userProfile,
			Message:   "User profile updated successfully",
			Timestamp: time.Now(),
			Status:    "success",
		}

		// Check if AI summary is empty and trigger generation asynchronously
		if userProfile.AISummary == "" {
			go generateAndSaveAISummary(userProfile.UserEmail)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Create new user profile
	if err := db.Create(&userProfile).Error; err != nil {
		log.Printf("Error creating user profile: %v", err)
		http.Error(w, "Failed to create user profile", http.StatusInternalServerError)
		return
	}

	response := UserProfileResponse{
		Data:      userProfile,
		Message:   "User profile created successfully",
		Timestamp: time.Now(),
		Status:    "success",
	}

	// Trigger AI summary generation asynchronously for new profiles
	go generateAndSaveAISummary(userProfile.UserEmail)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetUserProfileHandler handles GET requests to retrieve a user profile
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email parameter is required", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	var userProfile models.UserProfile
	if err := db.Where("user_email = ? AND deleted_at IS NULL", email).First(&userProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User profile not found", http.StatusNotFound)
		} else {
			log.Printf("Error retrieving user profile: %v", err)
			http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		}
		return
	}

	response := UserProfileResponse{
		Data:      userProfile,
		Message:   "User profile retrieved successfully",
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

// UpdateUserProfileHandler handles PUT requests to update a user profile
func UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userProfile models.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&userProfile); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if userProfile.UserEmail == "" {
		http.Error(w, "user_email is required", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	// Check if user profile exists
	var existingProfile models.UserProfile
	if err := db.Where("user_email = ? AND deleted_at IS NULL", userProfile.UserEmail).First(&existingProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User profile not found", http.StatusNotFound)
		} else {
			log.Printf("Error checking existing user profile: %v", err)
			http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		}
		return
	}

	// Update the profile
	if err := db.Model(&existingProfile).Updates(userProfile).Error; err != nil {
		log.Printf("Error updating user profile: %v", err)
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	// Get the updated profile
	if err := db.Where("user_email = ? AND deleted_at IS NULL", userProfile.UserEmail).First(&userProfile).Error; err != nil {
		log.Printf("Error retrieving updated user profile: %v", err)
		http.Error(w, "Failed to retrieve updated user profile", http.StatusInternalServerError)
		return
	}

	response := UserProfileResponse{
		Data:      userProfile,
		Message:   "User profile updated successfully",
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

// DeleteUserProfileHandler handles DELETE requests to delete a user profile
func DeleteUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email parameter is required", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	// Check if user profile exists
	var userProfile models.UserProfile
	if err := db.Where("user_email = ? AND deleted_at IS NULL", email).First(&userProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User profile not found", http.StatusNotFound)
		} else {
			log.Printf("Error checking existing user profile: %v", err)
			http.Error(w, "Failed to delete user profile", http.StatusInternalServerError)
		}
		return
	}

	// Soft delete the profile
	if err := db.Delete(&userProfile).Error; err != nil {
		log.Printf("Error deleting user profile: %v", err)
		http.Error(w, "Failed to delete user profile", http.StatusInternalServerError)
		return
	}

	response := Response{
		Message:   "User profile deleted successfully",
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

// ListUserProfilesHandler handles GET requests to list user profiles
func ListUserProfilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Convert limit and offset to integers with proper error handling
	var limit, offset int
	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter. Must be a valid integer", http.StatusBadRequest)
			return
		}
		if limit <= 0 {
			http.Error(w, "Limit must be greater than 0", http.StatusBadRequest)
			return
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Invalid offset parameter. Must be a valid integer", http.StatusBadRequest)
			return
		}
		if offset < 0 {
			http.Error(w, "Offset must be greater than or equal to 0", http.StatusBadRequest)
			return
		}
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	query := db.Where("deleted_at IS NULL")

	// Apply pagination if provided
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	var userProfiles []models.UserProfile
	if err := query.Find(&userProfiles).Error; err != nil {
		log.Printf("Error retrieving user profiles: %v", err)
		http.Error(w, "Failed to retrieve user profiles", http.StatusInternalServerError)
		return
	}

	response := UserProfileListResponse{
		Data:      userProfiles,
		Count:     len(userProfiles),
		Message:   "User profiles retrieved successfully",
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
