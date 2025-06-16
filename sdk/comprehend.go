package sdk

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

var Comprehend *comprehend.Comprehend

func InitComprehend() {
	Comprehend = comprehend.New(AWSSession)
}

func DetectLanguage(text string) (string, error) {
	result, err := Comprehend.DetectDominantLanguage(&comprehend.DetectDominantLanguageInput{
		Text: aws.String(text),
	})
	if err != nil {
		return "en", err
	}
	return *result.Languages[0].LanguageCode, nil
}

func DetectSentiment(text string) (string, error) {
	result, err := Comprehend.DetectSentiment(&comprehend.DetectSentimentInput{
		Text:         aws.String(text),
		LanguageCode: aws.String("en"),
	})
	if err != nil {
		return "UNKNOWN", err
	}
	return *result.Sentiment, nil
}
