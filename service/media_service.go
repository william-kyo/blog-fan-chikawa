package service

import (
	"blog-fanchiikawa-service/db"
	"blog-fanchiikawa-service/repository"
	"blog-fanchiikawa-service/sdk"
	"errors"
	"fmt"
	"log"
)

type MediaService interface {
	CreateImage(filename, originFilename, bucket, objectKey string, uploaded bool) error
	DetectAndSaveImageLabels() error
}

type mediaService struct {
	imageRepo      repository.ImageRepository
	labelRepo      repository.LabelRepository
	imageLabelRepo repository.ImageLabelRepository
	transactionMgr repository.TransactionManager
}

func NewMediaService(imageRepo repository.ImageRepository, labelRepo repository.LabelRepository, imageLabelRepo repository.ImageLabelRepository, transactionMgr repository.TransactionManager) MediaService {
	return &mediaService{
		imageRepo:      imageRepo,
		labelRepo:      labelRepo,
		imageLabelRepo: imageLabelRepo,
		transactionMgr: transactionMgr,
	}
}

func (s *mediaService) CreateImage(filename, originFilename, bucket, objectKey string, uploaded bool) error {
	newImage := &db.Image{
		Filename:       filename,
		OriginFilename: originFilename,
		Bucket:         bucket,
		ObjectKey:      objectKey,
		Uploaded:       uploaded,
	}
	err := s.imageRepo.Create(newImage)
	return err
}

func (s *mediaService) DetectAndSaveImageLabels() error {
	undetectedImages, err := s.imageRepo.GetByLabelDetected(false)
	if err != nil {
		return fmt.Errorf("failed to fetch undetected image")
	}

	log.Printf("Found %d undetected images to process", len(undetectedImages))

	for _, image := range undetectedImages {
		log.Printf("Processing image ID: %d, Bucket: %s, ObjectKey: %s", image.ID, image.Bucket, image.ObjectKey)

		labels, err := sdk.DetectLabels(image.Bucket, image.ObjectKey)
		if err != nil {
			log.Printf("Failed to detect labels for image ID %d: %v", image.ID, err)
			continue // Continue processing other images instead of returning error directly
		}

		err = s.SaveImageLabels(image.ID, labels)
		if err != nil {
			log.Printf("Failed to save labels for image ID %d: %v", image.ID, err)
			continue // Continue processing other images
		}
	}
	return nil
}

func (s *mediaService) SaveImageLabels(id int64, labels []string) error {
	log.Printf("Starting SaveImageLabels for image ID: %d with %d labels: %v", id, len(labels), labels)

	// Remove duplicate labels
	uniqueLabels := make(map[string]bool)
	var deduplicatedLabels []string
	for _, label := range labels {
		if !uniqueLabels[label] {
			uniqueLabels[label] = true
			deduplicatedLabels = append(deduplicatedLabels, label)
		}
	}
	labels = deduplicatedLabels
	log.Printf("After deduplication: %d unique labels: %v", len(labels), labels)

	if len(labels) == 0 {
		log.Printf("No labels detected for image ID: %d, marking as detected", id)
		affected, err := s.imageRepo.UpdateLabelDetected(id, true)
		if affected == 0 {
			log.Printf("Failed to update image ID: %d", id)
		}
		return err
	}

	err := s.transactionMgr.WithTransaction(func() error {
		for _, labelName := range labels {
			log.Printf("Processing label '%s' for image ID: %d", labelName, id)

			// Check if label exists
			label, err := s.labelRepo.GetByName(labelName)
			if err != nil {
				log.Printf("Error getting label '%s': %v", labelName, err)
				return err
			}

			// If label doesn't exist, create new label
			if label == nil {
				log.Printf("Label '%s' not found, creating new label", labelName)
				newLabel := &db.Label{Name: labelName}
				err = s.labelRepo.Create(newLabel)
				if err != nil {
					log.Printf("Failed to create label '%s': %v", labelName, err)
					return err
				}
				label = newLabel
			}

			// Check if same image-label relationship already exists
			existingRelation, err := s.imageLabelRepo.GetByImageAndLabel(id, label.ID)
			if err != nil {
				log.Printf("Error checking existing image-label relationship: %v", err)
				return err
			}
			if existingRelation != nil {
				log.Printf("Image-label relationship already exists: ImageID=%d, LabelID=%d, skipping", id, label.ID)
				continue // Skip existing relationship
			}

			log.Printf("Creating image-label relationship: ImageID=%d, LabelID=%d", id, label.ID)
			err = s.imageLabelRepo.Create(&db.ImageLabel{
				ImageID: id,
				LabelID: label.ID,
			})
			if err != nil {
				log.Printf("Failed to create image-label relationship: %v", err)
				return err
			}
		}

		// After all labels are saved, mark image as detected
		log.Printf("All labels saved, marking image ID: %d as label detected", id)
		affected, err := s.imageRepo.UpdateLabelDetected(id, true)
		if affected == 0 {
			log.Printf("Warning: UpdateLabelDetected affected 0 rows for image ID: %d", id)
			return errors.New("failed to update image label_detected status")
		}
		if err != nil {
			log.Printf("Error updating label_detected status for image ID: %d: %v", id, err)
			return err
		}

		log.Printf("Successfully saved %d labels for image ID: %d", len(labels), id)
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed for image ID: %d: %v", id, err)
	} else {
		log.Printf("SaveImageLabels completed successfully for image ID: %d", id)
	}

	return err
}
