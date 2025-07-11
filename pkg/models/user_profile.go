package models

import (
	"gorm.io/gorm"
	"time"
)

type UserProfile struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserEmail   string         `json:"user_email" gorm:"index;not null"`
	AISummary   string         `json:"ai_summary" gorm:"type:text"`
	LastUpdated time.Time      `json:"last_updated" gorm:"autoUpdateTime"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
