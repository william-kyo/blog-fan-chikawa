package service

import (
	"blog-fanchiikawa-service/db"
	"blog-fanchiikawa-service/repository"
)

type MediaService interface {
	CreateImage(filename string, originFilename string, uploaded bool) error
}

type mediaService struct {
	imageRepo repository.ImageRepository
}

func NewMediaService(imageRepo repository.ImageRepository) MediaService {
	return &mediaService{
		imageRepo: imageRepo,
	}
}

func (s *mediaService) CreateImage(filename string, originFilename string, uploaded bool) error {
	newImage := &db.Image{
		Filename:       filename,
		OriginFilename: originFilename,
		Uploaded:       uploaded,
	}
	err := s.imageRepo.Create(newImage)
	return err
}
