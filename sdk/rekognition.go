package sdk

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

var Rekognition *rekognition.Rekognition

func InitRekognition() {
	Rekognition = rekognition.New(AWSSession)
}

// DetectLabels detects labels in S3 images
func DetectLabels(bucketName, objectKey string) ([]string, error) {
	// Build detect labels request
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(bucketName),
				Name:   aws.String(objectKey),
			},
		},
		MaxLabels:     aws.Int64(20),     // Return maximum 20 labels
		MinConfidence: aws.Float64(75.0), // Minimum confidence 75%
	}

	// Call API
	result, err := Rekognition.DetectLabels(input)
	if err != nil {
		log.Printf("Failed to call DetectLabels: %v", err)
		return nil, fmt.Errorf("failed to call detect labels: %w", err)
	}

	// Output results
	log.Printf("Image: s3://%s/%s\n", bucketName, objectKey)
	log.Printf("Detected %d labels:\n\n", len(result.Labels))

	var results []string
	for i, label := range result.Labels {
		log.Printf("%d. %s (Confidence: %.2f%%)\n",
			i+1, *label.Name, *label.Confidence)

		results = append(results, *label.Name)

		// If there are subcategories, print them too
		if len(label.Categories) > 0 {
			log.Printf("   Category: ")
			for j, category := range label.Categories {
				if j > 0 {
					log.Printf(", ")
				}
				log.Printf("%s", *category.Name)
				results = append(results, *category.Name)
			}
			log.Println()
		}

		// If there are instance information, print them too
		if len(label.Instances) > 0 {
			log.Printf("   Instance count: %d\n", len(label.Instances))
			for _, instance := range label.Instances {
				if instance.BoundingBox != nil {
					bbox := instance.BoundingBox
					log.Printf("   - Bounding box: (%.3f, %.3f, %.3f, %.3f) Confidence: %.2f%%\n",
						*bbox.Left, *bbox.Top, *bbox.Width, *bbox.Height,
						*instance.Confidence)
				}
			}
		}
		log.Println()
	}

	return results, nil
}

func DetectText(bucketName, objectKey string) ([]string, error) {
	// Build detect text request
	input := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(bucketName),
				Name:   aws.String(objectKey),
			},
		},
	}

	// Call API
	result, err := Rekognition.DetectText(input)
	if err != nil {
		log.Printf("Failed to call DetectText: %v", err)
		return nil, fmt.Errorf("failed to call detect text: %w", err)
	}

	// Output results
	log.Printf("Image: s3://%s/%s\n", bucketName, objectKey)
	log.Printf("Detected %d text detections:\n\n", len(result.TextDetections))

	var results []string
	for i, textDetection := range result.TextDetections {
		// Only process word-level detections, skip line-level
		if *textDetection.Type == "WORD" {
			log.Printf("%d. Text: %s (Confidence: %.2f%%)\n",
				i+1, *textDetection.DetectedText, *textDetection.Confidence)

			results = append(results, *textDetection.DetectedText)

			// If there are bounding box information, print them too
			if textDetection.Geometry != nil && textDetection.Geometry.BoundingBox != nil {
				bbox := textDetection.Geometry.BoundingBox
				log.Printf("   - Bounding box: (%.3f, %.3f, %.3f, %.3f)\n",
					*bbox.Left, *bbox.Top, *bbox.Width, *bbox.Height)
			}
		}
		log.Println()
	}

	return results, nil
}

// CustomLabelResult represents a custom label detection result
type CustomLabelResult struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

// DetectCustomLabels uses custom model to detect labels in S3 images
func DetectCustomLabels(bucketName, objectKey, projectVersionArn string) ([]CustomLabelResult, error) {
	// Build detect custom labels request
	input := &rekognition.DetectCustomLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(bucketName),
				Name:   aws.String(objectKey),
			},
		},
		ProjectVersionArn: aws.String(projectVersionArn),
		MaxResults:        aws.Int64(10),       // Return maximum 10 labels
		MinConfidence:     aws.Float64(50.0),  // Minimum confidence 50%
	}

	// Call API
	result, err := Rekognition.DetectCustomLabels(input)
	if err != nil {
		log.Printf("Failed to call DetectCustomLabels: %v", err)
		return nil, fmt.Errorf("failed to call detect custom labels: %w", err)
	}

	// Output results
	log.Printf("Image: s3://%s/%s\n", bucketName, objectKey)
	log.Printf("Detected %d custom labels:\n\n", len(result.CustomLabels))

	var results []CustomLabelResult
	for i, label := range result.CustomLabels {
		log.Printf("%d. %s (Confidence: %.2f%%)\n",
			i+1, *label.Name, *label.Confidence)

		results = append(results, CustomLabelResult{
			Name:       *label.Name,
			Confidence: *label.Confidence,
		})
	}

	// Sort by confidence (highest first) and return top 2
	if len(results) > 1 {
		// Simple bubble sort for small arrays
		for i := 0; i < len(results)-1; i++ {
			for j := 0; j < len(results)-1-i; j++ {
				if results[j].Confidence < results[j+1].Confidence {
					results[j], results[j+1] = results[j+1], results[j]
				}
			}
		}
	}

	// Return top 2 results
	if len(results) > 2 {
		results = results[:2]
	}

	return results, nil
}
