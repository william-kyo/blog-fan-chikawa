package service

import (
	"blog-fanchiikawa-service/sdk"
	"context"
	"fmt"
	"io"
)

// CommentReplyService defines the interface for comment reply generation operations
type CommentReplyService interface {
	GenerateCommentRepliesFromUpload(ctx context.Context, imageFile io.ReadSeeker, originalComment string) (*sdk.CommentReplyResponse, error)
}

// commentReplyService implements CommentReplyService interface
type commentReplyService struct {
	anthropicService *sdk.AnthropicService
}

// NewCommentReplyService creates a new CommentReplyService instance
func NewCommentReplyService() CommentReplyService {
	return &commentReplyService{
		anthropicService: sdk.NewAnthropicService(),
	}
}

// GenerateCommentRepliesFromUpload generates comment replies from uploaded image and comment text
func (s *commentReplyService) GenerateCommentRepliesFromUpload(
	ctx context.Context,
	imageFile io.ReadSeeker,
	originalComment string,
) (*sdk.CommentReplyResponse, error) {
	// Read image data
	imageData, err := io.ReadAll(imageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	// Validate image size (max 10MB)
	maxSize := 10 * 1024 * 1024 // 10MB
	if len(imageData) > maxSize {
		return nil, fmt.Errorf("image file too large (max 10MB)")
	}

	// Validate image format by checking the data
	if len(imageData) < 8 {
		return nil, fmt.Errorf("invalid image file")
	}

	// Basic validation for common image formats
	if !s.isValidImageFormat(imageData) {
		return nil, fmt.Errorf("unsupported image format (supported: JPEG, PNG, GIF, WebP)")
	}

	// Validate comment
	if originalComment == "" {
		return nil, fmt.Errorf("original comment cannot be empty")
	}

	if len(originalComment) > 1000 {
		return nil, fmt.Errorf("comment too long (max 1000 characters)")
	}

	// Create request
	req := sdk.CommentReplyRequest{
		ImageData:       imageData,
		OriginalComment: originalComment,
	}

	// Generate replies
	response, err := s.anthropicService.GenerateCommentReplies(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate comment replies: %w", err)
	}

	return response, nil
}

// isValidImageFormat checks if the image data represents a valid image format
func (s *commentReplyService) isValidImageFormat(data []byte) bool {
	if len(data) < 8 {
		return false
	}

	// JPEG
	if data[0] == 0xFF && data[1] == 0xD8 {
		return true
	}

	// PNG
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return true
	}

	// GIF
	if len(data) >= 6 {
		gif87a := string(data[0:6]) == "GIF87a"
		gif89a := string(data[0:6]) == "GIF89a"
		if gif87a || gif89a {
			return true
		}
	}

	// WebP
	if len(data) >= 12 {
		riff := string(data[0:4]) == "RIFF"
		webp := string(data[8:12]) == "WEBP"
		if riff && webp {
			return true
		}
	}

	return false
}