package sdk

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
)

var Translate *translate.Translate

func InitTranslate() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
		Profile: "fanchiikawa",
	})
	if err != nil {
		panic(err)
	}
	Translate = translate.New(sess)
}

func TranslateText(text string, sourceLanguage string, targetLanguage string) (string, error) {
	result, err := Translate.Text(&translate.TextInput{
		Text:               aws.String(text),
		SourceLanguageCode: aws.String(sourceLanguage),
		TargetLanguageCode: aws.String(targetLanguage),
	})
	if err != nil {
		return "", err
	}
	return *result.TranslatedText, nil
}
