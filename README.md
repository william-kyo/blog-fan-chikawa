# Blog Fanchiikawa Service

A GraphQL-based microservice built with Go, featuring text-to-speech, language detection, translation, and user management capabilities powered by AWS services.

## 🏗️ Architecture

This project follows a clean layered architecture pattern following Go best practices:

```
┌─────────────────────────────────────────────────────────────┐
│                    Controller Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐                 │
│  │ GraphQL         │  │ Resolver        │                 │
│  │ Resolvers       │  │ (Thin Layer)    │                 │
│  └─────────────────┘  └─────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────────┐
│                     Service Layer                          │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐      │
│  │ User         │ │ Language     │ │ Speech       │      │
│  │ Service      │ │ Service      │ │ Service      │      │
│  └──────────────┘ └──────────────┘ └──────────────┘      │
│  ┌──────────────┐ ┌──────────────┐                       │
│  │ Translate    │ │ Storage      │                       │
│  │ Service      │ │ Service      │                       │
│  └──────────────┘ └──────────────┘                       │
│                               │                           │
│           AWS SDK Integration │                           │
└─────────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐                 │
│  │ User            │  │ Transaction     │                 │
│  │ Repository      │  │ Manager         │                 │
│  └─────────────────┘  └─────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────────┐
│                     Data Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐                 │
│  │ XORM Models     │  │ MySQL Database  │                 │
│  │ (db/models.go)  │  │                 │                 │
│  └─────────────────┘  └─────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
blog-fanchiikawa-service/
├── 🎯 Controller Layer
│   ├── graph/                   # GraphQL schema and generated code
│   │   ├── schema.graphqls      # GraphQL schema definition
│   │   ├── schema.resolvers.go  # GraphQL resolvers (thin layer)
│   │   └── model/               # GraphQL models
│   └── resolver/                # Business logic resolvers
│       ├── resolver.go          # Dependency injection
│       ├── mutation.go          # Mutation handlers
│       └── query.go             # Query handlers
│
├── 🔧 Service Layer
│   └── service/                 # Business logic & AWS integration
│       ├── user_service.go      # User management logic
│       ├── language_service.go  # Language detection (AWS Comprehend)
│       ├── translate_service.go # Translation (AWS Translate)  
│       ├── speech_service.go    # Text-to-speech (AWS Polly)
│       └── storage_service.go   # File storage (AWS S3)
│
├── 📦 Repository Layer  
│   └── repository/              # Data access layer
│       ├── interfaces.go        # Repository contracts
│       ├── user_repository.go   # User data operations
│       ├── user_device_repository.go # Device data operations
│       └── transaction_manager.go # Database transactions
│
├── 🗄️ Data Layer
│   ├── db/                      # Database models and connection
│   │   ├── db.go               # XORM engine setup
│   │   └── models.go           # Database models
│   └── sdk/                    # AWS SDK configuration
│       ├── aws_config.go       # Centralized AWS session
│       ├── comprehend.go       # Language detection
│       ├── translate.go        # Text translation
│       ├── polly.go           # Text-to-speech
│       └── s3.go              # File storage
│
└── 📋 Configuration
    ├── server.go               # Application entry point
    ├── .env                    # Environment variables
    └── .gitignore             # Git ignore rules
```

## 🎨 Layer Responsibilities

### 🎯 **Controller Layer** (`/graph`, `/resolver`)
- **Purpose**: Handle GraphQL requests and responses
- **Responsibilities**:
  - Parse GraphQL queries/mutations
  - Parameter validation and transformation
  - Delegate to service layer
  - Return formatted responses

### 🔧 **Service Layer** (`/service`)
- **Purpose**: Implement business logic and coordinate operations
- **Responsibilities**:
  - Core business rules and validation
  - AWS service integration
  - Cross-service coordination
  - Transaction orchestration
  - Error handling and logging

### 📦 **Repository Layer** (`/repository`)
- **Purpose**: Abstract data access operations
- **Responsibilities**:
  - Database operations with XORM
  - Data query optimization
  - Transaction management
  - Data model conversion

### 🗄️ **Data Layer** (`/db`, `/sdk`)
- **Purpose**: Handle data persistence and external services
- **Responsibilities**:
  - Database connection management
  - AWS SDK configuration
  - Data model definitions
  - Infrastructure setup

## 🛠️ Key Features

### GraphQL API
- **User Management**: Login/registration with device tracking
- **Language Services**: 
  - Language detection using AWS Comprehend
  - Text translation using AWS Translate
  - Sentiment analysis
- **Speech Services**: Text-to-speech conversion using AWS Polly
- **Storage**: File management using AWS S3
- **Image Processing**: 
  - Image label detection using AWS Rekognition
  - Text extraction from images and PDFs using AWS Textract
- **Chat Services**: 
  - Real-time chat with AWS Lex bots
  - Chat session management
  - Message history and persistence
  - WebSocket real-time communication

### Architecture Benefits
- ✅ **Single Responsibility**: Each layer has one clear purpose
- ✅ **Testability**: Business logic can be unit tested independently
- ✅ **Reusability**: Services can be used by different interfaces
- ✅ **Maintainability**: Clean separation of concerns
- ✅ **Scalability**: Easy to extend and modify

## 🚀 Getting Started

### Prerequisites
- Go 1.19+
- MySQL database
- AWS account with configured credentials

### Environment Setup

#### 1. AWS Configuration
Ensure you have AWS credentials configured. You can either:
- Configure AWS CLI: `aws configure --profile your-profile-name`
- Set up AWS credentials file in `~/.aws/credentials`
- Use IAM roles (for EC2/Lambda deployment)

#### 2. Environment Variables
Copy `.env.example` to `.env` and configure:
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```bash
# Required AWS Configuration
AWS_PROFILE=your-aws-profile
AWS_DEFAULT_REGION=ap-northeast-1

# Required AWS Lex Configuration
AWS_LEX_BOT_NAME=YourBotName
AWS_LEX_BOT_ID=your-bot-id-here
AWS_LEX_BOT_ALIAS=TSTALIASID
AWS_LEX_LOCALE_ID=en_US

# Optional Database Configuration
# DB_HOST=localhost
# DB_PORT=3306
# DB_USER=root
# DB_PASSWORD=
# DB_NAME=fanchiikawa

# Optional Application Configuration
# PORT=8080
# DEBUG=false
```

#### 3. AWS Lex Bot Setup
Before using the chat functionality, you need to create an AWS Lex bot:

1. **Create a Lex Bot in AWS Console:**
   - Go to AWS Lex Console
   - Create a new bot or use an existing one
   - Note down the Bot ID, Bot Alias, and ensure it supports your desired locale

2. **Required Permissions:**
   Your AWS credentials need the following permissions:
   ```json
   {
       "Version": "2012-10-17",
       "Statement": [
           {
               "Effect": "Allow",
               "Action": [
                   "lex:RecognizeText",
                   "lex:RecognizeUtterance"
               ],
               "Resource": "arn:aws:lex:*:*:bot/*"
           }
       ]
   }
   ```

3. **Update Environment Variables:**
   - `AWS_LEX_BOT_NAME`: Display name for your bot
   - `AWS_LEX_BOT_ID`: The unique bot ID from AWS Lex
   - `AWS_LEX_BOT_ALIAS`: Bot alias (default: TSTALIASID for test alias)
   - `AWS_LEX_LOCALE_ID`: Locale for the bot (e.g., en_US, ja_JP)

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd blog-fanchiikawa-service

# Install dependencies
go mod download

# Run the application
go run server.go
```

### Usage

The application provides multiple interfaces:
- **GraphQL Playground**: `http://localhost:8080/` - For API exploration and testing
- **Chat Interface**: `http://localhost:8080/chat/` - Modern web chat interface with auto-configured Lex settings
- **WebSocket Endpoint**: `ws://localhost:8080/ws` - For real-time communication

#### Chat Interface Features
- **Auto-Configuration**: Bot settings are automatically loaded from server environment variables
- **Simplified Setup**: Users only need to provide User ID and Chat Title
- **Real-time Messaging**: WebSocket-based instant communication
- **Chat History**: Persistent message storage and retrieval
- **Responsive Design**: Works on desktop and mobile devices

#### Example Queries

**User Management:**
```graphql
mutation {
  login(input: {
    nickname: "john"
    email: "john@example.com" 
    deviceId: "device-123"
  }) {
    id
    nickname
    email
  }
}
```

**Text-to-Speech:**
```graphql
mutation {
  textToSpeech(input: {text: "Hello world"})
}
```

**Language Detection:**
```graphql
mutation {
  detectLanguage(input: "Hello world")
}
```

**Get Lex Configuration:**
```graphql
query {
  lexConfig {
    botName
    botId
    botAlias
    localeId
  }
}
```

**Create Chat Session:**
```graphql
mutation {
  createChat(input: {
    userId: 1
    title: "My Chat Session"
    # botName, botId, botAlias, localeId are optional - will use environment variables if not provided
  }) {
    id
    title
    botName
    sessionId
  }
}
```

**Send Message to Lex Bot:**
```graphql
mutation {
  sendMessage(input: {
    chatId: 1
    message: "Hello, how can you help me?"
  }) {
    id
    content
    isUser
    intent
    sentAt
  }
}
```

**Get Chat History:**
```graphql
query {
  chatHistory(chatId: 1) {
    chat {
      id
      title
      botName
    }
    messages {
      id
      content
      isUser
      intent
      sentAt
    }
  }
}
```

**Get User's Chats:**
```graphql
query {
  userChats(userId: 1) {
    id
    title
    botName
    sessionId
    createdAt
  }
}
```

## 🧪 Testing

The layered architecture enables comprehensive testing:

- **Unit Tests**: Test service layer business logic independently
- **Integration Tests**: Test repository layer with database
- **End-to-End Tests**: Test complete GraphQL workflows

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## 🏗️ Development

### Adding New GraphQL Operations

When adding new GraphQL mutations or queries, follow these steps:

#### 1. **Update GraphQL Schema**
Edit `graph/schema.graphqls` to add your new operation:

```graphql
# Add new input types if needed
input NewFeatureInput {
  text: String!
  option: String
}

# Add to Mutation or Query type
type Mutation {
  # ... existing mutations
  newFeature(input: NewFeatureInput!): String!
}
```

#### 2. **Generate GraphQL Code**
Run gqlgen to generate the required GraphQL code:

```bash
# Generate GraphQL resolvers and models
go run github.com/99designs/gqlgen generate

# Alternative: if you have gqlgen installed globally
gqlgen generate
```

This will update:
- `graph/generated.go` - GraphQL execution code
- `graph/model/models_gen.go` - GraphQL input/output models
- `graph/schema.resolvers.go` - Resolver method stubs

#### 3. **Implement Business Logic in Service Layer**
Create or update the appropriate service:

```go
// service/new_feature_service.go
func (s *newFeatureService) ProcessNewFeature(input string) (string, error) {
    // Implement business logic here
    return result, nil
}
```

#### 4. **Create Resolver Implementation**
Update the resolver to delegate to your service:

```go
// resolver/mutation.go or resolver/query.go
func (r *Resolver) NewFeature(ctx context.Context, input model.NewFeatureInput) (string, error) {
    return r.NewFeatureService.ProcessNewFeature(input.Text)
}
```

#### 5. **Update Dependency Injection**
Add the new service to `server.go`:

```go
// Initialize new service
newFeatureService := service.NewNewFeatureService()

// Add to resolver
resolverInstance := resolver.NewResolver(
    // ... existing services
    newFeatureService,
)
```

#### 6. **Test Your Changes**
```bash
# Build and test
go run server.go

# Test in GraphQL playground at http://localhost:8080/
```

### Adding New Features (General)
1. **Define interfaces** in repository layer if data access needed
2. **Implement business logic** in service layer
3. **Create thin resolvers** in controller layer
4. **Update dependency injection** in server.go
5. **Generate GraphQL code** if schema changes are needed

### GraphQL Schema Management
- **Schema Location**: `graph/schema.graphqls`
- **Generated Models**: `graph/model/models_gen.go`
- **Resolver Stubs**: `graph/schema.resolvers.go`
- **Generated Code**: `graph/generated.go`

### Development Commands
```bash
# Generate GraphQL code
go run github.com/99designs/gqlgen generate

# Run application
go run server.go

# Run tests
go test ./...

# Build for production
go build -o blog-fanchiikawa-service
```

### Code Style
- Follow Go conventions and best practices
- Use interfaces for dependencies
- Keep resolvers thin (delegate to services)
- Centralize business logic in services
- Use meaningful error messages
- Always regenerate GraphQL code after schema changes

## 📚 Technology Stack

- **Language**: Go 1.19+
- **GraphQL**: gqlgen
- **Database**: MySQL with XORM ORM
- **Cloud Services**: AWS (Comprehend, Translate, Polly, S3)
- **Architecture**: Clean Architecture / Layered Architecture

## 📄 License

This project is licensed under the MIT License.