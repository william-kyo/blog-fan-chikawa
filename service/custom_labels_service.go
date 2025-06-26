package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"blog-fanchiikawa-service/graph/model"
	"blog-fanchiikawa-service/sdk"
)

// UploadFileData represents uploaded file data
type UploadFileData struct {
	ReadSeeker io.ReadSeeker
	Filename   string
	Size       int64
}

// CustomLabelsService handles custom labels detection operations
type CustomLabelsService interface {
	UploadAndDetectFromGraphQL(uploadData *UploadFileData) (*CustomLabelsResponse, error)
	DetectFromS3Key(s3Key string) (*CustomLabelsResponse, error)
	GenerateUploadURL(filename string) (*S3UploadURLResponse, error)
	// New methods that return GraphQL models directly
	UploadAndDetectForResolver(uploadData *UploadFileData) (*model.CustomLabelsResult, error)
	DetectFromS3KeyForResolver(s3Key string) (*model.CustomLabelsResult, error)
}

// S3UploadURLResponse represents presigned upload URL response
type S3UploadURLResponse struct {
	UploadURL string            `json:"uploadUrl"`
	Key       string            `json:"key"`
	Fields    map[string]string `json:"fields"`
}

type customLabelsService struct{}

// CustomLabelsResponse represents the response from custom labels detection
type CustomLabelsResponse struct {
	ImageURL string                      `json:"imageUrl"`
	S3Key    string                      `json:"s3Key"`
	Labels   []sdk.CustomLabelResult     `json:"labels"`
}

// NewCustomLabelsService creates a new custom labels service
func NewCustomLabelsService() CustomLabelsService {
	return &customLabelsService{}
}

// UploadAndDetectFromGraphQL handles GraphQL upload and performs custom labels detection
func (s *customLabelsService) UploadAndDetectFromGraphQL(uploadData *UploadFileData) (*CustomLabelsResponse, error) {
	// Validate input
	if uploadData.ReadSeeker == nil {
		return nil, fmt.Errorf("no file data provided")
	}
	
	if uploadData.Filename == "" {
		return nil, fmt.Errorf("filename is required")
	}

	// Create internal multipart file wrapper
	multipartFile := &serviceFileWrapper{
		ReadSeeker: uploadData.ReadSeeker,
		filename:   uploadData.Filename,
		size:       uploadData.Size,
	}

	return s.uploadAndDetect(multipartFile, uploadData.Filename)
}

// uploadAndDetect uploads image to S3 and performs custom labels detection
func (s *customLabelsService) uploadAndDetect(file multipart.File, filename string) (*CustomLabelsResponse, error) {
	// Get configuration from environment variables
	bucketName := os.Getenv("REKOGNITION_S3_BUCKET")
	if bucketName == "" {
		return nil, fmt.Errorf("REKOGNITION_S3_BUCKET environment variable not set")
	}

	projectVersionArn := os.Getenv("REKOGNITION_PROJECT_VERSION_ARN")
	if projectVersionArn == "" {
		return nil, fmt.Errorf("REKOGNITION_PROJECT_VERSION_ARN environment variable not set")
	}

	// Step 1: Upload file to S3
	uploadResult, err := sdk.UploadFileForCustomLabels(file, filename, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Step 2: Generate presigned URL for display
	imageURL, err := sdk.GeneratePresignedURL(bucketName, uploadResult.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Step 3: Detect custom labels
	labels, err := sdk.DetectCustomLabels(bucketName, uploadResult.Key, projectVersionArn)
	if err != nil {
		return nil, fmt.Errorf("failed to detect custom labels: %w", err)
	}

	return &CustomLabelsResponse{
		ImageURL: imageURL,
		S3Key:    uploadResult.Key,
		Labels:   labels,
	}, nil
}

// serviceFileWrapper wraps ReadSeeker to implement multipart.File interface
type serviceFileWrapper struct {
	io.ReadSeeker
	filename string
	size     int64
}

func (w *serviceFileWrapper) Close() error {
	if closer, ok := w.ReadSeeker.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (w *serviceFileWrapper) ReadAt(p []byte, off int64) (n int, err error) {
	if readerAt, ok := w.ReadSeeker.(io.ReaderAt); ok {
		return readerAt.ReadAt(p, off)
	}
	// Fallback: seek to offset and read
	if _, err := w.Seek(off, io.SeekStart); err != nil {
		return 0, err
	}
	return w.Read(p)
}

// DetectFromS3Key performs custom labels detection on existing S3 object
func (s *customLabelsService) DetectFromS3Key(s3Key string) (*CustomLabelsResponse, error) {
	// Get configuration from environment variables
	bucketName := os.Getenv("REKOGNITION_S3_BUCKET")
	if bucketName == "" {
		return nil, fmt.Errorf("REKOGNITION_S3_BUCKET environment variable not set")
	}

	projectVersionArn := os.Getenv("REKOGNITION_PROJECT_VERSION_ARN")
	if projectVersionArn == "" {
		return nil, fmt.Errorf("REKOGNITION_PROJECT_VERSION_ARN environment variable not set")
	}

	// Validate S3 key
	if s3Key == "" {
		return nil, fmt.Errorf("S3 key is required")
	}

	// Generate presigned URL for display
	imageURL, err := sdk.GeneratePresignedURL(bucketName, s3Key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Detect custom labels
	labels, err := sdk.DetectCustomLabels(bucketName, s3Key, projectVersionArn)
	if err != nil {
		return nil, fmt.Errorf("failed to detect custom labels: %w", err)
	}

	return &CustomLabelsResponse{
		ImageURL: imageURL,
		S3Key:    s3Key,
		Labels:   labels,
	}, nil
}

// GenerateUploadURL generates presigned URL for direct S3 upload
func (s *customLabelsService) GenerateUploadURL(filename string) (*S3UploadURLResponse, error) {
	// Get configuration from environment variables
	bucketName := os.Getenv("REKOGNITION_S3_BUCKET")
	if bucketName == "" {
		return nil, fmt.Errorf("REKOGNITION_S3_BUCKET environment variable not set")
	}

	// Validate filename
	if filename == "" {
		return nil, fmt.Errorf("filename is required")
	}

	// Generate presigned upload URL
	result, err := sdk.GeneratePresignedUploadURL(filename, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return &S3UploadURLResponse{
		UploadURL: result.UploadURL,
		Key:       result.Key,
		Fields:    result.Fields,
	}, nil
}

// UploadAndDetectForResolver handles GraphQL upload and returns GraphQL model directly
func (s *customLabelsService) UploadAndDetectForResolver(uploadData *UploadFileData) (*model.CustomLabelsResult, error) {
	response, err := s.UploadAndDetectFromGraphQL(uploadData)
	if err != nil {
		return nil, err
	}
	return s.convertToGraphQLResult(response), nil
}

// DetectFromS3KeyForResolver performs detection and returns GraphQL model directly
func (s *customLabelsService) DetectFromS3KeyForResolver(s3Key string) (*model.CustomLabelsResult, error) {
	response, err := s.DetectFromS3Key(s3Key)
	if err != nil {
		return nil, err
	}
	return s.convertToGraphQLResult(response), nil
}

// convertToGraphQLResult converts service response to GraphQL model
func (s *customLabelsService) convertToGraphQLResult(response *CustomLabelsResponse) *model.CustomLabelsResult {
	var labels []*model.CustomLabel
	for _, label := range response.Labels {
		labels = append(labels, &model.CustomLabel{
			Name:       label.Name,
			Confidence: label.Confidence,
		})
	}

	return &model.CustomLabelsResult{
		ImageURL: response.ImageURL,
		S3Key:    response.S3Key,
		Labels:   labels,
	}
}

// Ensure serviceFileWrapper implements multipart.File interface
var _ multipart.File = (*serviceFileWrapper)(nil)