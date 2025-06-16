package sdk

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var AWSSession *session.Session

func InitAWSSession() {
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "ap-northeast-1" // fallback default
	}

	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "fanchiikawa" // fallback default
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		Profile: profile,
	})
	if err != nil {
		panic(err)
	}

	AWSSession = sess
}