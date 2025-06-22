package service

import (
	"blog-fanchiikawa-service/db"
	"blog-fanchiikawa-service/repository"
	"blog-fanchiikawa-service/sdk"
	"errors"
	"fmt"
	"log"
	"strings"
)

type MediaService interface {
	CreateImage(filename, originFilename, fileExtension, bucket, objectKey string, uploaded bool) error
	DetectAndSaveImageLabels() error
	DetectAndSaveImageText() error
}

type mediaService struct {
	imageRepo            repository.ImageRepository
	labelRepo            repository.LabelRepository
	imageLabelRepo       repository.ImageLabelRepository
	textKeywordRepo      repository.TextKeywordRepository
	imageTextKeywordRepo repository.ImageTextKeywordRepository
	transactionMgr       repository.TransactionManager
}

func NewMediaService(imageRepo repository.ImageRepository, labelRepo repository.LabelRepository, imageLabelRepo repository.ImageLabelRepository, textKeywordRepo repository.TextKeywordRepository, imageTextKeywordRepo repository.ImageTextKeywordRepository, transactionMgr repository.TransactionManager) MediaService {
	return &mediaService{
		imageRepo:            imageRepo,
		labelRepo:            labelRepo,
		imageLabelRepo:       imageLabelRepo,
		textKeywordRepo:      textKeywordRepo,
		imageTextKeywordRepo: imageTextKeywordRepo,
		transactionMgr:       transactionMgr,
	}
}

func (s *mediaService) CreateImage(filename, originFilename, fileExtension, bucket, objectKey string, uploaded bool) error {
	newImage := &db.Image{
		Filename:       filename,
		OriginFilename: originFilename,
		FileExtension:  fileExtension,
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

func (s *mediaService) DetectAndSaveImageText() error {
	undetectedImages, err := s.imageRepo.GetByTextDetected(false)
	if err != nil {
		return fmt.Errorf("failed to fetch undetected image text")
	}

	log.Printf("Found %d undetected files to process text", len(undetectedImages))
	for _, image := range undetectedImages {
		log.Printf("Processing file ID: %d, Extension: %s, Bucket: %s, ObjectKey: %s", 
			image.ID, image.FileExtension, image.Bucket, image.ObjectKey)

		var textKeywords []string
		var err error

		// Determine processing method based on file extension
		fileExt := strings.ToLower(image.FileExtension)
		if fileExt == ".pdf" {
			// Use Textract for PDF files
			log.Printf("Using Textract for PDF file ID: %d", image.ID)
			textKeywords, err = sdk.DetectDocumentText(image.Bucket, image.ObjectKey)
		} else if fileExt == ".jpg" || fileExt == ".jpeg" || fileExt == ".png" || fileExt == ".gif" || fileExt == ".bmp" {
			// Use Rekognition for image files
			log.Printf("Using Rekognition for image file ID: %d", image.ID)
			textKeywords, err = sdk.DetectText(image.Bucket, image.ObjectKey)
		} else {
			log.Printf("Unsupported file extension %s for file ID: %d, skipping", fileExt, image.ID)
			// Mark as processed even though we skipped it to avoid reprocessing
			_, updateErr := s.imageRepo.UpdateTextDetected(image.ID, true)
			if updateErr != nil {
				log.Printf("Failed to mark unsupported file as processed: %v", updateErr)
			}
			continue
		}

		if err != nil {
			log.Printf("Failed to detect text for file ID %d: %v", image.ID, err)
			continue // Continue processing other files instead of returning error directly
		}

		err = s.SaveImageTextKeywords(image.ID, textKeywords)
		if err != nil {
			log.Printf("Failed to save text for file ID %d: %v", image.ID, err)
			continue // Continue processing other files
		}
	}

	return nil
}

func (s *mediaService) SaveImageTextKeywords(id int64, textKeywords []string) error {
	log.Printf("Starting SaveImageTextKeywords for image ID: %d with %d textKeywords: %v", id, len(textKeywords), textKeywords)

	// Remove duplicate textKeywords
	uniqueLabels := make(map[string]bool)
	var deduplicatedLabels []string
	for _, label := range textKeywords {
		if !uniqueLabels[label] {
			uniqueLabels[label] = true
			deduplicatedLabels = append(deduplicatedLabels, label)
		}
	}
	textKeywords = deduplicatedLabels
	log.Printf("After deduplication: %d unique textKeywords: %v", len(textKeywords), textKeywords)

	if len(textKeywords) == 0 {
		log.Printf("No textKeywords detected for image ID: %d, marking as detected", id)
		affected, err := s.imageRepo.UpdateTextDetected(id, true)
		if affected == 0 {
			log.Printf("Failed to update image ID: %d", id)
		}
		return err
	}

	err := s.transactionMgr.WithTransaction(func() error {
		for _, keyword := range textKeywords {
			log.Printf("Processing textKeywords '%s' for image ID: %d", keyword, id)

			// Check if keyword exists
			k, err := s.textKeywordRepo.GetByKeyword(keyword)
			if err != nil {
				log.Printf("Error getting keyword '%s': %v", keyword, err)
				return err
			}

			// If keyword doesn't exist, create new keyword
			if k == nil {
				log.Printf("Keyword '%s' not found, creating new keyword", keyword)
				newKeyword := &db.TextKeyword{Keyword: keyword}
				err = s.textKeywordRepo.Create(newKeyword)
				if err != nil {
					log.Printf("Failed to create label '%s': %v", keyword, err)
					return err
				}
				k = newKeyword
			}

			// Check if same image-keyword relationship already exists
			existingRelation, err := s.imageTextKeywordRepo.GetByImageAndKeyword(id, k.ID)
			if err != nil {
				log.Printf("Error checking existing image-keyword relationship: %v", err)
				return err
			}
			if existingRelation != nil {
				log.Printf("Image-keyword relationship already exists: ImageID=%d, KeyWordID=%d, skipping", id, k.ID)
				continue // Skip existing relationship
			}

			log.Printf("Creating image-keyword relationship: ImageID=%d, KeyWordID=%d", id, k.ID)
			err = s.imageTextKeywordRepo.Create(&db.ImageTextKeyword{
				ImageID:       id,
				TextKeywordId: k.ID,
			})
			if err != nil {
				log.Printf("Failed to create image-keyword relationship: %v", err)
				return err
			}
		}

		// After all keywords are saved, mark image as detected
		log.Printf("All textKeywords saved, marking image ID: %d as textKeywords detected", id)
		affected, err := s.imageRepo.UpdateTextDetected(id, true)
		if affected == 0 {
			log.Printf("Warning: UpdateTextDetected affected 0 rows for image ID: %d", id)
			return errors.New("failed to update image text_detected status")
		}
		if err != nil {
			log.Printf("Error updating text_detected status for image ID: %d: %v", id, err)
			return err
		}

		log.Printf("Successfully saved %d textKeywords for image ID: %d", len(textKeywords), id)
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed for image ID: %d: %v", id, err)
	} else {
		log.Printf("SaveImageLabels completed successfully for image ID: %d", id)
	}

	return err
}
