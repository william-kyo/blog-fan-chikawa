package repository

import (
	"blog-fanchiikawa-service/db"
	"xorm.io/xorm"
)

type ChatMessageRepository interface {
	CreateMessage(message *db.ChatMessage) error
	GetMessagesByChatID(chatID int64) ([]*db.ChatMessage, error)
	GetRecentMessagesByChatID(chatID int64, limit int) ([]*db.ChatMessage, error)
	DeleteMessagesByChatID(chatID int64) error
}

type chatMessageRepository struct {
	engine *xorm.Engine
}

func NewChatMessageRepository(engine *xorm.Engine) ChatMessageRepository {
	return &chatMessageRepository{
		engine: engine,
	}
}

func (r *chatMessageRepository) CreateMessage(message *db.ChatMessage) error {
	_, err := r.engine.Insert(message)
	return err
}

func (r *chatMessageRepository) GetMessagesByChatID(chatID int64) ([]*db.ChatMessage, error) {
	var messages []*db.ChatMessage
	err := r.engine.Where("chat_id = ?", chatID).OrderBy("created_at ASC").Find(&messages)
	return messages, err
}

func (r *chatMessageRepository) GetRecentMessagesByChatID(chatID int64, limit int) ([]*db.ChatMessage, error) {
	var messages []*db.ChatMessage
	err := r.engine.Where("chat_id = ?", chatID).OrderBy("created_at DESC").Limit(limit).Find(&messages)
	return messages, err
}

func (r *chatMessageRepository) DeleteMessagesByChatID(chatID int64) error {
	_, err := r.engine.Where("chat_id = ?", chatID).Delete(&db.ChatMessage{})
	return err
}