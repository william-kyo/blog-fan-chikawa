# GraphQL Text-to-Speech Interface Implementation

## Goal
Create a GraphQL interface that takes unknown language text, detects the language using AWS Comprehend, then generates speech using AWS Polly and stores the audio file in S3.

## Plan
- [x] Add AWS Polly SDK integration following existing patterns
- [x] Update GraphQL schema to add textToSpeech mutation
- [x] Implement GraphQL resolver for textToSpeech
- [x] Add textToSpeech input type to schema
- [x] Initialize Polly service in server.go
- [x] Test the complete workflow

## Current Understanding
- Existing pattern: AWS services initialized in server.go and wrapped in sdk/ directory
- S3 bucket: `fan-ai-warehouse`, region: `ap-northeast-1`
- Speech files will be stored in `speech/` folder
- Existing Comprehend detectLanguage can be reused for language detection
- GraphQL mutations follow input type pattern with simple string returns

## Implementation Details

### Input
- Type: String (unknown language text)

### Processing Flow
1. Call existing `detectLanguage` from Comprehend to get language code
2. If language detection fails, return error "Unable to detect language type"
3. Use language code and original text to call Polly for speech synthesis
4. Save generated audio file to S3 bucket `fan-ai-warehouse` in `speech/` folder
5. Return the S3 key of the audio file

### Output
- Type: String (S3 key of the generated audio file)

## Review

✅ **Implementation completed successfully**. The GraphQL text-to-speech interface has been implemented with the following features:

### Files Created/Modified:
1. **sdk/polly.go** - New AWS Polly SDK integration with language mapping and S3 upload
2. **graph/schema.graphqls** - Added TextToSpeech input type and textToSpeech mutation  
3. **graph/schema.resolvers.go** - Added TextToSpeech resolver with language detection
4. **server.go** - Added Polly service initialization

### Key Features Implemented:
- **Language Detection**: Uses existing AWS Comprehend detectLanguage to identify text language
- **Language Mapping**: Maps Comprehend language codes to Polly-supported language codes
- **Voice Selection**: Automatically selects appropriate voice based on detected language
- **S3 Storage**: Saves generated MP3 files to `fan-ai-warehouse` bucket in `speech/` folder
- **Error Handling**: Returns "Unable to detect language type" for detection failures
- **File Naming**: Uses timestamp and language code format: `speech/{timestamp}_{language}.mp3`

### Test Results:
- ✅ English text: "Hello world" → `speech/1750044647_en.mp3`
- ✅ Chinese text: "你好世界" → `speech/1750044656_zh.mp3`
- ✅ GraphQL interface working correctly through HTTP API

### GraphQL Usage:
```graphql
mutation {
  textToSpeech(input: {text: "Your text here"})
}
```

The implementation follows all existing codebase patterns and successfully integrates with the current AWS infrastructure.