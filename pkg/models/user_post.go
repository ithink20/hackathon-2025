package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Comment represents a comment on a post
type Comment struct {
	ID         string `json:"id"`
	AuthorName string `json:"authorName"`
	AuthorImg  string `json:"authorImg"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
	Likes      int    `json:"likes"`
}

// PostMetadata contains tags and comments
type PostMetadata struct {
	Tags     []string  `json:"tags,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
}

// Value implements driver.Valuer interface for JSONB
func (pm PostMetadata) Value() (driver.Value, error) {
	return json.Marshal(pm)
}

// Scan implements sql.Scanner interface for JSONB
func (pm *PostMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, pm)
}

// UserPost represents a user post
type UserPost struct {
	PostID      string         `json:"post_id" gorm:"primaryKey;column:post_id"`
	PostType    string         `json:"post_type" gorm:"column:post_type"`
	Title       string         `json:"title" gorm:"not null"`
	Content     string         `json:"content" gorm:"not null"`
	AuthorName  string         `json:"authorName" gorm:"column:author_name;not null"`
	AuthorImage string         `json:"authorImg" gorm:"column:author_image"`
	AuthorId    string         `json:"author_id" gorm:"column:author_id"`
	Timestamp   int64          `json:"timestamp" gorm:"not null"`
	Metadata    PostMetadata   `json:"metaData" gorm:"type:jsonb"`
	Likes       int            `json:"likes" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (UserPost) TableName() string {
	return "user_post"
}

// CreatePostRequest represents the request payload for creating/updating a post
type CreatePostRequest struct {
	Title       string    `json:"title" validate:"required"`
	Type        string    `json:"type"`
	Content     string    `json:"content" validate:"required"`
	AuthorName  string    `json:"authorName" validate:"required"`
	AuthorImage string    `json:"authorImg"`
	AuthorID    string    `json:"authorId" validate:"required"`
	Tags        []string  `json:"tags,omitempty"`
	Comments    []Comment `json:"comments,omitempty"`
	Likes       int       `json:"likes"`
}

// PostResponse represents the response for post operations
type PostResponse struct {
	Post      *UserPost  `json:"post,omitempty"`
	Posts     []UserPost `json:"posts,omitempty"`
	Message   string     `json:"message"`
	Timestamp time.Time  `json:"timestamp"`
	Status    string     `json:"status"`
	Error     string     `json:"error,omitempty"`
}
