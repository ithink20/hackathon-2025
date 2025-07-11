package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	TemplateType    string         `json:"template_type" gorm:"index;not null"`
	TemplateContent string         `json:"template_content" gorm:"type:text"`
	LastUpdated     time.Time      `json:"last_updated" gorm:"autoUpdateTime"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (Template) TableName() string {
	return "templates"
}
