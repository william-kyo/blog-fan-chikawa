package resolver

import (
	"blog-fanchiikawa-service/service"
)

// Resolver holds all the services needed for GraphQL resolvers
type Resolver struct {
	UserService         service.UserService
	LanguageService     service.LanguageService
	TranslateService    service.TranslateService
	SpeechService       service.SpeechService
	StorageService      service.StorageService
	ChatService         service.ChatService
	ConfigService       service.ConfigService
	CustomLabelsService service.CustomLabelsService
	CommentReplyService service.CommentReplyService
}

// NewResolver creates a new Resolver instance with all services
func NewResolver(
	userService service.UserService,
	languageService service.LanguageService,
	translateService service.TranslateService,
	speechService service.SpeechService,
	storageService service.StorageService,
	chatService service.ChatService,
	configService service.ConfigService,
	customLabelsService service.CustomLabelsService,
	commentReplyService service.CommentReplyService,
) *Resolver {
	return &Resolver{
		UserService:         userService,
		LanguageService:     languageService,
		TranslateService:    translateService,
		SpeechService:       speechService,
		StorageService:      storageService,
		ChatService:         chatService,
		ConfigService:       configService,
		CustomLabelsService: customLabelsService,
		CommentReplyService: commentReplyService,
	}
}