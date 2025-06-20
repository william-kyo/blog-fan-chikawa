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

func (r *imageRepository) GetByID(id int64) (*db.Image, error) {
	var image db.Image
	has, err := db.Engine.Where("id = ?", id).Get(&image)
	if !has {
		return nil, nil
	}
	return &image, err
}

func (r *imageRepository) UpdateLabelDetected(id int64, labelDetected bool) (int64, error) {
	affected, err := db.Engine.ID(id).Cols("label_detected").Update(&db.Image{LabelDetected: labelDetected})
	return affected, err
}

func (r *imageRepository) GetByLabelDetected(labelDetected bool) ([]*db.Image, error) {
	var images []*db.Image
	err := db.Engine.Where("label_detected = ?", labelDetected).Find(&images)
	return images, err
}
