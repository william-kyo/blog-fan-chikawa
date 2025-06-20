package service

import (
	"blog-fanchiikawa-service/sdk"
	"fmt"
)

// SpeechService defines the interface for speech-related operations
type SpeechService interface {
	// TextToSpeech converts text to speech and returns the S3 key
	TextToSpeech(text string) (string, error)
}

// speechService implements SpeechService interface
type speechService struct {
	languageService LanguageService
}

// NewSpeechService creates a new SpeechService instance
func NewSpeechService(languageService LanguageService) SpeechService {
	return &speechService{
		languageService: languageService,
	}
}

// TextToSpeech converts text to speech and returns the S3 key
func (s *speechService) TextToSpeech(text string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("text cannot be empty")
	}

	// First detect the language
	languageCode, err := s.languageService.DetectLanguage(text)
	if err != nil {
		return "", fmt.Errorf("unable to detect language type")
	}

	// Generate speech and upload to S3
	s3Key, err := sdk.TextToSpeech(text, languageCode)
	if err != nil {
		return "", fmt.Errorf("failed to generate speech: %w", err)
	}

	return s3Key, nil
}
