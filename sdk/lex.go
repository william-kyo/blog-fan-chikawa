package sdk

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lexruntimev2"
)

type LexService struct {
	client *lexruntimev2.Client
}

func NewLexService() *LexService {
	cfg := GetAWSConfig()
	return &LexService{
		client: lexruntimev2.NewFromConfig(cfg),
	}
}

type LexRequest struct {
	BotId      string
	BotAliasId string
	LocaleId   string
	SessionId  string
	Text       string
}

type LexResponse struct {
	Messages       []string
	IntentName     string
	DialogState    string
	SessionState   map[string]interface{}
	Interpretations []map[string]interface{}
}

func (l *LexService) RecognizeText(ctx context.Context, req *LexRequest) (*LexResponse, error) {
	input := &lexruntimev2.RecognizeTextInput{
		BotId:      &req.BotId,
		BotAliasId: &req.BotAliasId,
		LocaleId:   &req.LocaleId,
		SessionId:  &req.SessionId,
		Text:       &req.Text,
	}

	result, err := l.client.RecognizeText(ctx, input)
	if err != nil {
		return nil, err
	}

	response := &LexResponse{
		Messages: make([]string, 0),
	}

	if result.Messages != nil {
		for _, msg := range result.Messages {
			if msg.Content != nil {
				response.Messages = append(response.Messages, *msg.Content)
			}
		}
	}

	if result.SessionState != nil {
		if result.SessionState.Intent != nil && result.SessionState.Intent.Name != nil {
			response.IntentName = *result.SessionState.Intent.Name
		}
		
		if result.SessionState.DialogAction != nil {
			response.DialogState = string(result.SessionState.DialogAction.Type)
		}
	}

	if result.Interpretations != nil {
		response.Interpretations = make([]map[string]interface{}, 0)
		for _, interp := range result.Interpretations {
			interpMap := make(map[string]interface{})
			if interp.Intent != nil && interp.Intent.Name != nil {
				interpMap["intent"] = *interp.Intent.Name
			}
			if interp.NluConfidence != nil {
				interpMap["confidence"] = interp.NluConfidence.Score
			}
			response.Interpretations = append(response.Interpretations, interpMap)
		}
	}

	return response, nil
}