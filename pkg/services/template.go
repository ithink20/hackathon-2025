package services

import (
	"fmt"
	"hackathon-2025/internal/database"
	"hackathon-2025/pkg/models"

	"gorm.io/gorm"
)

type TemplateService struct {
	db *gorm.DB
}

func NewTemplateService() *TemplateService {
	return &TemplateService{
		db: database.GetDB(),
	}
}

func (ts *TemplateService) GetTemplateContentByType(templateType string) (string, error) {
	var template models.Template

	err := ts.db.Where("template_type = ? AND deleted_at IS NULL", templateType).First(&template).Error
	if err != nil {
		return "", fmt.Errorf("failed to get template content for type %s: %w", templateType, err)
	}

	return template.TemplateContent, nil
}

func (ts *TemplateService) GetTemplateByType(templateType string) (*models.Template, error) {
	var template models.Template

	err := ts.db.Where("template_type = ? AND deleted_at IS NULL", templateType).First(&template).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get template for type %s: %w", templateType, err)
	}

	return &template, nil
}

func (ts *TemplateService) GetAllTemplates() ([]models.Template, error) {
	var templates []models.Template

	err := ts.db.Where("deleted_at IS NULL").Find(&templates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all templates: %w", err)
	}

	return templates, nil
}

func (ts *TemplateService) GetTemplatesByTypes(templateTypes []string) ([]models.Template, error) {
	var templates []models.Template

	err := ts.db.Where("template_type IN ? AND deleted_at IS NULL", templateTypes).Find(&templates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get templates for types %v: %w", templateTypes, err)
	}

	return templates, nil
}

func (ts *TemplateService) TemplateExists(templateType string) (bool, error) {
	var count int64

	err := ts.db.Model(&models.Template{}).Where("template_type = ? AND deleted_at IS NULL", templateType).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if template exists for type %s: %w", templateType, err)
	}

	return count > 0, nil
}
