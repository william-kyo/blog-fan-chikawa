package resolver

import (
	"blog-fanchiikawa-service/graph/model"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

// GenerateCommentReplies generates comment replies with image context
func (r *Resolver) GenerateCommentReplies(
	ctx context.Context,
	input model.GenerateCommentRepliesInput,
	file graphql.Upload,
) (*model.CommentReplyResponse, error) {
	// Open the uploaded file
	uploadedFile := file.File
	if uploadedFile == nil {
		return nil, fmt.Errorf("invalid file upload")
	}

	// Generate comment replies using the service
	response, err := r.CommentReplyService.GenerateCommentRepliesFromUpload(
		ctx,
		uploadedFile,
		input.OriginalComment,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate comment replies: %v", err)
	}

	// Convert SDK response to GraphQL model
	var replies []*model.CommentReply
	for _, reply := range response.Replies {
		replies = append(replies, &model.CommentReply{
			Style:   reply.Style,
			Content: reply.Content,
		})
	}

	return &model.CommentReplyResponse{
		Replies: replies,
	}, nil
}