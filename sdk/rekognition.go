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
