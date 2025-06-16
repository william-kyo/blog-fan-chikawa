package service

import (
	"blog-fanchiikawa-service/sdk"
	"fmt"
)

// TranslateService defines the interface for translation operations
type TranslateService interface {
	// TranslateText translates text from source language to target language
	TranslateText(text, sourceLanguage, targetLanguage string) (string, error)
}

// translateService implements TranslateService interface
type translateService struct{}

// NewTranslateService creates a new TranslateService instance
func NewTranslateService() TranslateService {
	return &translateService{}
}

// TranslateText translates text from source language to target language
func (s *translateService) TranslateText(text, sourceLanguage, targetLanguage string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("text cannot be empty")
	}
	
	if sourceLanguage == "" || targetLanguage == "" {
		return "", fmt.Errorf("source and target languages must be specified")
	}

	translatedText, err := sdk.TranslateText(text, sourceLanguage, targetLanguage)
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %w", err)
	}
	
	return translatedText, nil
}