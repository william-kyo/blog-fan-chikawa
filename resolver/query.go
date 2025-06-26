package resolver

import (
	"blog-fanchiikawa-service/graph/model"
	"context"
)

// Users handles the users query
func (r *Resolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.UserService.GetUsers(10) // Default limit of 10
}

// FetchLastData handles the fetchLastData query
func (r *Resolver) FetchLastData(ctx context.Context) (string, error) {
	return r.StorageService.GetLastData()
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