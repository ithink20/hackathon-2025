package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"hackathon-2025/internal/database"
	"hackathon-2025/pkg/models"
)

// UserPostHandler handles CRUD operations for user posts
func UserPostHandler(w http.ResponseWriter, r *http.Request) {
	opType := r.URL.Query().Get("op_type")

	switch opType {
	case "create":
		createPost(w, r)
	case "read":
		readPost(w, r)
	case "update":
		updatePost(w, r)
	case "delete":
		deletePost(w, r)
	case "list":
		listPosts(w, r)
	default:
		http.Error(w, "Invalid op_type. Must be one of: create, read, update, delete, list", http.StatusBadRequest)
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	// Generate unique post ID (you might want to use UUID in production)
	rand.Seed(time.Now().UnixNano())
	postID := fmt.Sprintf("%d", rand.Intn(900000)+100000)

	userPost := models.UserPost{
		PostID:      postID,
		Title:       req.Title,
		PostType:    req.Type,
		Content:     req.Content,
		AuthorName:  req.AuthorName,
		AuthorImage: req.AuthorImage,
		AuthorId:    req.AuthorID,
		Timestamp:   time.Now().Unix(),
		Metadata: models.PostMetadata{
			Tags:     req.Tags,
			Comments: req.Comments,
		},
		Likes: req.Likes,
	}

	if err := db.Create(&userPost).Error; err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	response := models.PostResponse{
		Post:      &userPost,
		Message:   "Post created successfully",
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

func readPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "post_id parameter is required", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	var userPost models.UserPost
	if err := db.Where("post_id = ? AND deleted_at IS NULL", postID).First(&userPost).Error; err != nil {
		log.Printf("Error getting post %s: %v", postID, err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	response := models.PostResponse{
		Post:      &userPost,
		Message:   "Post retrieved successfully",
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

func updatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "post_id parameter is required", http.StatusBadRequest)
		return
	}

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	var userPost models.UserPost
	if err := db.Where("post_id = ? AND deleted_at IS NULL", postID).First(&userPost).Error; err != nil {
		log.Printf("Error getting post %s: %v", postID, err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Update all fields with the provided values
	userPost.Title = req.Title
	userPost.Content = req.Content
	userPost.AuthorName = req.AuthorName
	userPost.AuthorImage = req.AuthorImage
	userPost.Metadata.Tags = req.Tags
	userPost.Metadata.Comments = req.Comments
	userPost.Likes = req.Likes
	userPost.PostType = req.Type
	userPost.AuthorId = req.AuthorID

	if err := db.Save(&userPost).Error; err != nil {
		log.Printf("Error updating post %s: %v", postID, err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	response := models.PostResponse{
		Post:      &userPost,
		Message:   "Post updated successfully",
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

func deletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "post_id parameter is required", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	var userPost models.UserPost
	if err := db.Where("post_id = ? AND deleted_at IS NULL", postID).First(&userPost).Error; err != nil {
		log.Printf("Error getting post %s: %v", postID, err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if err := db.Delete(&userPost).Error; err != nil {
		log.Printf("Error deleting post %s: %v", postID, err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	response := models.PostResponse{
		Message:   "Post deleted successfully",
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

func listPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	// Get query parameters for pagination and filtering
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	searchKeyword := r.URL.Query().Get("search")
	postType := r.URL.Query().Get("post_type")
	authorId := r.URL.Query().Get("author_id")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var userPosts []models.UserPost
	query := db.Where("deleted_at IS NULL")

	if searchKeyword != "" {
		searchPattern := "%" + searchKeyword + "%"
		query = query.Where(
			"title ILIKE ? OR content ILIKE ? OR metadata::text ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Add post_type filtering if provided
	if postType != "" {
		postTypePattern := "%" + postType + "%"
		query = query.Where("post_type ILIKE ?", postTypePattern)
	}

	if authorId != "" {
		query = query.Where("author_id = ?", authorId)
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&userPosts).Error; err != nil {
		log.Printf("Error getting posts: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	response := models.PostResponse{
		Posts:     userPosts,
		Message:   fmt.Sprintf("Retrieved %d posts", len(userPosts)),
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
