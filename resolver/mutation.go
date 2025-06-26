package resolver

import (
	"blog-fanchiikawa-service/graph/model"
	"blog-fanchiikawa-service/service"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

// Login handles the login mutation
func (r *Resolver) Login(ctx context.Context, input model.LoginUser) (*model.User, error) {
	return r.UserService.Login(input.Nickname, input.Email, input.DeviceID)
}

// DetectLanguage handles the detectLanguage mutation
func (r *Resolver) DetectLanguage(ctx context.Context, input string) (string, error) {
	return r.LanguageService.DetectLanguage(input)
}

// DetectSentiment handles the detectSentiment mutation
func (r *Resolver) DetectSentiment(ctx context.Context, input string) (string, error) {
	return r.LanguageService.DetectSentiment(input)
}

// TranslateText handles the translateText mutation
func (r *Resolver) TranslateText(ctx context.Context, input *model.TranslateText) (string, error) {
	return r.TranslateService.TranslateText(input.Text, input.SourceLanguage, input.TargetLanguage)
}

// TextToSpeech handles the textToSpeech mutation
func (r *Resolver) TextToSpeech(ctx context.Context, input model.TextToSpeech) (string, error) {
	return r.SpeechService.TextToSpeech(input.Text)
}

// UploadAndDetectCustomLabels handles file upload and custom labels detection
func (r *Resolver) UploadAndDetectCustomLabels(ctx context.Context, file graphql.Upload) (*model.CustomLabelsResult, error) {
	// Basic validation
	if file.File == nil {
		return nil, fmt.Errorf("no file provided")
	}

	// Convert GraphQL upload to service data structure
	uploadData := &service.UploadFileData{
		ReadSeeker: file.File,
		Filename:   file.Filename,
		Size:       file.Size,
	}

	// Call service layer - now returns GraphQL model directly
	return r.CustomLabelsService.UploadAndDetectForResolver(uploadData)
}

// DetectCustomLabelsFromS3 handles custom labels detection from S3 key
func (r *Resolver) DetectCustomLabelsFromS3(ctx context.Context, input model.DetectCustomLabelsInput) (*model.CustomLabelsResult, error) {
	// Call service layer - now returns GraphQL model directly
	return r.CustomLabelsService.DetectFromS3KeyForResolver(input.S3Key)
}

