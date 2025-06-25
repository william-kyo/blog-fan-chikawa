# AWS SDK Version Management Guide

## Why Do We Need Two AWS SDK Versions?

### AWS SDK v1 (`github.com/aws/aws-sdk-go`)
- **Purpose**: Existing services (S3, Comprehend, Translate, Polly, Rekognition, Textract)
- **Stability**: Mature and stable, widely used
- **Configuration Type**: `*session.Session`

### AWS SDK v2 (`github.com/aws/aws-sdk-go-v2`)
- **Purpose**: New services (Lex Runtime V2)
- **Reason**: Lex Runtime V2 API is only available in SDK v2
- **Configuration Type**: `aws.Config`
- **Advantages**: Better performance, more modern API design

## Unified Configuration Management

### Problems Before Optimization
```go
// Duplicated configuration code
func InitAWSSession() {
    region := os.Getenv("AWS_DEFAULT_REGION")
    profile := os.Getenv("AWS_PROFILE")
    // ... SDK v1 initialization
}

func InitAWSConfigV2() {
    region := os.Getenv("AWS_DEFAULT_REGION") // Duplicated
    profile := os.Getenv("AWS_PROFILE")       // Duplicated
    // ... SDK v2 initialization
}
```

### Optimized Solution
```go
// Unified configuration retrieval
func getAWSCredentials() (region, profile string) {
    // Unified reading from environment variables
}

// Unified initialization method
func InitAWS() {
    region, profile := getAWSCredentials()
    // Initialize both SDK v1 and v2
}

// Backward compatible methods
func InitAWSSession() { /* legacy support */ }
func InitAWSConfigV2() { /* legacy support */ }
```

## Usage

### Simplified Server Initialization
```go
func main() {
    // Single call to initialize all AWS configurations
    sdk.InitAWS()
    
    // Other service initialization...
}
```

### Using Different Versions in Services
```go
// Service using SDK v1
func NewS3Service() {
    session := sdk.GetAWSSession()
    return s3.New(session)
}

// Service using SDK v2
func NewLexService() {
    config := sdk.GetAWSConfig()
    return lexruntimev2.NewFromConfig(config)
}
```

## Future Plans

1. **Progressive Migration**: Gradually migrate existing services to SDK v2
2. **Unified SDK**: Remove SDK v1 dependency when all services are migrated
3. **Backward Compatibility**: Maintain legacy initialization methods to ensure existing code isn't broken

## Summary

This design solves the code duplication problem while maintaining backward compatibility and supporting new Lex Runtime V2 functionality. Through the unified `InitAWS()` method, AWS configuration management is simplified.