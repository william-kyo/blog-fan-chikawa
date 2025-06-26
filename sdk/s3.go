package sdk

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var S3 *s3.S3
var S3Uploader *s3manager.Uploader
var prefix = "warehouse/"
var bucket = "fan-ai-warehouse"

func InitS3() {
	S3 = s3.New(AWSSession)
	S3Uploader = s3manager.NewUploader(AWSSession)
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
				s3Key := prefix + filepath.Base(filename)

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

// CustomLabelsUploadResult represents the result of a custom labels upload
type CustomLabelsUploadResult struct {
	Key      string `json:"key"`
	Location string `json:"location"`
	Bucket   string `json:"bucket"`
}

// UploadFileForCustomLabels uploads a multipart file to S3 bucket for custom labels detection
func UploadFileForCustomLabels(file multipart.File, filename string, bucketName string) (*CustomLabelsUploadResult, error) {
	// Generate unique key using timestamp and original filename
	ext := filepath.Ext(filename)
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("warehouse/image/%d_%s", timestamp, filename)

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Determine content type
	contentType := getContentType(ext)

	// Upload to S3
	result, err := S3Uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String(contentType),
		ACL:         aws.String("private"), // Private access
	})

	if err != nil {
		log.Printf("Failed to upload file to S3: %v", err)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("Successfully uploaded file to S3: %s", result.Location)

	return &CustomLabelsUploadResult{
		Key:      key,
		Location: result.Location,
		Bucket:   bucketName,
	}, nil
}

// getContentType determines content type based on file extension
func getContentType(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// GeneratePresignedURL generates a presigned URL for viewing the uploaded image
func GeneratePresignedURL(bucketName, key string) (string, error) {
	req, _ := S3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	// URL expires in 1 hour
	url, err := req.Presign(1 * time.Hour)
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}

// S3PresignedUploadResult represents presigned upload URL data
type S3PresignedUploadResult struct {
	UploadURL string            `json:"uploadUrl"`
	Key       string            `json:"key"`
	Fields    map[string]string `json:"fields"`
}

// GeneratePresignedUploadURL generates a presigned URL for direct upload to S3
func GeneratePresignedUploadURL(filename, bucketName string) (*S3PresignedUploadResult, error) {
	// Generate unique key
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("warehouse/image/%d_%s", timestamp, filename)

	// Create a simple presigned PUT URL for easier frontend implementation
	req, _ := S3.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		ACL:    aws.String("private"),
	})

	// Generate presigned URL (15 minutes expiry)
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Printf("Failed to generate presigned upload URL: %v", err)
		return nil, fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	// For PUT requests, we don't need additional fields
	fieldsMap := make(map[string]string)

	return &S3PresignedUploadResult{
		UploadURL: url,
		Key:       key,
		Fields:    fieldsMap,
	}, nil
}
