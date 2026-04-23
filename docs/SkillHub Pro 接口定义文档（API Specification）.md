# SkillHub Pro 接口定义文档（API Specification）

## 文档信息

| 项目     | 内容                     |
| -------- | ------------------------ |
| 项目名称 | SkillHub Pro（技能宝库） |
| 文档类型 | API 接口定义说明书       |
| 版本     | V1.0                     |
| 创建日期 | 2026-04-23               |
| 基础路径 | `/api/v1`                |
| 协议     | HTTPS                    |
| 数据格式 | JSON                     |
| 字符编码 | UTF-8                    |


## 1. 概述

本文档详细定义 SkillHub Pro 平台的所有 RESTful API 接口，包括请求方法、路径、参数、请求体结构、响应格式及错误码说明。接口设计遵循 REST 风格，使用标准 HTTP 状态码，并统一响应格式。

### 1.1 基础 URL

- **开发环境**：`https://dev-api.skillhub.pro/v1`
- **生产环境**：`https://api.skillhub.pro/v1`

### 1.2 统一响应格式

所有接口统一返回 JSON 格式数据，结构如下：

**成功响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "request_id": "uuid"
}
```

**错误响应**：
```json
{
  "code": 1001,
  "message": "Invalid parameter: email",
  "data": null,
  "request_id": "uuid"
}
```

### 1.3 认证方式

| 接口类型 | 认证方式                        | 说明                     |
| -------- | ------------------------------- | ------------------------ |
| 公开接口 | 无需认证                        | 如技能列表、详情、搜索   |
| 用户接口 | Bearer Token (JWT)              | 登录后获取，有效期 7 天  |
| API 接口 | `X-API-Key` 请求头              | 在个人中心生成的 API Key |
| 管理接口 | Bearer Token (JWT) + 管理员权限 | 需要 `is_admin=true`     |

### 1.4 错误码定义

| 错误码 | HTTP状态码 | 说明                            |
| ------ | ---------- | ------------------------------- |
| 0      | 200        | 请求成功                        |
| 1001   | 400        | 参数错误                        |
| 1002   | 401        | 认证失败（未登录或 Token 无效） |
| 1003   | 403        | 权限不足                        |
| 1004   | 404        | 资源不存在                      |
| 1005   | 409        | 资源冲突                        |
| 2001   | 500        | 数据库错误                      |
| 2002   | 500        | 外部服务错误（GitHub API 等）   |
| 2003   | 503        | 服务暂时不可用                  |
| 3001   | 400        | 路由匹配失败（无相关技能）      |
| 3002   | 500        | LLM 调用失败                    |

### 1.5 分页参数

列表接口统一支持分页，请求参数如下：

| 参数 | 类型 | 必填 | 说明                        | 示例 |
| ---- | ---- | ---- | --------------------------- | ---- |
| page | int  | 否   | 页码，默认 1                | 2    |
| size | int  | 否   | 每页数量，默认 20，最大 100 | 20   |

分页响应结构：
```json
{
  "code": 0,
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "size": 20,
      "total": 10234,
      "total_pages": 512
    }
  }
}
```


## 2. 公开接口

### 2.1 技能相关

#### 2.1.1 获取技能列表

**GET** `/skills`

**请求参数**

| 参数            | 类型   | 必填 | 说明                                                    | 示例            |
| --------------- | ------ | ---- | ------------------------------------------------------- | --------------- |
| page            | int    | 否   | 页码                                                    | 1               |
| size            | int    | 否   | 每页数量                                                | 20              |
| category        | string | 否   | 分类 slug                                               | `data-analysis` |
| source_type     | string | 否   | 来源：`official` / `github` / `community`               | `github`        |
| min_stars       | int    | 否   | 最小 GitHub 星数                                        | 100             |
| sort            | string | 否   | 排序字段：`stars` / `downloads` / `updated` / `created` | `stars`         |
| order           | string | 否   | 排序方向：`asc` / `desc`，默认 `desc`                   | `desc`          |
| security_status | string | 否   | 安全状态：`safe` / `warning`                            | `safe`          |

**响应示例**
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "excel-trend-analyzer",
        "description": "分析 Excel 数据趋势并生成可视化图表",
        "source_type": "official",
        "source_url": "https://github.com/anthropics/skills",
        "github_stars": 12500,
        "github_forks": 1200,
        "category": {
          "id": "category-uuid",
          "name": "数据分析",
          "slug": "data-analysis"
        },
        "author": "Anthropic",
        "version": "1.2.0",
        "download_count": 8430,
        "rating_avg": 4.7,
        "rating_count": 128,
        "security_status": "safe",
        "created_at": "2026-03-15T10:00:00Z",
        "updated_at": "2026-04-20T08:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 20,
      "total": 10234,
      "total_pages": 512
    }
  },
  "request_id": "abc-123"
}
```

---

#### 2.1.2 获取技能详情

**GET** `/skills/{id}`

**路径参数**

| 参数 | 类型 | 必填 | 说明         |
| ---- | ---- | ---- | ------------ |
| id   | UUID | 是   | 技能唯一标识 |

**响应示例**
```json
{
  "code": 0,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "excel-trend-analyzer",
    "description": "分析 Excel 数据趋势并生成可视化图表",
    "source_type": "official",
    "source_url": "https://github.com/anthropics/skills",
    "github_stars": 12500,
    "github_forks": 1200,
    "category": {
      "id": "category-uuid",
      "name": "数据分析",
      "slug": "data-analysis"
    },
    "author": "Anthropic",
    "version": "1.2.0",
    "download_count": 8430,
    "view_count": 15620,
    "rating_avg": 4.7,
    "rating_count": 128,
    "security_status": "safe",
    "skill_md_content": "---\nname: Excel Trend Analyzer\ndescription: ...\n---\n\n# Instructions\n...",
    "skill_md_url": "https://oss.skillhub.pro/skills/uuid/SKILL.md",
    "metadata": {
      "tags": ["excel", "data-analysis", "visualization"],
      "license": "MIT"
    },
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-04-20T08:30:00Z"
  },
  "request_id": "abc-123"
}
```

---

#### 2.1.3 搜索技能（全文/语义）

**POST** `/skills/search`

**请求体**
```json
{
  "query": "分析Excel销售趋势",
  "search_type": "semantic",  // 可选：fulltext / semantic，默认 fulltext
  "filters": {
    "category": "data-analysis",
    "source_type": ["official", "github"],
    "min_stars": 50,
    "security_status": "safe"
  },
  "page": 1,
  "size": 20
}
```

**响应格式**：与列表接口相同，`items` 中包含匹配结果。

---

#### 2.1.4 获取热门技能

**GET** `/skills/trending`

**请求参数**

| 参数  | 类型 | 必填 | 说明                       |
| ----- | ---- | ---- | -------------------------- |
| limit | int  | 否   | 返回数量，默认 20，最大 50 |

**响应示例**
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "id": "uuid",
        "name": "excel-trend-analyzer",
        "description": "...",
        "github_stars": 12500,
        "download_count": 8430,
        "rating_avg": 4.7,
        "hot_score": 9850.5
      }
    ]
  }
}
```

---

#### 2.1.5 获取分类列表

**GET** `/categories`

**响应示例**
```json
{
  "code": 0,
  "data": [
    {
      "id": "uuid",
      "name": "数据分析",
      "slug": "data-analysis",
      "description": "数据处理、分析与可视化相关技能",
      "icon": "chart-bar",
      "parent_id": null,
      "sort_order": 10,
      "skill_count": 1234
    },
    {
      "id": "uuid-2",
      "name": "文档处理",
      "slug": "document-processing",
      "parent_id": null,
      "skill_count": 567
    }
  ]
}
```

### 2.2 统计相关

#### 2.2.1 获取平台统计

**GET** `/stats`

**响应示例**
```json
{
  "code": 0,
  "data": {
    "total_skills": 10234,
    "total_users": 5432,
    "total_installs": 124500,
    "total_api_calls": 2340000,
    "categories_count": 24,
    "updated_at": "2026-04-23T10:00:00Z"
  }
}
```


## 3. 用户认证接口

### 3.1 注册

**POST** `/auth/register`

**请求体**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123",
  "name": "张三"
}
```

**响应**（成功返回用户信息及 JWT Token）
```json
{
  "code": 0,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "张三",
      "avatar_url": "",
      "is_admin": false
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 604800
  }
}
```

### 3.2 邮箱登录

**POST** `/auth/login`

**请求体**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123"
}
```

**响应**：同上，返回用户信息及 JWT Token。

### 3.3 GitHub OAuth 登录

**GET** `/auth/github`
重定向至 GitHub 授权页面。

**GET** `/auth/github/callback`
授权后回调，返回 JWT Token（可重定向至前端页面并携带 token）。

### 3.4 获取当前用户信息

**GET** `/auth/me`
需要认证：Bearer Token

**响应示例**
```json
{
  "code": 0,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "张三",
    "avatar_url": "https://avatars.githubusercontent.com/u/xxx",
    "is_admin": false,
    "created_at": "2026-01-15T10:00:00Z"
  }
}
```

### 3.5 刷新 Token

**POST** `/auth/refresh`
需要认证：Bearer Token（可使用即将过期的 Token）

**响应**：返回新的 `access_token`。


## 4. 用户功能接口（需认证）

### 4.1 收藏管理

#### 4.1.1 获取收藏列表

**GET** `/user/favorites`
**响应**：分页返回技能列表。

#### 4.1.2 添加收藏

**POST** `/user/favorites`
**请求体**
```json
{
  "skill_id": "uuid"
}
```

#### 4.1.3 取消收藏

**DELETE** `/user/favorites/{skill_id}`

### 4.2 评分与评论

#### 4.2.1 提交评分/评论

**POST** `/skills/{skill_id}/ratings`
**请求体**
```json
{
  "rating": 5,
  "comment": "非常实用的技能，准确率高。"
}
```

#### 4.2.2 获取技能评论列表

**GET** `/skills/{skill_id}/ratings`
支持分页。

### 4.3 API Key 管理

#### 4.3.1 获取 API Key 列表

**GET** `/user/api-keys`
**响应示例**
```json
{
  "code": 0,
  "data": [
    {
      "id": "uuid",
      "key_prefix": "sk_live_****abcd",
      "name": "My CLI Tool",
      "last_used_at": "2026-04-20T08:30:00Z",
      "expires_at": null,
      "is_active": true,
      "created_at": "2026-03-01T10:00:00Z"
    }
  ]
}
```

#### 4.3.2 创建 API Key

**POST** `/user/api-keys`
**请求体**
```json
{
  "name": "My CLI Tool",
  "expires_at": "2026-12-31T23:59:59Z"  // 可选，不填表示永不过期
}
```
**响应**：返回完整 Key（仅此时可见）
```json
{
  "code": 0,
  "data": {
    "id": "uuid",
    "key": "sk_live_xxxx_example_key_placeholder_xxxxx",
    "key_prefix": "sk_live_****abcd",
    "name": "My CLI Tool"
  }
}
```

#### 4.3.3 吊销 API Key

**DELETE** `/user/api-keys/{id}`


## 5. 智能路由接口

### 5.1 匹配技能

**POST** `/router/match`

**请求体**
```json
{
  "query": "帮我分析这个Excel表格里的销售趋势，并生成一份PPT报告",
  "top_k": 3,
  "include_content": false,
  "filters": {
    "source_type": ["official", "github"],
    "min_stars": 50
  }
}
```

| 字段            | 类型   | 必填 | 说明                              |
| --------------- | ------ | ---- | --------------------------------- |
| query           | string | 是   | 自然语言查询                      |
| top_k           | int    | 否   | 返回最佳匹配数量，默认 3，最大 10 |
| include_content | bool   | 否   | 是否包含 SKILL.md URL，默认 false |
| filters         | object | 否   | 过滤条件（同技能列表过滤）        |

**响应示例**
```json
{
  "code": 0,
  "data": {
    "matches": [
      {
        "skill_id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "excel-trend-analyzer",
        "description": "分析Excel数据趋势并生成可视化图表",
        "match_score": 0.952,
        "match_reason": "该技能擅长数据处理与可视化",
        "skill_md_url": "https://oss.skillhub.pro/skills/.../SKILL.md"
      },
      {
        "skill_id": "uuid-2",
        "name": "ppt-report-generator",
        "description": "从数据生成专业PPT报告",
        "match_score": 0.874,
        "match_reason": "该技能可生成演示文稿"
      }
    ],
    "meta": {
      "embedding_time_ms": 120,
      "vector_search_time_ms": 45,
      "rerank_time_ms": 38,
      "total_time_ms": 203
    }
  },
  "request_id": "router-abc-123"
}
```

### 5.2 执行技能（端到端）

**POST** `/router/execute`

**请求体**
```json
{
  "query": "分析附件中的销售数据并生成周报",
  "top_k": 1,
  "execution_config": {
    "model": "claude-3.5-sonnet",
    "max_tokens": 4096,
    "temperature": 0.2,
    "stream": false
  },
  "context": {
    "user_id": "optional-identifier",
    "session_id": "session-123"
  },
  "files": [
    {
      "name": "sales_march.xlsx",
      "url": "https://user-upload.example.com/sales.xlsx"
    }
  ]
}
```

| 字段                         | 类型   | 必填 | 说明                                      |
| ---------------------------- | ------ | ---- | ----------------------------------------- |
| query                        | string | 是   | 自然语言任务描述                          |
| top_k                        | int    | 否   | 使用前多少名技能组装上下文，默认 1        |
| execution_config.model       | string | 否   | 指定 LLM 模型（默认 `claude-3.5-sonnet`） |
| execution_config.max_tokens  | int    | 否   | 最大输出 token 数，默认 4096              |
| execution_config.temperature | float  | 否   | 采样温度，默认 0.2                        |
| execution_config.stream      | bool   | 否   | 是否流式输出（SSE），默认 false           |
| context                      | object | 否   | 自定义上下文信息                          |
| files                        | array  | 否   | 附件文件信息（需可公网访问的 URL）        |

**响应示例**
```json
{
  "code": 0,
  "data": {
    "result": "根据分析，3月总销售额为$125,430，环比增长8.2%。以下为详细趋势图与周报内容...",
    "selected_skills": [
      {
        "skill_id": "uuid-1",
        "name": "excel-trend-analyzer",
        "contribution": "primary",
        "match_score": 0.952
      }
    ],
    "execution_meta": {
      "match_time_ms": 203,
      "llm_time_ms": 2340,
      "total_tokens_used": 1850,
      "request_id": "exec-abc-123"
    }
  },
  "request_id": "exec-abc-123"
}
```

**流式响应格式（`execution_config.stream=true`）**

Content-Type: `text/event-stream`
```
data: {"type":"match","matches":[...]}

data: {"type":"chunk","content":"根据分析"}

data: {"type":"chunk","content":"，3月总销售额为$125,430..."}

data: {"type":"done","meta":{"tokens_used":1850}}
```


## 6. 插件专用接口

### 6.1 获取热门技能（插件版）

**GET** `/plugin/skills/popular`
无需认证（但携带 JWT 可获取个性化推荐）。

**响应示例**
```json
{
  "code": 0,
  "data": {
    "items": [ ... ],
    "updated_at": "2026-04-23T10:00:00Z"
  }
}
```

### 6.2 获取技能下载信息

**GET** `/plugin/skills/{id}/download`
用于插件一键安装。

**响应示例**
```json
{
  "code": 0,
  "data": {
    "download_url": "https://github.com/anthropics/skills.git",
    "type": "git",           // git 或 zip
    "version": "1.2.0",
    "recommended_path": ".claude/skills/"
  }
}
```

### 6.3 获取用户已安装技能列表

**GET** `/plugin/user/installed`
需要认证：Bearer Token

**响应示例**
```json
{
  "code": 0,
  "data": {
    "skills": [
      {
        "skill_id": "uuid",
        "name": "excel-trend-analyzer",
        "installed_at": "2026-04-01T08:00:00Z",
        "local_version": "1.2.0",
        "auto_update": true
      }
    ]
  }
}
```

### 6.4 同步用户安装状态

**POST** `/plugin/user/installed`
需要认证：Bearer Token

**请求体**
```json
{
  "skill_id": "uuid",
  "action": "install",        // install / uninstall / update
  "local_version": "1.2.0"
}
```

### 6.5 上下文推荐

**POST** `/plugin/recommend`
需要认证：Bearer Token

根据 IDE 当前文件、光标位置推荐相关技能。

**请求体**
```json
{
  "file_path": "src/analysis/sales_report.py",
  "file_language": "python",
  "cursor_context": "df = pd.read_excel('sales.xlsx')\n# TODO: analyze trend",
  "top_k": 5
}
```

**响应示例**
```json
{
  "code": 0,
  "data": {
    "recommendations": [
      {
        "skill_id": "uuid",
        "name": "excel-trend-analyzer",
        "description": "分析Excel数据趋势",
        "reason": "检测到 Excel 文件读取操作"
      },
      {
        "skill_id": "uuid-2",
        "name": "python-code-optimizer",
        "description": "优化 Python 代码性能",
        "reason": "当前文件为 Python"
      }
    ]
  }
}
```


## 7. 管理接口（需管理员权限）

### 7.1 同步任务管理

#### 7.1.1 触发同步任务

**POST** `/admin/sync/trigger`
**请求体**
```json
{
  "type": "official",    // official / github_discover / single
  "repo_url": ""         // type=single 时必填
}
```

#### 7.1.2 获取同步任务列表

**GET** `/admin/sync/tasks`

#### 7.1.3 获取同步任务详情

**GET** `/admin/sync/tasks/{id}`

### 7.2 技能审核

#### 7.2.1 获取待审核技能列表

**GET** `/admin/skills/pending`

#### 7.2.2 审核技能

**POST** `/admin/skills/{id}/review`
**请求体**
```json
{
  "action": "approve",    // approve / reject
  "reason": "符合规范"
}
```

### 7.3 系统配置

#### 7.3.1 获取配置

**GET** `/admin/configs`

#### 7.3.2 更新配置

**PUT** `/admin/configs`
**请求体**（示例）
```json
{
  "github": {
    "min_stars": 100,
    "sync_interval_hours": 24
  },
  "embedding": {
    "model": "BAAI/bge-m3"
  }
}
```

### 7.4 统计看板

**GET** `/admin/stats/dashboard`
返回详细运营数据。


## 8. 附录

### 8.1 通用请求头

| 头部            | 说明                | 示例                   |
| --------------- | ------------------- | ---------------------- |
| `Authorization` | JWT 认证            | `Bearer eyJhbGciOi...` |
| `X-API-Key`     | API Key 认证        | `sk_live_xxxxxxxx`     |
| `Content-Type`  | 请求体类型          | `application/json`     |
| `Accept`        | 响应类型            | `application/json`     |
| `X-Request-ID`  | 请求追踪 ID（可选） | `client-gen-uuid`      |

### 8.2 速率限制

| 套餐   | 速率限制           |
| ------ | ------------------ |
| 免费版 | 100 次/天          |
| 专业版 | 10,000 次/天       |
| 企业版 | 无限制（可自定义） |

超过限制返回 HTTP 429 状态码。

### 8.3 示例代码（cURL）

**语义搜索**
```bash
curl -X POST https://api.skillhub.pro/v1/skills/search \
  -H "Content-Type: application/json" \
  -d '{"query": "分析Excel数据", "search_type": "semantic"}'
```

**智能路由执行**
```bash
curl -X POST https://api.skillhub.pro/v1/router/execute \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk_live_xxxx" \
  -d '{
    "query": "帮我生成周报",
    "execution_config": {"model": "claude-3.5-sonnet"}
  }'
```

---

以上为 SkillHub Pro 平台的完整接口定义文档，开发团队应基于此规范进行前后端开发和联调。