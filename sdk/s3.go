package sdk

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3 *s3.S3
var prefix = "warehouse/"
var bucket = "fan-ai-warehouse"

func InitS3() {
	S3 = s3.New(AWSSession)
}

func GetLastData() (string, error) {

	// Get the last date from S3 bucket
	result, err := S3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"), // Use delimiter to get only folders
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
	folderName := strings.TrimPrefix(*lastFolder, prefix)
	folderName = strings.TrimSuffix(folderName, "/")

	return folderName, nil
}

type UploadResult struct {
	Filename string
	S3Bucket string
	S3Key    string
	Error    error
}

func UploadMultipleFiles(uploadfiles []string) []UploadResult {
	jobs := make(chan string, len(uploadfiles))
	ret := make(chan UploadResult, len(uploadfiles))

	var wg sync.WaitGroup
	for range len(uploadfiles) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for filename := range jobs {
				// Generate S3 key name
				s3Key := prefix + "image/" + filepath.Base(filename)

				// Upload file
				err := UploadFile(s3Key, filename)
				ret <- UploadResult{
					Filename: filename,
					S3Bucket: bucket,
					S3Key:    s3Key,
					Error:    err,
				}
			}

		}()
	}

	for _, file := range uploadfiles {
		jobs <- file
	}

	// Wait for all tasks done
	go func() {
		wg.Wait()
		close(ret)
	}()

	close(jobs)

	var uploadresults []UploadResult
	for r := range ret {
		uploadresults = append(uploadresults, r)
	}

	return uploadresults
}

func UploadFile(key string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Printf("Upload result: %s for key: %s", err, key)
	}
	return err
}
