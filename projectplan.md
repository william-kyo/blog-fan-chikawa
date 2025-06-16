# Database Analysis and XORM Refactoring Plan

## Current Database Implementation Analysis

### Current Database Setup
- **Database Driver**: MySQL using `github.com/go-sql-driver/mysql v1.9.2`
- **Connection Method**: Standard `database/sql` package with raw SQL queries
- **Database Name**: `fanchiikawa`
- **Connection Details**: 
  - Host: localhost:3306
  - User: root
  - Password: (empty)
  - Parse Time: enabled

### Database Tables (Inferred from Models)
Based on the GraphQL models and SQL queries, the database has the following tables:

1. **user table**
   - `id` (INT64, PRIMARY KEY, AUTO_INCREMENT)
   - `nickname` (VARCHAR)
   - `email` (VARCHAR)
   - `created_at` (TIMESTAMP)
   - `updated_at` (TIMESTAMP)

2. **user_device table**
   - `id` (INT64, PRIMARY KEY, AUTO_INCREMENT)
   - `user_id` (INT64, FOREIGN KEY to user.id)
   - `device_id` (VARCHAR)
   - `created_at` (TIMESTAMP)
   - `updated_at` (TIMESTAMP)

### Current SQL Queries Found
1. **SELECT queries**:
   - `SELECT * FROM user WHERE email = ?` (in Login resolver)
   - `SELECT * FROM user WHERE id = ?` (in Login resolver)
   - `SELECT * FROM user limit 10` (in Users resolver)

2. **INSERT queries**:
   - `INSERT INTO user (nickname, email) VALUES (?, ?)` (in Login resolver)
   - `INSERT INTO user_device (user_id, device_id) VALUES (?, ?)` (in Login resolver)

3. **Transaction Usage**:
   - Transaction used in Login resolver for creating user and device atomically

### Current Database Integration Points
1. **`/Users/kyo/IdeaProject/blog-fanchiikawa-service/db/db.go`**:
   - Database connection initialization
   - Global `MySQL` variable for database access

2. **`/Users/kyo/IdeaProject/blog-fanchiikawa-service/graph/schema.resolvers.go`**:
   - Login mutation (lines 19-81): User creation and device registration
   - Users query (lines 127-143): Fetch users with limit

3. **`/Users/kyo/IdeaProject/blog-fanchiikawa-service/server.go`**:
   - Database initialization on startup (line 22)

### Issues with Current Implementation
1. **Raw SQL queries**: Prone to SQL injection if not handled carefully
2. **No ORM benefits**: No automatic struct mapping, relationship handling
3. **Manual transaction management**: Verbose transaction handling
4. **No database migrations**: No schema versioning or migration system
5. **Hard-coded connection**: Database credentials are hard-coded
6. **No connection pooling configuration**: Using default connection pool settings

## XORM Refactoring Plan

### Todo Items

- [x] 1. Add XORM dependency to go.mod
- [x] 2. Create XORM model structs with proper tags
- [x] 3. Refactor database connection to use XORM engine
- [x] 4. Create database migration/sync functionality
- [x] 5. Refactor Login resolver to use XORM
- [x] 6. Refactor Users query resolver to use XORM
- [x] 7. Add proper error handling and logging
- [x] 8. Test all GraphQL operations
- [x] 9. Update database initialization in server.go
- [x] 10. Clean up unused raw SQL code

### Implementation Strategy
1. **Incremental approach**: Replace one resolver at a time to minimize risk
2. **Maintain backward compatibility**: Ensure GraphQL schema remains unchanged
3. **Add proper struct tags**: Use XORM tags for database mapping
4. **Implement connection pooling**: Configure XORM engine with proper settings
5. **Add database migrations**: Use XORM's sync functionality for schema management

### Benefits of XORM Migration
1. **Type safety**: Compile-time checking of database operations
2. **Automatic mapping**: Struct to table mapping with tags
3. **Built-in transactions**: Simplified transaction handling
4. **Schema synchronization**: Automatic table creation/updates
5. **Connection management**: Built-in connection pooling
6. **Query builder**: More readable and maintainable queries

## Review Section

✅ **XORM Refactoring completed successfully!**

### Implementation Summary
Successfully migrated from raw SQL (`database/sql`) to XORM ORM with the following accomplishments:

### Files Created/Modified:
1. **db/models.go** - New XORM model structs with proper tags
2. **db/db.go** - Refactored database connection to use XORM engine  
3. **graph/schema.resolvers.go** - Updated Login and Users resolvers to use XORM
4. **go.mod** - Added XORM dependency

### Key Features Implemented:
- **XORM Models**: Proper struct tags for database mapping
- **Connection Pooling**: Configured max idle (10) and open (100) connections
- **Schema Synchronization**: Automatic table sync with `Engine.Sync2()`
- **Transaction Support**: Simplified transaction handling in Login resolver
- **Type Safety**: Compile-time checking of database operations
- **Debug Logging**: Configurable SQL logging via DEBUG environment variable

### Test Results:
- ✅ Users query: Successfully retrieves user list using `Engine.Limit(10).Find()`
- ✅ Login mutation: Creates new users and devices in transaction
- ✅ Schema sync: Automatically syncs database structure on startup
- ✅ Other services: AWS services (textToSpeech, etc.) remain unaffected

### Benefits Achieved:
- **Eliminated raw SQL**: Replaced manual query building with XORM methods
- **Better error handling**: XORM provides more descriptive error messages
- **Simplified transactions**: Session-based transaction management
- **Automatic mapping**: Direct struct-to-table mapping without manual scanning
- **Connection management**: Built-in connection pooling and configuration

### GraphQL API Compatibility:
All existing GraphQL operations continue to work exactly as before - the refactoring is completely transparent to API consumers.

The migration successfully modernizes the data layer while maintaining full backward compatibility.