package repository

import (
	"blog-fanchiikawa-service/db"
	"xorm.io/xorm"
)

type ChatRepository interface {
	CreateChat(chat *db.Chat) error
	GetChatByID(id int64) (*db.Chat, error)
	GetChatsByUserID(userID int64) ([]*db.Chat, error)
	UpdateChat(chat *db.Chat) error
	DeleteChat(id int64) error
}

type chatRepository struct {
	engine *xorm.Engine
}

func NewChatRepository(engine *xorm.Engine) ChatRepository {
	return &chatRepository{
		engine: engine,
	}
}

func (r *chatRepository) CreateChat(chat *db.Chat) error {
	_, err := r.engine.Insert(chat)
	return err
}

func (r *chatRepository) GetChatByID(id int64) (*db.Chat, error) {
	chat := &db.Chat{}
	has, err := r.engine.ID(id).Get(chat)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return chat, nil
}

func (r *chatRepository) GetChatsByUserID(userID int64) ([]*db.Chat, error) {
	var chats []*db.Chat
	err := r.engine.Where("user_id = ?", userID).OrderBy("created_at DESC").Find(&chats)
	return chats, err
}

func (r *chatRepository) UpdateChat(chat *db.Chat) error {
	_, err := r.engine.ID(chat.ID).Update(chat)
	return err
}

func (r *chatRepository) DeleteChat(id int64) error {
	_, err := r.engine.ID(id).Delete(&db.Chat{})
	return err
}