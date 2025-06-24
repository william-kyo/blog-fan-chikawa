# AWS SDK版本管理说明

## 为什么需要两个AWS SDK版本？

### AWS SDK v1 (`github.com/aws/aws-sdk-go`)
- **用途**: 现有服务 (S3, Comprehend, Translate, Polly, Rekognition, Textract)
- **稳定性**: 成熟稳定，广泛使用
- **配置类型**: `*session.Session`

### AWS SDK v2 (`github.com/aws/aws-sdk-go-v2`)
- **用途**: 新服务 (Lex Runtime V2)
- **原因**: Lex Runtime V2 API只在SDK v2中提供
- **配置类型**: `aws.Config`
- **优势**: 更好的性能，更现代的API设计

## 统一配置管理

### 优化前的问题
```go
// 重复的配置代码
func InitAWSSession() {
    region := os.Getenv("AWS_DEFAULT_REGION")
    profile := os.Getenv("AWS_PROFILE")
    // ... SDK v1 初始化
}

func InitAWSConfigV2() {
    region := os.Getenv("AWS_DEFAULT_REGION") // 重复
    profile := os.Getenv("AWS_PROFILE")       // 重复
    // ... SDK v2 初始化
}
```

### 优化后的解决方案
```go
// 统一的配置获取
func getAWSCredentials() (region, profile string) {
    // 统一从环境变量读取配置
}

// 统一初始化方法
func InitAWS() {
    region, profile := getAWSCredentials()
    // 同时初始化SDK v1和v2
}

// 保持向后兼容的方法
func InitAWSSession() { /* legacy support */ }
func InitAWSConfigV2() { /* legacy support */ }
```

## 使用方式

### 简化的服务器初始化
```go
func main() {
    // 一次调用初始化所有AWS配置
    sdk.InitAWS()
    
    // 其他服务初始化...
}
```

### 服务中使用不同版本
```go
// 使用SDK v1的服务
func NewS3Service() {
    session := sdk.GetAWSSession()
    return s3.New(session)
}

// 使用SDK v2的服务
func NewLexService() {
    config := sdk.GetAWSConfig()
    return lexruntimev2.NewFromConfig(config)
}
```

## 未来计划

1. **渐进式迁移**: 可以逐步将现有服务迁移到SDK v2
2. **统一SDK**: 当所有服务都迁移后，可以移除SDK v1依赖
3. **向后兼容**: 保持旧的初始化方法，确保不破坏现有代码

## 总结

这种设计既解决了代码重复问题，又保持了向后兼容性，同时支持新的Lex Runtime V2功能。通过统一的`InitAWS()`方法，简化了AWS配置管理。