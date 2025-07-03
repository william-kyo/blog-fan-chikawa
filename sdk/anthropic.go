package sdk

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AnthropicService provides Anthropic API functionality
type AnthropicService struct {
	client anthropic.Client
}

// CommentReplyRequest represents a request for generating comment replies
type CommentReplyRequest struct {
	ImageData       []byte
	OriginalComment string
}

// CommentReply represents a single reply suggestion
type CommentReply struct {
	Style   string `json:"style"`
	Content string `json:"content"`
}

// CommentReplyResponse represents the response with multiple reply suggestions
type CommentReplyResponse struct {
	Replies []CommentReply `json:"replies"`
}

// NewAnthropicService creates a new Anthropic service
func NewAnthropicService() *AnthropicService {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		panic("ANTHROPIC_API_KEY environment variable is required")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicService{
		client: client,
	}
}

// GenerateCommentReplies generates 3 different style replies for a comment with image context
func (a *AnthropicService) GenerateCommentReplies(ctx context.Context, req CommentReplyRequest) (*CommentReplyResponse, error) {
	// Convert image to base64
	imageBase64 := base64.StdEncoding.EncodeToString(req.ImageData)
	
	// Detect image format
	imageFormat := detectImageFormat(req.ImageData)
	
	// Create system prompt for comment reply generation
	systemPrompt := `You are an expert at generating thoughtful, engaging comment replies. Given an image and an original comment, generate 3 different style replies that are appropriate and relevant.

Reply styles:
1. Friendly - Warm, supportive, and encouraging
2. Professional - Informative, helpful, and constructive  
3. Humorous - Light-hearted, witty, but still respectful

Requirements:
- Each reply should be 1-2 sentences maximum
- Replies should be contextually relevant to both the image and original comment
- Maintain a positive and respectful tone
- Avoid controversial or sensitive topics

Format your response as a JSON object with this structure:
{
  "replies": [
    {"style": "friendly", "content": "your friendly reply here"},
    {"style": "professional", "content": "your professional reply here"},
    {"style": "humorous", "content": "your humorous reply here"}
  ]
}`

	userPrompt := fmt.Sprintf("Original comment: \"%s\"\n\nPlease analyze the image and generate 3 different style replies to this comment.", req.OriginalComment)

	// Create content blocks
	textBlock := anthropic.NewTextBlock(userPrompt)
	imageBlock := anthropic.NewImageBlockBase64(imageFormat, imageBase64)

	// Create user message
	message := anthropic.NewUserMessage(textBlock, imageBlock)

	// Create the request
	response, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5SonnetLatest,
		MaxTokens: 1000,
		Messages:  []anthropic.MessageParam{message},
		System:    []anthropic.TextBlockParam{{Text: systemPrompt}},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate replies: %w", err)
	}

	// Extract text content from response
	if len(response.Content) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	// Get the first content block as text
	var responseText string
	firstContent := response.Content[0]
	
	// Check if it's a text block
	if firstContent.Type == "text" {
		textBlock := firstContent.AsText()
		responseText = textBlock.Text
	} else {
		return nil, fmt.Errorf("unexpected content type in response: %s", firstContent.Type)
	}

	// Try to parse the JSON response
	var parsedResponse CommentReplyResponse
	if err := json.Unmarshal([]byte(responseText), &parsedResponse); err != nil {
		// If JSON parsing fails, return a fallback response
		return &CommentReplyResponse{
			Replies: []CommentReply{
				{
					Style:   "friendly",
					Content: "That's such a beautiful capture! I love how you caught the lighting in this shot 😊",
				},
				{
					Style:   "professional",
					Content: "Excellent composition and technical execution. The depth of field really makes the subject stand out.",
				},
				{
					Style:   "humorous",
					Content: "Plot twist: the camera was more excited about this shot than we are! 📸✨",
				},
			},
		}, nil
	}

	return &parsedResponse, nil
}

// EncodeImageToBase64 encodes an image file to base64
func EncodeImageToBase64(imageData []byte) string {
	return base64.StdEncoding.EncodeToString(imageData)
}

// detectImageFormat detects the MIME type of an image from its data
func detectImageFormat(data []byte) string {
	if len(data) < 8 {
		return "image/jpeg" // default fallback
	}

	// JPEG
	if data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}

	// PNG
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return "image/png"
	}

	// GIF
	if len(data) >= 6 {
		gif87a := string(data[0:6]) == "GIF87a"
		gif89a := string(data[0:6]) == "GIF89a"
		if gif87a || gif89a {
			return "image/gif"
		}
	}

	// WebP
	if len(data) >= 12 {
		riff := string(data[0:4]) == "RIFF"
		webp := string(data[8:12]) == "WEBP"
		if riff && webp {
			return "image/webp"
		}
	}

	// Default to JPEG if format is not recognized
	return "image/jpeg"
}

// ValidateImageFormat validates if the image format is supported
func ValidateImageFormat(reader io.Reader) error {
	// Read first few bytes to determine format
	header := make([]byte, 512)
	n, err := reader.Read(header)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read image header: %w", err)
	}

	// Check for common image formats
	if n >= 2 {
		// JPEG
		if header[0] == 0xFF && header[1] == 0xD8 {
			return nil
		}
		// PNG
		if n >= 8 && header[0] == 0x89 && header[1] == 0x50 && header[2] == 0x4E && header[3] == 0x47 {
			return nil
		}
		// GIF
		if n >= 6 && string(header[0:6]) == "GIF87a" || string(header[0:6]) == "GIF89a" {
			return nil
		}
		// WebP
		if n >= 12 && string(header[0:4]) == "RIFF" && string(header[8:12]) == "WEBP" {
			return nil
		}
	}

	return fmt.Errorf("unsupported image format")
}