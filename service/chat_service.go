package service

import (
	"blog-fanchiikawa-service/db"
	"blog-fanchiikawa-service/repository"
	"blog-fanchiikawa-service/sdk"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
)

// ChatService defines the interface for chat business logic
type ChatService interface {
	CreateChat(req *CreateChatRequest) (*ChatResponse, error)
	SendMessage(ctx context.Context, req *SendMessageRequest) (*MessageResponse, error)
	GetChatHistory(chatID int64) (*ChatHistoryResponse, error)
	GetUserChats(userID int64) ([]*ChatResponse, error)
	DeleteChat(chatID int64) error
}

// chatService implements ChatService interface
type chatService struct {
	chatRepo        repository.ChatRepository
	chatMessageRepo repository.ChatMessageRepository
	lexService      *sdk.LexService
	snowflakeNode   *snowflake.Node
}

func NewChatService(chatRepo repository.ChatRepository, chatMessageRepo repository.ChatMessageRepository, lexService *sdk.LexService) ChatService {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatal("Failed to create snowflake node:", err)
	}

	return &chatService{
		chatRepo:        chatRepo,
		chatMessageRepo: chatMessageRepo,
		lexService:      lexService,
		snowflakeNode:   node,
	}
}

type CreateChatRequest struct {
	UserID   int64  `json:"userId"`
	Title    string `json:"title"`
	BotName  string `json:"botName,omitempty"`
	BotId    string `json:"botId,omitempty"`
	BotAlias string `json:"botAlias,omitempty"`
	LocaleId string `json:"localeId,omitempty"`
}

type ChatResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userId"`
	Title     string `json:"title"`
	BotName   string `json:"botName"`
	SessionId string `json:"sessionId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type SendMessageRequest struct {
	ChatID  int64  `json:"chatId"`
	Message string `json:"message"`
}

type MessageResponse struct {
	ID       int64  `json:"id"`
	ChatID   int64  `json:"chatId"`
	Content  string `json:"content"`
	IsUser   bool   `json:"isUser"`
	Intent   string `json:"intent"`
	SentAt   string `json:"sentAt"`
}

type ChatHistoryResponse struct {
	Chat     *ChatResponse      `json:"chat"`
	Messages []*MessageResponse `json:"messages"`
}

func (s *chatService) CreateChat(req *CreateChatRequest) (*ChatResponse, error) {
	sessionId := s.snowflakeNode.Generate().String()

	// Use environment variables as defaults if not provided
	botName := req.BotName
	if botName == "" {
		botName = os.Getenv("AWS_LEX_BOT_NAME")
	}

	botId := req.BotId
	if botId == "" {
		botId = os.Getenv("AWS_LEX_BOT_ID")
	}

	botAlias := req.BotAlias
	if botAlias == "" {
		botAlias = os.Getenv("AWS_LEX_BOT_ALIAS")
		if botAlias == "" {
			botAlias = "TSTALIASID" // Default test alias
		}
	}

	localeId := req.LocaleId
	if localeId == "" {
		localeId = os.Getenv("AWS_LEX_LOCALE_ID")
		if localeId == "" {
			localeId = "en_US" // Default locale
		}
	}

	chat := &db.Chat{
		UserID:    req.UserID,
		Title:     req.Title,
		BotName:   botName,
		BotId:     botId,
		BotAlias:  botAlias,
		LocaleId:  localeId,
		SessionId: sessionId,
	}

	if err := s.chatRepo.CreateChat(chat); err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	return &ChatResponse{
		ID:        chat.ID,
		UserID:    chat.UserID,
		Title:     chat.Title,
		BotName:   chat.BotName,
		SessionId: chat.SessionId,
		CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chat.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *chatService) SendMessage(ctx context.Context, req *SendMessageRequest) (*MessageResponse, error) {
	chat, err := s.chatRepo.GetChatByID(req.ChatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}
	if chat == nil {
		return nil, fmt.Errorf("chat not found")
	}

	userMessage := &db.ChatMessage{
		ChatID:  req.ChatID,
		Content: req.Message,
		IsUser:  true,
	}

	if err := s.chatMessageRepo.CreateMessage(userMessage); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	lexReq := &sdk.LexRequest{
		BotId:      chat.BotId,
		BotAliasId: chat.BotAlias,
		LocaleId:   chat.LocaleId,
		SessionId:  chat.SessionId,
		Text:       req.Message,
	}

	lexResp, err := s.lexService.RecognizeText(ctx, lexReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get Lex response: %w", err)
	}

	botResponse := ""
	if len(lexResp.Messages) > 0 {
		botResponse = lexResp.Messages[0]
	}

	botMessage := &db.ChatMessage{
		ChatID:  req.ChatID,
		Content: botResponse,
		IsUser:  false,
		Intent:  lexResp.IntentName,
	}

	if err := s.chatMessageRepo.CreateMessage(botMessage); err != nil {
		return nil, fmt.Errorf("failed to save bot message: %w", err)
	}

	return &MessageResponse{
		ID:      botMessage.ID,
		ChatID:  botMessage.ChatID,
		Content: botMessage.Content,
		IsUser:  botMessage.IsUser,
		Intent:  botMessage.Intent,
		SentAt:  botMessage.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *chatService) GetChatHistory(chatID int64) (*ChatHistoryResponse, error) {
	chat, err := s.chatRepo.GetChatByID(chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}
	if chat == nil {
		return nil, fmt.Errorf("chat not found")
	}

	messages, err := s.chatMessageRepo.GetMessagesByChatID(chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	chatResp := &ChatResponse{
		ID:        chat.ID,
		UserID:    chat.UserID,
		Title:     chat.Title,
		BotName:   chat.BotName,
		SessionId: chat.SessionId,
		CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chat.UpdatedAt.Format(time.RFC3339),
	}

	messageResps := make([]*MessageResponse, len(messages))
	for i, msg := range messages {
		messageResps[i] = &MessageResponse{
			ID:      msg.ID,
			ChatID:  msg.ChatID,
			Content: msg.Content,
			IsUser:  msg.IsUser,
			Intent:  msg.Intent,
			SentAt:  msg.CreatedAt.Format(time.RFC3339),
		}
	}

	return &ChatHistoryResponse{
		Chat:     chatResp,
		Messages: messageResps,
	}, nil
}

func (s *chatService) GetUserChats(userID int64) ([]*ChatResponse, error) {
	chats, err := s.chatRepo.GetChatsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user chats: %w", err)
	}

	responses := make([]*ChatResponse, len(chats))
	for i, chat := range chats {
		responses[i] = &ChatResponse{
			ID:        chat.ID,
			UserID:    chat.UserID,
			Title:     chat.Title,
			BotName:   chat.BotName,
			SessionId: chat.SessionId,
			CreatedAt: chat.CreatedAt.Format(time.RFC3339),
			UpdatedAt: chat.UpdatedAt.Format(time.RFC3339),
		}
	}

	return responses, nil
}

func (s *chatService) DeleteChat(chatID int64) error {
	if err := s.chatMessageRepo.DeleteMessagesByChatID(chatID); err != nil {
		return fmt.Errorf("failed to delete chat messages: %w", err)
	}

	if err := s.chatRepo.DeleteChat(chatID); err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}