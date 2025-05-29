# 分享功能 API 文档

## 概述

分享功能允许用户创建分享链接，无需登录即可访问分享的内容，并提供分享统计功能。目前支持AI日报的分享。

**基础URL**: `http://localhost:8080`

## API 列表

| API | 方法 | 路径 | 说明 | 是否需要登录 |
|-----|------|------|------|-------------|
| [创建分享](#1-创建分享) | POST | `/community/share/create` | 创建分享链接 | **是** |
| [访问分享内容](#2-访问分享内容) | GET | `/community/share/{token}` | 通过token访问分享内容 | 否 |
| [获取分享统计](#3-获取分享统计) | GET | `/community/share/{token}/stats` | 获取分享的统计信息 | **是** |
| [获取支持的业务类型](#4-获取支持的业务类型) | GET | `/community/share/business-types` | 获取支持分享的业务类型列表 | 否 |

*注：需要登录的接口需要在请求头中携带有效的Authorization token

---

## 1. 创建分享

创建一个分享链接，如果已存在相同的分享则返回现有分享。**此接口需要用户登录**。

### 请求

**URL**: `POST /community/share/create`

**Content-Type**: `application/json`

**请求头**:

| 参数名 | 类型 | 是否必需 | 说明 | 示例 |
|--------|------|----------|------|------|
| Authorization | string | 是 | 用户登录token | Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9... |

**请求体参数**:

| 参数名 | 类型 | 是否必需 | 说明 | 示例 |
|--------|------|----------|------|------|
| business_type | string | 是 | 业务类型，目前只支持"ai_news" | "ai_news" |
| business_id | integer | 是 | 业务对象ID，如AI新闻的ID | 123 |
| expire_days | integer | 否 | 过期天数，0表示永久有效 | 7 |

**请求示例**:
```json
{
  "business_type": "ai_news",
  "business_id": 123,
  "expire_days": 7
}
```

### 响应

**成功响应 (200)**:
```json
{
  "code": 200,
  "message": "分享创建成功",
  "data": {
    "share_token": "abc123def456789",
    "share_url": "/community/share/abc123def456789",
    "expire_at": "2024-01-30T10:30:00Z"
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| share_token | string | 分享令牌，16位字符串 |
| share_url | string | 分享链接路径 |
| expire_at | string/null | 过期时间，ISO 8601格式，null表示永久有效 |

**错误响应**:

| 错误码 | 说明 | 响应示例 |
|--------|------|----------|
| 400 | 参数错误 | `{"code": 400, "message": "参数错误: business_type is required"}` |
| 400 | 不支持的业务类型 | `{"code": 400, "message": "创建分享失败"}` |
| 401 | 未登录或token无效 | `{"code": 401, "message": "token为空"}` |
| 403 | 用户被拉黑 | `{"code": 403, "message": "你已涉嫌违规社区文化，已被纳入小黑屋"}` |
| 500 | 服务器内部错误 | `{"code": 500, "message": "创建分享失败"}` |

---

## 2. 访问分享内容

通过分享token获取具体的分享内容，无需登录。

### 请求

**URL**: `GET /community/share/{token}`

**路径参数**:

| 参数名 | 类型 | 是否必需 | 说明 | 示例 |
|--------|------|----------|------|------|
| token | string | 是 | 分享令牌 | abc123def456789 |

**请求示例**:
```
GET /community/share/abc123def456789
```

### 响应

**成功响应 (200) - AI新闻内容**:
```json
{
  "code": 200,
  "message": "",
  "data": {
    "id": 123,
    "title": "ChatGPT发布重大更新",
    "content": "详细的新闻内容...",
    "summary": "OpenAI发布了ChatGPT的重大更新",
    "category": "AI技术",
    "tags": "ChatGPT,OpenAI,人工智能",
    "publish_date": "2024-01-15",
    "created_at": "2024-01-15T08:30:00Z",
    "comment_count": 25,
    "share_type": "ai_news"
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | integer | 内容ID |
| title | string | 标题 |
| content | string | 详细内容 |
| summary | string | 摘要 |
| category | string | 分类 |
| tags | string | 标签，多个标签用逗号分隔 |
| publish_date | string | 发布日期 |
| created_at | string | 创建时间 |
| comment_count | integer | 评论数量 |
| share_type | string | 分享类型标识 |

**错误响应**:

| 错误码 | 说明 | 响应示例 |
|--------|------|----------|
| 400 | 无效的分享链接 | `{"code": 400, "message": "无效的分享链接"}` |
| 400 | 分享链接无效或已过期 | `{"code": 400, "message": "分享链接无效或已过期"}` |
| 400 | 获取内容失败 | `{"code": 400, "message": "获取内容失败"}` |
| 400 | 不支持的分享类型 | `{"code": 400, "message": "不支持的分享类型"}` |

---

## 3. 获取分享统计

获取指定分享的统计信息，包括浏览量、创建时间等。**此接口需要用户登录**。

### 请求

**URL**: `GET /community/share/{token}/stats`

**请求头**:

| 参数名 | 类型 | 是否必需 | 说明 | 示例 |
|--------|------|----------|------|------|
| Authorization | string | 是 | 用户登录token | Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9... |

**路径参数**:

| 参数名 | 类型 | 是否必需 | 说明 | 示例 |
|--------|------|----------|------|------|
| token | string | 是 | 分享令牌 | abc123def456789 |

**请求示例**:
```
GET /community/share/abc123def456789/stats
```

### 响应

**成功响应 (200)**:
```json
{
  "code": 200,
  "message": "",
  "data": {
    "total_views": 150,
    "created_at": "2024-01-15T10:30:00Z",
    "expire_at": "2024-01-30T10:30:00Z"
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| total_views | integer | 总浏览量 |
| created_at | string | 分享创建时间，ISO 8601格式 |
| expire_at | string/null | 过期时间，ISO 8601格式，null表示永久有效 |

**错误响应**:

| 错误码 | 说明 | 响应示例 |
|--------|------|----------|
| 400 | 无效的分享token | `{"code": 400, "message": "无效的分享token"}` |
| 400 | 分享链接无效或已过期 | `{"code": 400, "message": "分享链接无效或已过期"}` |
| 401 | 未登录或token无效 | `{"code": 401, "message": "token为空"}` |
| 403 | 用户被拉黑 | `{"code": 403, "message": "你已涉嫌违规社区文化，已被纳入小黑屋"}` |

---

## 4. 获取支持的业务类型

获取系统支持分享的业务类型列表。

### 请求

**URL**: `GET /community/share/business-types`

**请求示例**:
```
GET /community/share/business-types
```

### 响应

**成功响应 (200)**:
```json
{
  "code": 200,
  "message": "",
  "data": [
    {
      "type": "ai_news",
      "description": "AI日报"
    }
  ]
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| type | string | 业务类型标识 |
| description | string | 业务类型的中文描述 |

---

## 错误码说明

| HTTP状态码 | 错误码 | 说明 |
|-----------|--------|------|
| 200 | 200 | 请求成功 |
| 400 | 400 | 请求参数错误或业务逻辑错误 |
| 500 | 500 | 服务器内部错误 |

---

## 使用示例

### 完整的分享流程示例

```bash
# 1. 获取支持的业务类型
curl -X GET "http://localhost:8080/community/share/business-types"

# 2. 创建AI新闻分享
curl -X POST "http://localhost:8080/community/share/create" \
  -H "Content-Type: application/json" \
  -d '{
    "business_type": "ai_news",
    "business_id": 123,
    "expire_days": 7
  }'

# 假设返回的token是: abc123def456789

# 3. 访问分享内容
curl -X GET "http://localhost:8080/community/share/abc123def456789"

# 4. 查看分享统计
curl -X GET "http://localhost:8080/community/share/abc123def456789/stats"
```

### JavaScript 调用示例

```javascript
// 创建分享
async function createShare(businessId, expireDays = 7) {
  const response = await fetch('/community/share/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      business_type: 'ai_news',
      business_id: businessId,
      expire_days: expireDays
    })
  });
  
  const result = await response.json();
  if (result.code === 200) {
    return result.data;
  } else {
    throw new Error(result.message);
  }
}

// 获取分享内容
async function getSharedContent(token) {
  const response = await fetch(`/community/share/${token}`);
  const result = await response.json();
  
  if (result.code === 200) {
    return result.data;
  } else {
    throw new Error(result.message);
  }
}

// 获取分享统计
async function getShareStats(token) {
  const response = await fetch(`/community/share/${token}/stats`);
  const result = await response.json();
  
  if (result.code === 200) {
    return result.data;
  } else {
    throw new Error(result.message);
  }
}

// 使用示例
try {
  // 创建分享
  const shareData = await createShare(123, 7);
  console.log('分享创建成功:', shareData.share_url);
  
  // 获取分享内容
  const content = await getSharedContent(shareData.share_token);
  console.log('分享内容:', content.title);
  
  // 获取统计
  const stats = await getShareStats(shareData.share_token);
  console.log('浏览量:', stats.total_views);
} catch (error) {
  console.error('操作失败:', error.message);
}
```

---

## 注意事项

1. **防刷机制**: 同一IP在5分钟内多次访问同一分享链接，只计算一次浏览量
2. **过期检查**: 访问过期的分享链接会返回错误
3. **业务类型限制**: 目前只支持"ai_news"类型的分享
4. **重复创建**: 同一用户对同一内容重复创建分享会返回已存在的分享链接
5. **异步统计**: 访问记录是异步处理的，不会影响内容获取的响应速度
6. **无需登录**: 所有分享相关的接口都不需要登录验证 