package service

import (
	"blog-fanchiikawa-service/sdk"
	"fmt"
)

// LanguageService defines the interface for language-related operations
type LanguageService interface {
	// DetectLanguage detects the language of the given text
	DetectLanguage(text string) (string, error)
	
	// DetectSentiment analyzes the sentiment of the given text
	DetectSentiment(text string) (string, error)
}

// languageService implements LanguageService interface
type languageService struct{}

// NewLanguageService creates a new LanguageService instance
func NewLanguageService() LanguageService {
	return &languageService{}
}

// DetectLanguage detects the language of the given text
func (s *languageService) DetectLanguage(text string) (string, error) {
	language, err := sdk.DetectLanguage(text)
	if err != nil {
		return "", fmt.Errorf("failed to detect language: %w", err)
	}
	return language, nil
}

// DetectSentiment analyzes the sentiment of the given text
func (s *languageService) DetectSentiment(text string) (string, error) {
	sentiment, err := sdk.DetectSentiment(text)
	if err != nil {
		return "", fmt.Errorf("failed to detect sentiment: %w", err)
	}
	return sentiment, nil
}