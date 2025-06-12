package sdk

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3 *s3.S3

func InitS3() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
		Profile: "fanchiikawa",
	})
	if err != nil {
		panic(err)
	}

	S3 = s3.New(sess)
}

func GetLastData() (string, error) {
	// Get the last date from S3 bucket
	result, err := S3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String("fan-ai-warehouse"),
		Prefix:    aws.String(""),
		Delimiter: aws.String("/"), // 使用分隔符来只获取文件夹
	})
	if err != nil {
		return "", err
	}

	// If no common prefixes (folders) found, return empty string
	if len(result.CommonPrefixes) == 0 {
		return "", nil
	}

	// Sort the prefixes (folders) by name
	sort.Slice(result.CommonPrefixes, func(i, j int) bool {
		return *result.CommonPrefixes[i].Prefix < *result.CommonPrefixes[j].Prefix
	})

	// Get the last folder (most recent)
	lastFolder := result.CommonPrefixes[len(result.CommonPrefixes)-1].Prefix

	// Remove trailing slash
	folderName := strings.TrimSuffix(*lastFolder, "/")

	return folderName, nil
}
