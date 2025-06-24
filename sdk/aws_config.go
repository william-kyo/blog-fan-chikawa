package sdk

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go-v2/config"
	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
)

var AWSSession *session.Session
var AWSConfigV2 awsv2.Config

// getAWSCredentials returns common AWS configuration parameters
func getAWSCredentials() (region, profile string) {
	region = os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "ap-northeast-1" // fallback default
	}

	profile = os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "fanchiikawa" // fallback default
	}
	
	return region, profile
}

// InitAWS initializes both SDK v1 and v2 configurations
func InitAWS() {
	region, profile := getAWSCredentials()

	// Initialize SDK v1 session for existing services
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

	// Initialize SDK v2 config for Lex Runtime V2
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		panic(err)
	}
	AWSConfigV2 = cfg
}

// Legacy methods for backward compatibility
func InitAWSSession() {
	if AWSSession == nil {
		InitAWS()
	}
}

func InitAWSConfigV2() {
	if AWSConfigV2.Region == "" {
		InitAWS()
	}
}

func GetAWSSession() *session.Session {
	if AWSSession == nil {
		InitAWSSession()
	}
	return AWSSession
}

func GetAWSConfig() awsv2.Config {
	if AWSConfigV2.Region == "" {
		InitAWSConfigV2()
	}
	return AWSConfigV2
}