package repository

import (
	"blog-fanchiikawa-service/db"
)

type imageRepository struct{}

func NewImageReposity() ImageRepository {
	return &imageRepository{}
}

func (r *imageRepository) Create(image *db.Image) error {
	_, err := db.Engine.Insert(image)
	return err
}

func (r *imageRepository) GetByID(id int64) ([]*db.Image, error) {
	var image []*db.Image
	err := db.Engine.Where("id = ?", id).Find(&image)
	return image, err
}
