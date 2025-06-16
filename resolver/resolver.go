package resolver

import (
	"blog-fanchiikawa-service/service"
)

// Resolver holds all the services needed for GraphQL resolvers
type Resolver struct {
	UserService      service.UserService
	LanguageService  service.LanguageService
	TranslateService service.TranslateService
	SpeechService    service.SpeechService
	StorageService   service.StorageService
}

// NewResolver creates a new Resolver instance with all services
func NewResolver(
	userService service.UserService,
	languageService service.LanguageService,
	translateService service.TranslateService,
	speechService service.SpeechService,
	storageService service.StorageService,
) *Resolver {
	return &Resolver{
		UserService:      userService,
		LanguageService:  languageService,
		TranslateService: translateService,
		SpeechService:    speechService,
		StorageService:   storageService,
	}
}