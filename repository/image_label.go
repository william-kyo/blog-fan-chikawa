package repository

import "blog-fanchiikawa-service/db"

type imageLabelRepository struct{}

func NewImageLabelRepository() ImageLabelRepository {
	return &imageLabelRepository{}
}

func (r *imageLabelRepository) Create(imageLabel *db.ImageLabel) error {
	_, err := db.Engine.Insert(imageLabel)
	return err
}

func (r *imageLabelRepository) GetByImageAndLabel(imageID, labelID int64) (*db.ImageLabel, error) {
	var imageLabel db.ImageLabel
	has, err := db.Engine.Where("image_id = ? AND label_id = ?", imageID, labelID).Get(&imageLabel)
	if !has {
		return nil, nil
	}
	return &imageLabel, err
}
