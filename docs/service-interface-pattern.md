# 服务接口设计模式

## 问题背景

之前的代码存在不一致的设计模式：

### 不一致的服务定义
```go
// ❌ 混合模式 - 不一致
type Resolver struct {
    UserService    service.UserService     // 接口类型
    ChatService    *service.ChatService    // 指针类型 ❌
    ConfigService  *service.ConfigService  // 指针类型 ❌
}
```

## 解决方案：统一接口模式

### ✅ 一致的接口设计
```go
// ✅ 统一接口模式
type Resolver struct {
    UserService    service.UserService    // 接口
    ChatService    service.ChatService    // 接口
    ConfigService  service.ConfigService  // 接口
}
```

## 接口设计模式的优势

### 1. **依赖注入 (Dependency Injection)**
```go
// 接口使得依赖注入更清晰
func NewResolver(
    userService service.UserService,      // 可以注入任何实现
    chatService service.ChatService,      // 可以注入任何实现
    configService service.ConfigService,  // 可以注入任何实现
) *Resolver
```

### 2. **可测试性 (Testability)**
```go
// 容易创建 mock 对象进行单元测试
type mockChatService struct{}
func (m *mockChatService) CreateChat(...) (*ChatResponse, error) {
    return &ChatResponse{ID: 123}, nil
}

// 测试时注入 mock
resolver := NewResolver(
    userService,
    &mockChatService{}, // 注入测试用的 mock
    configService,
)
```

### 3. **松耦合 (Loose Coupling)**
```go
// Resolver 只依赖接口，不依赖具体实现
type ChatService interface {
    CreateChat(req *CreateChatRequest) (*ChatResponse, error)
    SendMessage(ctx context.Context, req *SendMessageRequest) (*MessageResponse, error)
    // ... 其他方法
}

// 具体实现可以随时替换
type chatService struct { /* 实现细节 */ }
type advancedChatService struct { /* 不同的实现 */ }
```

### 4. **接口隔离 (Interface Segregation)**
```go
// 每个服务都有清晰的职责界限
type ConfigService interface {
    GetLexConfig() *LexConfig  // 只关心配置相关的方法
}

type ChatService interface {
    CreateChat(...) (...)      // 只关心聊天相关的方法
    SendMessage(...) (...)
}
```

## 实现模式

### 接口定义
```go
// 1. 定义公开接口
type ChatService interface {
    CreateChat(req *CreateChatRequest) (*ChatResponse, error)
    // ... 其他公开方法
}

// 2. 私有实现结构体
type chatService struct {
    chatRepo        repository.ChatRepository
    chatMessageRepo repository.ChatMessageRepository
    lexService      *sdk.LexService
}

// 3. 构造函数返回接口
func NewChatService(...) ChatService {
    return &chatService{ /* 初始化 */ }
}
```

### 方法实现
```go
// 4. 实现接口方法 (私有结构体)
func (s *chatService) CreateChat(req *CreateChatRequest) (*ChatResponse, error) {
    // 具体实现
}
```

## 代码质量提升

### 修改前 (不一致)
```go
❌ 混合类型，难以测试
chatService    *service.ChatService    // 具体类型
configService  *service.ConfigService  // 具体类型
```

### 修改后 (一致)
```go
✅ 统一接口，易于测试和维护
ChatService    service.ChatService     // 接口
ConfigService  service.ConfigService   // 接口
```

## 总结

通过统一使用接口模式：

1. **提高代码一致性** - 所有服务都遵循相同的设计模式
2. **增强可测试性** - 容易创建 mock 对象进行单元测试
3. **降低耦合度** - 依赖接口而非具体实现
4. **提升可维护性** - 更容易替换和扩展服务实现

这是Go语言中推荐的企业级应用设计模式！