package resolver

import (
	"blog-fanchiikawa-service/graph/model"
	"blog-fanchiikawa-service/service"
	"context"
	"time"
)

func (r *Resolver) CreateChat(ctx context.Context, input model.CreateChatInput) (*model.Chat, error) {
	var botName, botId, botAlias, localeId string
	
	if input.BotName != nil {
		botName = *input.BotName
	}
	if input.BotID != nil {
		botId = *input.BotID
	}
	if input.BotAlias != nil {
		botAlias = *input.BotAlias
	}
	if input.LocaleID != nil {
		localeId = *input.LocaleID
	}

	req := &service.CreateChatRequest{
		UserID:   input.UserID,
		Title:    input.Title,
		BotName:  botName,
		BotId:    botId,
		BotAlias: botAlias,
		LocaleId: localeId,
	}

	chatResp, err := r.ChatService.CreateChat(req)
	if err != nil {
		return nil, err
	}

	createdAt, _ := time.Parse(time.RFC3339, chatResp.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, chatResp.UpdatedAt)

	return &model.Chat{
		ID:        chatResp.ID,
		UserID:    chatResp.UserID,
		Title:     chatResp.Title,
		BotName:   chatResp.BotName,
		SessionID: chatResp.SessionId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (r *Resolver) SendMessage(ctx context.Context, input model.SendMessageInput) (*model.ChatMessage, error) {
	req := &service.SendMessageRequest{
		ChatID:  input.ChatID,
		Message: input.Message,
	}

	msgResp, err := r.ChatService.SendMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	sentAt, _ := time.Parse(time.RFC3339, msgResp.SentAt)

	return &model.ChatMessage{
		ID:      msgResp.ID,
		ChatID:  msgResp.ChatID,
		Content: msgResp.Content,
		IsUser:  msgResp.IsUser,
		Intent:  &msgResp.Intent,
		SentAt:  sentAt,
	}, nil
}

func (r *Resolver) DeleteChat(ctx context.Context, chatID int64) (bool, error) {
	err := r.ChatService.DeleteChat(chatID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Resolver) UserChats(ctx context.Context, userID int64) ([]*model.Chat, error) {
	chats, err := r.ChatService.GetUserChats(userID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Chat, len(chats))
	for i, chat := range chats {
		createdAt, _ := time.Parse(time.RFC3339, chat.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, chat.UpdatedAt)

		result[i] = &model.Chat{
			ID:        chat.ID,
			UserID:    chat.UserID,
			Title:     chat.Title,
			BotName:   chat.BotName,
			SessionID: chat.SessionId,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return result, nil
}

func (r *Resolver) ChatHistory(ctx context.Context, chatID int64) (*model.ChatHistory, error) {
	history, err := r.ChatService.GetChatHistory(chatID)
	if err != nil {
		return nil, err
	}

	createdAt, _ := time.Parse(time.RFC3339, history.Chat.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, history.Chat.UpdatedAt)

	chat := &model.Chat{
		ID:        history.Chat.ID,
		UserID:    history.Chat.UserID,
		Title:     history.Chat.Title,
		BotName:   history.Chat.BotName,
		SessionID: history.Chat.SessionId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	messages := make([]*model.ChatMessage, len(history.Messages))
	for i, msg := range history.Messages {
		sentAt, _ := time.Parse(time.RFC3339, msg.SentAt)

		messages[i] = &model.ChatMessage{
			ID:      msg.ID,
			ChatID:  msg.ChatID,
			Content: msg.Content,
			IsUser:  msg.IsUser,
			Intent:  &msg.Intent,
			SentAt:  sentAt,
		}
	}

	return &model.ChatHistory{
		Chat:     chat,
		Messages: messages,
	}, nil
}

func (r *Resolver) LexConfig(ctx context.Context) (*model.LexConfig, error) {
	config := r.ConfigService.GetLexConfig()
	
	return &model.LexConfig{
		BotName:  config.BotName,
		BotID:    config.BotId,
		BotAlias: config.BotAlias,
		LocaleID: config.LocaleId,
	}, nil
}