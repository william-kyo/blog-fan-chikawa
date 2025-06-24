package service

import (
	"os"
)

// ConfigService defines the interface for configuration business logic
type ConfigService interface {
	GetLexConfig() *LexConfig
}

type LexConfig struct {
	BotName  string `json:"botName"`
	BotId    string `json:"botId"`
	BotAlias string `json:"botAlias"`
	LocaleId string `json:"localeId"`
}

// configService implements ConfigService interface
type configService struct{}

func NewConfigService() ConfigService {
	return &configService{}
}

func (s *configService) GetLexConfig() *LexConfig {
	return &LexConfig{
		BotName:  os.Getenv("AWS_LEX_BOT_NAME"),
		BotId:    os.Getenv("AWS_LEX_BOT_ID"),
		BotAlias: getEnvWithDefault("AWS_LEX_BOT_ALIAS", "TSTALIASID"),
		LocaleId: getEnvWithDefault("AWS_LEX_LOCALE_ID", "en_US"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}