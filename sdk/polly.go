package sdk

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var Polly *polly.Polly

func InitPolly() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
		Profile: "fanchiikawa",
	})
	if err != nil {
		panic(err)
	}
	Polly = polly.New(sess)
}

func TextToSpeech(text, languageCode string) (string, error) {
	// Map language code to Polly supported language code
	pollyLanguageCode := mapToPollyLanguageCode(languageCode)
	
	// Generate audio using Polly
	input := &polly.SynthesizeSpeechInput{
		Text:         aws.String(text),
		OutputFormat: aws.String("mp3"),
		VoiceId:      getVoiceIdByLanguage(pollyLanguageCode),
		LanguageCode: aws.String(pollyLanguageCode),
	}

	result, err := Polly.SynthesizeSpeech(input)
	if err != nil {
		return "", err
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("speech/%d_%s.mp3", timestamp, languageCode)

	// Upload to S3 using the existing session
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
		Profile: "fanchiikawa",
	})
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("fan-ai-warehouse"),
		Key:    aws.String(filename),
		Body:   result.AudioStream,
	})
	if err != nil {
		return "", err
	}

	return filename, nil
}

func mapToPollyLanguageCode(comprehendLanguageCode string) string {
	languageMap := map[string]string{
		"en":    "en-US",
		"zh":    "cmn-CN",
		"zh-CN": "cmn-CN",
		"ja":    "ja-JP",
		"ko":    "ko-KR",
		"fr":    "fr-FR",
		"de":    "de-DE",
		"es":    "es-ES",
		"it":    "it-IT",
		"pt":    "pt-PT",
		"ru":    "ru-RU",
		"ar":    "arb",
		"hi":    "hi-IN",
	}

	if pollyCode, exists := languageMap[comprehendLanguageCode]; exists {
		return pollyCode
	}
	// Default to English if language not supported
	return "en-US"
}

func getVoiceIdByLanguage(languageCode string) *string {
	voiceMap := map[string]string{
		"en-US":  "Joanna",
		"cmn-CN": "Zhiyu",
		"ja-JP":  "Mizuki",
		"ko-KR":  "Seoyeon",
		"fr-FR":  "Celine",
		"de-DE":  "Marlene",
		"es-ES":  "Conchita",
		"it-IT":  "Carla",
		"pt-PT":  "Ines",
		"ru-RU":  "Tatyana",
		"arb":    "Zeina",
		"hi-IN":  "Aditi",
	}

	if voice, exists := voiceMap[languageCode]; exists {
		return aws.String(voice)
	}
	// Default to English voice if language not supported
	return aws.String("Joanna")
}