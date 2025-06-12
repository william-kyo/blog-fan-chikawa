package sdk

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

var Comprehend *comprehend.Comprehend

func InitComprehend() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
		Profile: "fanchiikawa",
	})
	if err != nil {
		panic(err)
	}
	Comprehend = comprehend.New(sess)
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
