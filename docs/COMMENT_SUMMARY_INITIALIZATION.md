# AI评论总结数据初始化指南

## 概述
这是一个用于为现有文章和章节生成AI评论总结的初始化工具。由于AI评论总结是新功能，需要为历史数据生成总结。

## 测试用例说明

### 🔍 干跑模式测试
```bash
go test ./server/service -v -run TestDryRunInitialization
```
- **用途**: 统计需要初始化的数据量，不实际生成总结
- **建议**: 在正式初始化前先运行此测试，了解数据规模

### 📚 文章总结初始化
```bash
go test ./server/service -v -run TestInitializeCommentSummariesForArticles
```
- **用途**: 为所有符合条件的文章生成AI评论总结
- **条件**: 评论数量 ≥ 3 条且尚未有总结的文章

### 📖 章节总结初始化
```bash
go test ./server/service -v -run TestInitializeCommentSummariesForSections
```
- **用途**: 为所有符合条件的章节生成AI评论总结
- **条件**: 评论数量 ≥ 3 条且尚未有总结的章节

### 🎯 全量初始化
```bash
go test ./server/service -v -run TestInitializeAllCommentSummaries
```
- **用途**: 一次性初始化文章和章节的所有评论总结
- **建议**: 数据量大时慎用，建议分批进行

## 配置参数

在 `comment_summary_init_test.go` 文件顶部可以调整以下参数：

```go
const (
    BATCH_SIZE = 10 // 每批处理的数量
    MIN_COMMENT_COUNT = 3  // 最少评论数量要求
    MAX_COMMENT_COUNT = 50 // 评论数量警告阈值
    MAX_CONCURRENT = 3 // 最大并发数
    DELAY_BETWEEN_REQUESTS = 1 * time.Second // 请求间延迟
)

// 可选：指定特定ID列表进行初始化
var (
    FILTER_ARTICLE_IDS = []int{} // 指定文章ID，空表示处理所有
    FILTER_SECTION_IDS = []int{} // 指定章节ID，空表示处理所有
)
```

## 使用建议

### 🚀 首次使用流程
1. **环境检查**: 确保LLM服务配置正确
2. **数据统计**: 先运行干跑模式了解数据规模
3. **小批量测试**: 设置特定ID列表进行小范围测试
4. **全量处理**: 确认无误后进行全量初始化

### ⚠️ 注意事项
- **API费用**: 每次AI调用都会产生费用，请合理评估成本
- **处理时间**: 大量数据可能需要数小时，建议在业务低峰期进行
- **网络稳定性**: 确保网络连接稳定，避免中途中断
- **重复运行**: 已有总结的数据会自动跳过，支持重复运行

### 🛠️ 自定义配置示例

#### 只初始化特定文章
```go
var (
    FILTER_ARTICLE_IDS = []int{1, 2, 3, 10, 25} // 只处理这些文章
    FILTER_SECTION_IDS = []int{} // 不处理章节
)
```

#### 高频率处理（小心使用）
```go
const (
    BATCH_SIZE = 20 // 增大批次
    DELAY_BETWEEN_REQUESTS = 500 * time.Millisecond // 减少延迟
)
```

#### 保守处理（推荐）
```go
const (
    BATCH_SIZE = 5 // 小批次
    DELAY_BETWEEN_REQUESTS = 2 * time.Second // 增加延迟
)
```

## 输出示例

```
📊 文章初始化统计:
   - 需要处理的文章数量: 45
   - 最小评论数量要求: 3
   - 批量处理大小: 10
   - 请求间延迟: 1s

🚀 处理批次 1-10 (共 45 篇文章)
🔄 处理文章: 如何学习Go语言 (ID: 12)
✅ 文章总结生成成功 (耗时: 3.2s, 评论数: 8, 总结长度: 456字符)
...

🏁 文章评论总结初始化完成!
═══════════════════════════════════════
📊 最终统计:
   - 总处理数量: 45
   - 成功数量: 43
   - 失败数量: 2
   - 总耗时: 15m30s
   - 平均耗时: 20.7s
```

## 故障排除

### 常见问题
1. **LLM服务连接失败**: 检查网络和API配置
2. **数据库连接问题**: 检查数据库配置和权限
3. **内存不足**: 减少 `BATCH_SIZE` 参数
4. **请求频率过高**: 增加 `DELAY_BETWEEN_REQUESTS`

### 恢复策略
- 重复运行测试会自动跳过已完成的总结
- 可以通过设置特定ID列表来处理失败的条目
- 检查错误日志定位具体问题

## 监控建议
- 观察API调用费用
- 监控服务器资源使用情况
- 关注错误率和处理时间
- 定期检查生成的总结质量