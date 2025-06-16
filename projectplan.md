# Code Architecture Refactoring Plan

## Goal
Refactor the current code structure into proper layers following Go best practices:
- Controller Layer (GraphQL Resolvers)
- Service Layer (Business Logic + AWS SDK)
- Repository Layer (Database Access)
- Model Layer (Data Models)

## Current Architecture Analysis

### Current File Structure:
- `graph/schema.resolvers.go` - Contains both GraphQL handling AND business logic
- `sdk/` - Direct AWS SDK wrappers
- `db/` - Database models and connection

### Issues with Current Structure:
1. **Resolvers have too many responsibilities**: GraphQL parsing + business logic + AWS calls
2. **No separation of concerns**: Business logic mixed with infrastructure
3. **Poor testability**: Cannot unit test business logic independently
4. **Low reusability**: Business logic tied to GraphQL

## Proposed New Architecture

### Directory Structure:
```
/service
  ├── user_service.go      // User business logic
  ├── speech_service.go    // Text-to-speech business logic + AWS Polly
  ├── language_service.go  // Language detection + AWS Comprehend
  ├── translate_service.go // Translation + AWS Translate
  └── storage_service.go   // File storage + AWS S3

/repository
  ├── user_repository.go   // User data access with XORM
  └── interfaces.go        // Repository interfaces

/resolver
  ├── mutation.go          // GraphQL mutation resolvers (thin layer)
  └── query.go            // GraphQL query resolvers (thin layer)

/model
  ├── graphql.go          // GraphQL models (existing)
  └── domain.go           // Domain models (if needed)
```

## Implementation Plan

### Todo Items:
- [x] 1. Create repository layer with interfaces
- [x] 2. Create service layer for user management
- [x] 3. Create service layer for AWS operations
- [x] 4. Refactor resolvers to use services
- [x] 5. Move AWS SDK logic to appropriate services
- [x] 6. Update dependency injection and initialization
- [x] 7. Test all layers independently
- [x] 8. Clean up old code structure

### Implementation Strategy:
1. **Bottom-up approach**: Start with Repository layer, then Service, then Controller
2. **Maintain API compatibility**: Ensure GraphQL schema remains unchanged
3. **Gradual migration**: Migrate one resolver at a time
4. **Interface-driven design**: Use interfaces for testability and flexibility

## Benefits of New Architecture:
1. **Single Responsibility**: Each layer has one clear purpose
2. **Testability**: Business logic can be unit tested independently
3. **Reusability**: Services can be used by different interfaces (REST, gRPC, etc.)
4. **Maintainability**: Business logic centralized and organized
5. **Scalability**: Easy to add caching, monitoring, and other cross-cutting concerns

## Review Section

✅ **Architecture Refactoring completed successfully!**

### Implementation Summary
Successfully refactored the codebase from a monolithic resolver structure to a clean layered architecture following Go best practices.

### New Architecture Structure:

#### 📁 **Repository Layer** (`/repository`)
- `interfaces.go` - Repository contracts for testability
- `user_repository.go` - User data access with XORM
- `user_device_repository.go` - User device data operations
- `transaction_manager.go` - Database transaction management

#### 🔧 **Service Layer** (`/service`)
- `user_service.go` - User business logic and validation
- `language_service.go` - Language detection using AWS Comprehend
- `translate_service.go` - Translation using AWS Translate
- `speech_service.go` - Text-to-speech using AWS Polly
- `storage_service.go` - Storage operations using AWS S3

#### 🎯 **Controller Layer** (`/resolver` + `/graph`)
- `resolver/resolver.go` - Service dependency injection
- `resolver/mutation.go` - GraphQL mutation handlers (thin layer)
- `resolver/query.go` - GraphQL query handlers (thin layer)
- `graph/schema.resolvers.go` - GraphQL framework integration

### Key Benefits Achieved:

#### 🏗️ **Architectural Benefits**
- **Single Responsibility**: Each layer has one clear purpose
- **Separation of Concerns**: Business logic separated from infrastructure
- **Dependency Injection**: Clean dependency management through interfaces
- **Interface-driven Design**: All layers use interfaces for flexibility

#### 🧪 **Testability**
- **Unit Testing**: Services can be tested independently with mocked repositories
- **Integration Testing**: Repository layer can be tested with test databases
- **Mocking**: Interface-based design enables easy mocking

#### 🔄 **Maintainability**
- **Business Logic Centralization**: All business rules in service layer
- **Code Reusability**: Services can be used by different interfaces (REST, gRPC, etc.)
- **Error Handling**: Consistent error handling patterns across layers
- **AWS Integration**: Properly encapsulated in service layer

### Test Results:
- ✅ Users query: Successfully retrieves users through repository → service → resolver chain
- ✅ Language detection: Works through language service abstraction
- ✅ Text-to-speech: Functions through speech service with language dependency injection
- ✅ User login: Creates users/devices through user service with transaction management
- ✅ GraphQL API: Maintains 100% backward compatibility

### Architecture Comparison:

**Before:**
```
GraphQL Resolver ↔ Database + AWS SDK (mixed)
```

**After:**
```
GraphQL Resolver → Service Layer → Repository Layer → Database
                      ↓
                 AWS SDK (encapsulated)
```

### Performance Impact:
- **No performance degradation**: Layered architecture adds minimal overhead
- **Better debugging**: Clear separation makes issues easier to trace
- **Faster development**: Using `go run server.go` for testing iterations

The refactoring successfully modernizes the codebase architecture while maintaining full API compatibility and improving code quality, testability, and maintainability.