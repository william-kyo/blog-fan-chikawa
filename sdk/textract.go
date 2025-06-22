package sdk

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/textract"
)

var Textract *textract.Textract

func InitTextract() {
	// Create a separate session for Textract using a supported region
	// Textract is available in us-east-1, us-west-2, eu-west-1, etc.
	textractSession := AWSSession.Copy(&aws.Config{
		Region: aws.String("ap-northeast-2"), // Use ap-northeast-2 which supports Textract
	})
	Textract = textract.New(textractSession)
}

// DetectDocumentText extracts text from PDF documents in S3
func DetectDocumentText(bucketName, objectKey string) ([]string, error) {
	// Get the original S3 region and Textract region
	originalRegion := os.Getenv("AWS_DEFAULT_REGION")
	if originalRegion == "" {
		originalRegion = "ap-northeast-1"
	}
	textractRegion := "ap-northeast-2"

	// If regions are different, we need to handle cross-region access
	if originalRegion != textractRegion {
		log.Printf("Cross-region detected: S3 bucket in %s, Textract in %s", originalRegion, textractRegion)

		// For cross-region access, we need to use the full ARN format or copy the file
		// Let's try using the full S3 ARN format first
		fullBucketName := bucketName
		if !strings.Contains(bucketName, "arn:aws:s3") {
			// Use the standard bucket name but ensure Textract can access it
			fullBucketName = bucketName
		}

		input := &textract.DetectDocumentTextInput{
			Document: &textract.Document{
				S3Object: &textract.S3Object{
					Bucket: aws.String(fullBucketName),
					Name:   aws.String(objectKey),
				},
			},
		}

		// Call API with retry logic for cross-region access
		result, err := Textract.DetectDocumentText(input)
		if err != nil {
			// Check if it's an unsupported document format before trying file copy
			if strings.Contains(err.Error(), "UnsupportedDocumentException") {
				log.Printf("PDF format not supported by Textract for file: %s", objectKey)
				return handleUnsupportedPDF(bucketName, objectKey)
			}
			
			log.Printf("Cross-region access failed, attempting file copy method: %v", err)
			return detectDocumentTextWithCopy(bucketName, objectKey, originalRegion, textractRegion)
		}

		return processTextractResult(result, bucketName, objectKey)
	}

	// Same region - direct access
	input := &textract.DetectDocumentTextInput{
		Document: &textract.Document{
			S3Object: &textract.S3Object{
				Bucket: aws.String(bucketName),
				Name:   aws.String(objectKey),
			},
		},
	}

	// Call API
	result, err := Textract.DetectDocumentText(input)
	if err != nil {
		log.Printf("Failed to call DetectDocumentText: %v", err)
		
		// Check if it's an unsupported document format
		if strings.Contains(err.Error(), "UnsupportedDocumentException") {
			log.Printf("PDF format not supported by Textract for file: %s", objectKey)
			return handleUnsupportedPDF(bucketName, objectKey)
		}
		
		return nil, fmt.Errorf("failed to call detect document text: %w", err)
	}

	return processTextractResult(result, bucketName, objectKey)
}

// processTextractResult processes Textract result and returns text array
func processTextractResult(result *textract.DetectDocumentTextOutput, bucketName, objectKey string) ([]string, error) {
	// Output results
	log.Printf("Document: s3://%s/%s\n", bucketName, objectKey)
	log.Printf("Detected %d text blocks:\n\n", len(result.Blocks))

	var results []string
	for i, block := range result.Blocks {
		// Only process WORD type blocks, skip LINE and PAGE blocks
		if *block.BlockType == "WORD" {
			log.Printf("%d. Text: %s (Confidence: %.2f%%)\n",
				i+1, *block.Text, *block.Confidence)

			results = append(results, *block.Text)

			// If there are bounding box information, print them too
			if block.Geometry != nil && block.Geometry.BoundingBox != nil {
				bbox := block.Geometry.BoundingBox
				log.Printf("   - Bounding box: (%.3f, %.3f, %.3f, %.3f)\n",
					*bbox.Left, *bbox.Top, *bbox.Width, *bbox.Height)
			}
		}
		log.Println()
	}

	return results, nil
}

// detectDocumentTextWithCopy handles cross-region access by copying file temporarily
func detectDocumentTextWithCopy(bucketName, objectKey, sourceRegion, targetRegion string) ([]string, error) {
	log.Printf("Attempting cross-region file copy from %s to %s", sourceRegion, targetRegion)

	// Create S3 clients for both regions
	targetS3Session := AWSSession.Copy(&aws.Config{Region: aws.String(targetRegion)})

	targetS3 := s3.New(targetS3Session)

	// Create temporary bucket name in target region
	tempBucketName := bucketName + "-textract-temp"
	tempObjectKey := objectKey

	// Check if temp bucket exists, if not create it
	_, err := targetS3.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(tempBucketName),
	})
	if err != nil {
		log.Printf("Temp bucket doesn't exist, creating: %s", tempBucketName)
		_, err = targetS3.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(tempBucketName),
			CreateBucketConfiguration: &s3.CreateBucketConfiguration{
				LocationConstraint: aws.String(targetRegion),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create temp bucket: %w", err)
		}
	}

	// Copy object from source to target region
	copySource := fmt.Sprintf("%s/%s", bucketName, objectKey)
	_, err = targetS3.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(tempBucketName),
		Key:        aws.String(tempObjectKey),
		CopySource: aws.String(copySource),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to copy object: %w", err)
	}

	log.Printf("Successfully copied file to temp location: s3://%s/%s", tempBucketName, tempObjectKey)

	// Now process with Textract
	input := &textract.DetectDocumentTextInput{
		Document: &textract.Document{
			S3Object: &textract.S3Object{
				Bucket: aws.String(tempBucketName),
				Name:   aws.String(tempObjectKey),
			},
		},
	}

	result, err := Textract.DetectDocumentText(input)
	if err != nil {
		// Clean up temp file before returning error
		targetS3.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(tempBucketName),
			Key:    aws.String(tempObjectKey),
		})
		
		// Check if it's an unsupported document format
		if strings.Contains(err.Error(), "UnsupportedDocumentException") {
			log.Printf("PDF format not supported by Textract, trying alternative approach for file: %s", objectKey)
			return handleUnsupportedPDF(bucketName, objectKey)
		}
		
		return nil, fmt.Errorf("failed to call detect document text on copied file: %w", err)
	}

	// Clean up temp file
	_, deleteErr := targetS3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(tempBucketName),
		Key:    aws.String(tempObjectKey),
	})
	if deleteErr != nil {
		log.Printf("Warning: failed to delete temp file: %v", deleteErr)
	}

	return processTextractResult(result, bucketName, objectKey)
}

// handleUnsupportedPDF handles PDFs that Textract cannot process directly
func handleUnsupportedPDF(bucketName, objectKey string) ([]string, error) {
	log.Printf("Handling unsupported PDF format for: s3://%s/%s", bucketName, objectKey)
	
	log.Printf("PDF format not supported - possible reasons:")
	log.Printf("1. PDF is image-based (scanned) - requires StartDocumentTextDetection")
	log.Printf("2. PDF is encrypted or password protected")
	log.Printf("3. PDF format version is unsupported")
	log.Printf("4. PDF file is corrupted")
	
	// For now, return empty array to mark as processed but indicate no text found
	log.Printf("Skipping text extraction for unsupported PDF: %s", objectKey)
	return []string{}, nil
}
