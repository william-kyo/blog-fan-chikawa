package repository

import (
	"blog-fanchiikawa-service/db"
)

type textKeywordRepository struct{}

func NewTextKeywordRepository() TextKeywordRepository {
	return &textKeywordRepository{}
}

func (r *textKeywordRepository) Create(textKeyword *db.TextKeyword) error {
	_, err := db.Engine.Insert(textKeyword)
	return err
}

func (r *textKeywordRepository) GetByKeyword(keyword string) (*db.TextKeyword, error) {
	var result db.TextKeyword
	has, err := db.Engine.Where("keyword = ?", keyword).Get(&result)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &result, err
}

type imageTextKeywordRepository struct{}

func NewImageTextKeywordRepository() ImageTextKeywordRepository {
	return &imageTextKeywordRepository{}
}

func (r *imageTextKeywordRepository) Create(imageLabel *db.ImageTextKeyword) error {
	_, err := db.Engine.Insert(imageLabel)
	return err
}

func (r *imageTextKeywordRepository) GetByImageAndKeyword(imageID, textKeywordID int64) (*db.ImageTextKeyword, error) {
	var imageTextKeyword db.ImageTextKeyword
	has, err := db.Engine.Where("image_id = ? AND text_keyword_id = ?", imageID, textKeywordID).Get(&imageTextKeyword)
	if !has {
		return nil, nil
	}
	return &imageTextKeyword, err
}
