package repository

import (
	"blog-fanchiikawa-service/db"
)

type labelRepository struct{}

func NewLabelRepository() LabelRepository {
	return &labelRepository{}
}

func (r *labelRepository) Create(label *db.Label) error {
	_, err := db.Engine.Insert(label)
	return err
}

func (r *labelRepository) GetByName(labelName string) (*db.Label, error) {
	var result db.Label
	has, err := db.Engine.Where("name = ?", labelName).Get(&result)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &result, err
}
