package resolver

import (
	"blog-fanchiikawa-service/graph/model"
	"blog-fanchiikawa-service/service"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

// UploadAndDetectCustomLabels handles file upload and custom labels detection
func (r *Resolver) UploadAndDetectCustomLabels(ctx context.Context, file graphql.Upload) (*model.CustomLabelsResult, error) {
	// Basic validation
	if file.File == nil {
		return nil, fmt.Errorf("no file provided")
	}

	// Convert GraphQL upload to service data structure
	uploadData := &service.UploadFileData{
		ReadSeeker: file.File,
		Filename:   file.Filename,
		Size:       file.Size,
	}

	// Call service layer - now returns GraphQL model directly
	return r.CustomLabelsService.UploadAndDetectForResolver(uploadData)
}

// DetectCustomLabelsFromS3 handles custom labels detection from S3 key
func (r *Resolver) DetectCustomLabelsFromS3(ctx context.Context, input model.DetectCustomLabelsInput) (*model.CustomLabelsResult, error) {
	// Call service layer - now returns GraphQL model directly
	return r.CustomLabelsService.DetectFromS3KeyForResolver(input.S3Key)
}

// GenerateS3UploadURL generates presigned URL for S3 upload
func (r *Resolver) GenerateS3UploadURL(ctx context.Context, filename string) (*model.S3PresignedURL, error) {
	// Call service layer
	response, err := r.CustomLabelsService.GenerateUploadURL(filename)
	if err != nil {
		return nil, err
	}

	// Convert fields map to GraphQL model
	var fields []*model.S3Field
	for name, value := range response.Fields {
		fields = append(fields, &model.S3Field{
			Name:  name,
			Value: value,
		})
	}

	return &model.S3PresignedURL{
		UploadURL: response.UploadURL,
		Key:       response.Key,
		Fields:    fields,
	}, nil
}