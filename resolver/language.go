package resolver

import (
	"blog-fanchiikawa-service/graph/model"
	"context"
)

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