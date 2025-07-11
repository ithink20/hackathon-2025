package models

import (
	"time"

	"gorm.io/gorm"
)

type ConfluenceSearchResponse struct {
	Results   []ConfluencePage `json:"results"`
	Size      int              `json:"size"`
	Start     int              `json:"start"`
	TotalSize int              `json:"totalSize"`
}

type ConfluencePage struct {
	Content ConfluenceContent `json:"content"`
}

type ConfluenceContent struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Title  string `json:"title"`
}

type PageInfo struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content,omitempty"`
	Link    string `json:"link"`
}

type UserPage struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserEmail   string         `json:"user_email" gorm:"index;not null"`
	PageID      string         `json:"page_id" gorm:"uniqueIndex;not null"`
	PageType    string         `json:"page_type" gorm:"not null"`
	PageTitle   string         `json:"page_title" gorm:"not null"`
	PageContent string         `json:"page_content" gorm:"type:longtext"`
	LastUpdated time.Time      `json:"last_updated" gorm:"autoUpdateTime"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	PageLink    string         `json:"page_link" gorm:"page_link"`
}

func (UserPage) TableName() string {
	return "user_pages"
}
