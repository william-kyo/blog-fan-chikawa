# Custom Labels Detection Feature

## Overview

This feature allows users to upload images to S3 and perform custom object detection using AWS Rekognition Custom Labels. The system returns the top 2 detected labels with confidence scores.

## Architecture

### Option 1: Direct S3 Upload (Recommended)
```
Frontend → GraphQL (Get Upload URL) → S3 (Direct Upload) → GraphQL (Detect via S3 Key) → Rekognition
```

### Option 2: Server Upload (Legacy)
```
Frontend → GraphQL Resolver → Custom Labels Service → SDK Layer → AWS S3 + Rekognition
```

### Layer Responsibilities

- **Frontend**: Direct S3 upload using presigned URLs, progress tracking
- **GraphQL Resolver**: Parameter conversion and validation
- **Service Layer**: Business logic, S3 URL generation, and AWS integration
- **SDK Layer**: Low-level AWS API calls and utilities

## Components

### 1. Backend Components

#### SDK Layer (`sdk/`)
- **rekognition.go**: Added `DetectCustomLabels()` function and `CustomLabelResult` struct
- **s3.go**: Added `UploadFileForCustomLabels()` and `GeneratePresignedURL()` functions

#### Service Layer (`service/`)
- **custom_labels_service.go**: Business logic for upload and detection workflow
  - Input validation and file processing
  - AWS service integration (S3 + Rekognition)
  - Internal file wrapper for multipart.File interface
  - Error handling and response formatting

#### GraphQL Layer
- **schema.graphqls**: Added `CustomLabel`, `CustomLabelsResult` types and `uploadAndDetectCustomLabels` mutation
- **resolver/mutation.go**: Lightweight resolver focused on parameter conversion and service calls

### 2. Frontend Components

#### Web Interface
- **web/custom-labels.html**: Complete Vue.js SPA with:
  - Drag-and-drop file upload
  - File validation (type and size)
  - Image preview
  - Results display with confidence scores
  - Modern, responsive UI design

## Configuration

### Environment Variables

Add these variables to your `.env` file:

```env
# AWS Rekognition Custom Labels Configuration
REKOGNITION_S3_BUCKET=your-s3-bucket-name
REKOGNITION_PROJECT_VERSION_ARN=arn:aws:rekognition:region:account:project/project-name/version/version-name/timestamp
```

### AWS Setup Requirements

1. **S3 Bucket**: Create a bucket for storing uploaded images
2. **Rekognition Custom Labels Project**: Train and deploy a custom model
3. **IAM Permissions**: Ensure your AWS credentials have access to:
   - S3: PutObject, GetObject
   - Rekognition: DetectCustomLabels

## API Usage

### Option 1: Direct S3 Upload (Recommended)

#### Step 1: Get Upload URL
```graphql
query GenerateS3UploadUrl($filename: String!) {
  generateS3UploadUrl(filename: $filename) {
    uploadUrl
    key
    fields {
      name
      value
    }
  }
}
```

#### Step 2: Upload to S3 (Frontend)
```javascript
// Direct PUT request to S3
await axios.put(uploadData.uploadUrl, file, {
  headers: { 'Content-Type': file.type }
});
```

#### Step 3: Detect Labels
```graphql
mutation DetectCustomLabelsFromS3($input: DetectCustomLabelsInput!) {
  detectCustomLabelsFromS3(input: $input) {
    imageUrl
    s3Key
    labels {
      name
      confidence
    }
  }
}
```

### Option 2: Server Upload (Legacy)

```graphql
mutation UploadAndDetectCustomLabels($file: Upload!) {
  uploadAndDetectCustomLabels(file: $file) {
    imageUrl
    s3Key
    labels {
      name
      confidence
    }
  }
}
```

### Response Format

```json
{
  "data": {
    "uploadAndDetectCustomLabels": {
      "imageUrl": "https://s3.amazonaws.com/bucket/presigned-url",
      "s3Key": "custom-labels/1640995200_image.jpg",
      "labels": [
        {
          "name": "Cat",
          "confidence": 95.67
        },
        {
          "name": "Animal",
          "confidence": 89.23
        }
      ]
    }
  }
}
```

## Features

### Frontend Features
- **File Upload**: Drag-and-drop or click to upload
- **Direct S3 Upload**: Client-side upload using presigned URLs
- **File Validation**: Type checking (images only) and size limits (10MB)
- **Progress Tracking**: Step-by-step progress indication (URL → Upload → Detect → Complete)
- **Results Display**: Image preview with detected labels and confidence scores
- **Error Handling**: User-friendly error messages

### Backend Features
- **Presigned URL Generation**: Secure, time-limited upload URLs
- **S3 Integration**: Secure file management with unique naming
- **Custom Labels Detection**: AWS Rekognition integration via S3 keys
- **Top Results**: Returns only the top 2 highest confidence labels
- **Presigned URLs**: Secure image viewing with time-limited access
- **Dual Upload Methods**: Both direct S3 and server upload supported
- **Type Safety**: Full GraphQL type definitions

## File Structure

```
blog-fanchiikawa-service/
├── sdk/
│   ├── rekognition.go          # Custom labels detection
│   └── s3.go                   # File upload and URL generation
├── service/
│   └── custom_labels_service.go # Business logic
├── resolver/
│   └── mutation.go             # GraphQL resolver with file upload
├── web/
│   └── custom-labels.html      # Frontend interface
├── graph/
│   └── schema.graphqls         # GraphQL schema
└── docs/
    └── custom-labels-feature.md # This documentation
```

## Usage

1. **Start the Server**:
   ```bash
   go run server.go
   ```

2. **Access the Interface**:
   ```
   http://localhost:8080/custom-labels/custom-labels.html
   ```

3. **Upload and Detect**:
   - Select or drag an image file
   - Click "Detect Labels"
   - View results with confidence scores

## Technical Details

### Image Processing Flow
1. File validation (client-side)
2. Upload to GraphQL endpoint (multipart/form-data)
3. S3 upload with unique key generation
4. Rekognition Custom Labels API call
5. Sort results by confidence (top 2)
6. Generate presigned URL for image display
7. Return structured response

### Error Handling
- Client-side validation for file type and size
- Server-side error handling for AWS API failures
- User-friendly error messages
- Graceful degradation

## Security

- **Private S3 Access**: Files uploaded with private ACL
- **Presigned URLs**: Time-limited access (1 hour expiration)
- **File Validation**: Type and size restrictions
- **Environment Variables**: Sensitive configuration externalized

## Performance

- **Unique Naming**: Timestamp-based keys prevent conflicts
- **Efficient Sorting**: Simple bubble sort for small result sets (max 10 labels)
- **Presigned URLs**: Direct S3 access without server proxy
- **Top Results Only**: Returns only top 2 labels to reduce response size