package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"hackathon-2025/internal/database"
	"hackathon-2025/pkg/models"

	"gorm.io/gorm"
)

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
		// User profile exists, update it
		if err := db.Model(&existingProfile).Updates(userProfile).Error; err != nil {
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
