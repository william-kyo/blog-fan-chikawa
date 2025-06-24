# Project Long-term TODO List

## 🎯 Learning Objectives
Based on AWS AI hands-on learning plan, master major AWS AI services using Go language, build multiple small runnable projects, and establish a solid practical foundation for AIF-C01 certification.

---

## 📅 Week 3: Image Recognition Practice

### Task Goals
- [x] Master Rekognition & Textract

### Specific Tasks
- [x] Write image upload schedule task
- [x] Call Rekognition for text detection and label recognition
- [x] Call Textract for OCR extraction
- [x] Implement file extension tracking in database
- [x] Create intelligent file processing based on format (images vs PDFs)
- [x] Handle cross-region AWS service access
- [x] Implement comprehensive error handling for unsupported formats


---

## 📅 Week 4: Chatbot Practice

### Task Goals
- [ ] Quickly build chatbot prototype

### Specific Tasks
- [x] Create Lex bot in AWS Console
- [ ] Implement chat client using Go SDK
- [ ] Integrate Web applicatin

---

## 📅 Week 5: Lightweight Model Deployment Practice

### Task Goals
- [ ] Familiarize with SageMaker deployment and inference process

### Specific Tasks
- [ ] Deploy built-in models using JumpStart
- [ ] Call SageMaker Endpoint with Go HTTP Client
- [ ] Simple form upload CSV for inference requests

---

## 🔧 Technical Improvement Plans

### Architecture Optimization
- [ ] Add user authentication and authorization system
- [ ] Implement caching layer (Redis)
- [ ] Add API rate limiting functionality
- [ ] Implement file upload feature optimization

### Quality Assurance
- [ ] Add unit test coverage
- [ ] Implement CI/CD pipeline
- [ ] Add monitoring and logging system
- [ ] Performance optimization and stress testing

### Documentation Enhancement
- [ ] API documentation generation
- [ ] Deployment guide
- [ ] Developer guide
- [ ] Architecture design documentation

---

## 🏆 Completed Projects

### ✅ Week 1-2 Achievements
- [x] Configure Go development environment with AWS SDK
- [x] Implement GraphQL service integrating multiple AWS AI services:
  - [x] Comprehend: sentiment analysis, language detection
  - [x] Translate: translation
  - [x] Polly: text-to-speech
- [x] Implement layered architecture refactoring
- [x] Migrate to XORM ORM
- [x] Clean up project structure

### ✅ Week 3 Achievements
- [x] Complete image recognition and text extraction system:
  - [x] Rekognition: image label detection and text recognition
  - [x] Textract: PDF document text extraction with cross-region support
  - [x] Scheduler system for automated file processing
  - [x] Database models for images, labels, text keywords and relationships
  - [x] Multi-format file processing (images: JPG/PNG, documents: PDF)
  - [x] Intelligent error handling for unsupported file formats
  - [x] File extension tracking and type-based routing
  - [x] Cross-region AWS service access with fallback mechanisms

---

## 📚 Learning Resources

- AWS Go SDK Examples: https://aws.github.io/aws-sdk-go-v2/docs/code-examples/
- AWS Console Free Tier: https://aws.amazon.com/free/
- AWS Official Labs: https://explore.skillbuilder.aws/learn

---

## 🎯 Next Focus

Current focus is Week 4 task: **Chatbot Practice**
- Move to conversational AI implementation
- Integrate Amazon Lex for natural language understanding
- Build interactive chat interfaces