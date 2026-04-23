# SkillHub Pro 概要设计文档

## 文档信息

| 项目     | 内容                             |
| -------- | -------------------------------- |
| 项目名称 | SkillHub Pro（技能宝库）         |
| 文档类型 | 概要设计说明书                   |
| 版本     | V1.0                             |
| 创建日期 | 2026-04-23                       |
| 目标读者 | 开发团队、技术评审、项目管理人员 |


## 1. 引言

### 1.1 编写目的

本文档旨在对 SkillHub Pro 系统进行概要设计，明确系统的整体架构、模块划分、技术选型、接口规范及部署方案，为后续的详细设计和编码实现提供指导和依据。

### 1.2 设计目标

- **高可用**：系统整体可用性 ≥ 99.5%，支持水平扩展。
- **高性能**：API 响应 P95 < 200ms，路由匹配 P95 < 500ms。
- **可扩展**：模块化设计，支持未来增加多 IDE 插件、多 LLM 网关。
- **安全性**：API 认证、技能安全扫描、防注入攻击。
- **开放性**：提供标准 RESTful API，支持第三方集成。

### 1.3 参考资料

- 《SkillHub Pro 产品需求说明文档 V3.0》
- 《VS Code Extension API 官方文档》
- 《Milvus 向量数据库文档》
- 《Gin Web Framework 文档》


## 2. 总体设计

### 2.1 系统逻辑架构

系统采用分层微服务架构，分为**客户端层、接入层、业务服务层、基础设施层**四层。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              客户端层（Client Layer）                         │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────────────┐  │
│  │   Web 前端       │  │  VS Code 插件    │  │  第三方 API 客户端       │  │
│  │  (Vue3 + TS)     │  │  (TypeScript)    │  │  (REST / SDK)           │  │
│  └──────────────────┘  └──────────────────┘  └──────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              接入层（Gateway Layer）                          │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                          API Gateway (Nginx / Kong)                     │ │
│  │               路由转发 | 负载均衡 | SSL终结 | 限流熔断 | 日志审计          │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           业务服务层（Business Service Layer）                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌────────┐ │
│  │ 技能管理服务 │ │ 用户认证服务 │ │  搜索服务   │ │  插件服务   │ │统计服务│ │
│  │ Skill Service│ │ Auth Service│ │Search Service│ │Plugin Service│ │Stats   │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ └────────┘ │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                    智能路由服务 (Intelligent Router Service)              │ │
│  │  意图识别 | 混合检索 | 重排序 | 上下文组装 | LLM 调用网关                    │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                    技能同步服务 (Skill Sync Service)                      │ │
│  │  GitHub 爬虫 | 仓库解析 | 向量化任务 | 安全扫描                            │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          基础设施层（Infrastructure Layer）                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────┐ │
│  │PostgreSQL│ │  Redis   │ │Meilisearch││  MinIO   │ │  Milvus  │ │RabbitMQ│ │
│  │ (主数据)  │ │ (缓存)   │ │(全文检索) │ │(文件存储) │ │(向量库)  │ │(消息)  │ │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ └───────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 核心数据流

#### 2.2.1 技能同步数据流

```
GitHub 仓库 / 官方源 → [爬虫服务] → SKILL.md 解析 → 元数据提取
                                              ↓
                                    LLM 增强描述生成 → 向量化
                                              ↓
                              PostgreSQL (元数据) + MinIO (源文件) + Milvus (向量)
```

#### 2.2.2 智能路由数据流

```
用户自然语言查询 → [路由服务] → Embedding 向量化
                    ↓
            Milvus 向量检索 (Top 50) + Meilisearch 关键词检索 (Top 50)
                    ↓
                RRF 融合排序 (Top 50)
                    ↓
            Cross-Encoder 重排序 (Top 3)
                    ↓
            从 MinIO 读取 SKILL.md → 上下文组装
                    ↓
            LLM 执行 → 返回结果
```

#### 2.2.3 VS Code 插件交互流

```
VS Code 插件 → [HTTPS] → API Gateway → 插件服务 / 技能服务
                    ↓
            用户登录 (OAuth) → JWT 颁发 → 本地加密存储
                    ↓
            技能安装：git clone / 下载 ZIP → 本地文件系统
                    ↓
            状态同步：POST /api/v1/plugin/user/installed
```


## 3. 模块设计

### 3.1 技能管理服务（Skill Service）

**职责**：技能的 CRUD、分类管理、版本管理、元数据维护。

**核心接口**：
| 方法 | 路径                    | 描述                 |
| ---- | ----------------------- | -------------------- |
| GET  | `/api/v1/skills`        | 分页获取技能列表     |
| GET  | `/api/v1/skills/{id}`   | 获取技能详情         |
| POST | `/api/v1/skills/search` | 全文搜索技能         |
| GET  | `/api/v1/categories`    | 获取所有分类         |
| POST | `/api/v1/skills/submit` | 提交新技能（需认证） |

**内部依赖**：PostgreSQL（主存储）、Meilisearch（搜索索引）、MinIO（文件读取）。

**设计要点**：
- 使用 GORM 作为 ORM，定义 Skill 模型与数据库表映射。
- 搜索查询同时走 Meilisearch，支持分词与拼音。
- 列表接口支持按分类、来源、星数排序。

### 3.2 技能同步服务（Skill Sync Service）

**职责**：从官方源和 GitHub 发现并同步技能，执行安全扫描和向量化。

**核心组件**：
- **GitHub Crawler**：基于 Topics、路径、Awesome 列表发现仓库。
- **Parser**：解析 SKILL.md，提取 YAML frontmatter 和正文。
- **Security Scanner**：调用 Semgrep 检测危险模式。
- **Embedding Worker**：异步生成增强描述并向量化。

**调度策略**：
- 官方仓库：每日 02:00 UTC 增量同步。
- GitHub 发现：每周日 03:00 UTC 全量扫描，每日增量更新。
- 用户提交：触发即时同步。

**内部数据流**：
```
Crawler → Parser → Security Scanner → Skill DB
                            ↘
                      Embedding Worker → Milvus
```

### 3.3 智能路由服务（Router Service）

**职责**：根据用户自然语言查询匹配最优技能，并支持端到端执行。

**核心子模块**：
| 子模块            | 功能                      | 技术实现                        |
| ----------------- | ------------------------- | ------------------------------- |
| Query Embedding   | 将查询文本向量化          | 调用 Embedding 模型 (BGE-M3)    |
| Hybrid Retriever  | 向量检索 + 关键词检索融合 | Milvus ANN + Meilisearch + RRF  |
| Reranker          | 对候选集精排              | Cross-Encoder (bge-reranker-v2) |
| Context Assembler | 组装发送给 LLM 的提示词   | 模板引擎 + SKILL.md 加载        |
| LLM Gateway       | 统一调用多种大模型        | LiteLLM 或自研适配层            |

**核心接口**：
| 方法 | 路径                       | 描述                    |
| ---- | -------------------------- | ----------------------- |
| POST | `/api/v1/router/match`     | 仅返回匹配的技能信息    |
| POST | `/api/v1/router/execute`   | 端到端执行并返回结果    |
| POST | `/api/v1/plugin/recommend` | 根据 IDE 上下文推荐技能 |

**设计要点**：
- 路由匹配结果可缓存（相同查询短时间内返回缓存）。
- 执行接口支持流式输出（SSE）。
- 所有路由调用记录存入 `router_logs` 表用于分析和优化。

### 3.4 用户认证服务（Auth Service）

**职责**：用户注册、登录、API Key 管理、JWT 签发与验证。

**认证方式**：
- 邮箱 + 密码（bcrypt 加密）
- GitHub OAuth 2.0

**API Key 管理**：
- 每个用户可生成多个 API Key，支持设置过期时间。
- API Key 使用 HMAC-SHA256 签名进行请求验证。

**JWT 设计**：
- 用于插件端会话保持，有效期 7 天。
- 存储在 VS Code SecretStorage 中。

### 3.5 插件服务（Plugin Service）

**职责**：为 VS Code 插件提供专用 API，包括热门推荐、下载地址、状态同步。

**核心接口**：
| 方法 | 路径                                  | 描述                               |
| ---- | ------------------------------------- | ---------------------------------- |
| GET  | `/api/v1/plugin/skills/popular`       | 获取热门技能                       |
| GET  | `/api/v1/plugin/skills/{id}/download` | 获取技能下载信息（Git URL 或 ZIP） |
| GET  | `/api/v1/plugin/user/installed`       | 获取用户已安装技能列表             |
| POST | `/api/v1/plugin/user/installed`       | 同步安装/卸载行为                  |
| POST | `/api/v1/plugin/recommend`            | IDE 上下文推荐                     |

**设计要点**：
- 下载信息优先返回 GitHub 仓库地址，次选 MinIO 预签名 URL。
- 状态同步接口需验证 JWT，更新 `user_installed_skills` 表。

### 3.6 统计服务（Stats Service）

**职责**：收集和聚合平台各类统计数据，提供看板数据。

**统计维度**：
- 技能总数、每日新增
- API 调用量（按 Key、按接口）
- 插件安装量、活跃用户
- 路由调用次数、匹配准确率（基于用户反馈）

**实现方案**：
- 使用 Prometheus 采集指标，Grafana 可视化。
- 业务统计通过分析日志和数据库记录定时计算。


## 4. 数据存储设计

### 4.1 关系型数据库（PostgreSQL）

**核心表**：
- `skills`：技能主表（详见 PRD 细化部分）
- `categories`：分类表
- `users`：用户表
- `api_keys`：API Key 表
- `user_installed_skills`：用户安装记录
- `router_logs`：路由调用日志

**读写分离**：
- 主库负责写操作，从库负责读操作（列表查询、统计报表）。

### 4.2 缓存（Redis）

**用途**：
- 热门技能榜单缓存（TTL 1 小时）
- 技能详情缓存（TTL 10 分钟）
- 路由匹配结果缓存（TTL 30 分钟，仅针对相同查询）
- API 限流计数器
- 会话存储（用户登录状态）

### 4.3 全文检索引擎（Meilisearch）

**索引字段**：
- `name`, `description`, `enhanced_description`, `category`, `tags`

**索引更新**：
- 技能创建/更新时异步同步。

### 4.4 向量数据库（Milvus）

**Collection Schema**：
| 字段          | 类型              | 说明                |
| ------------- | ----------------- | ------------------- |
| `id`          | Int64             | 自增主键            |
| `skill_id`    | VarChar           | 关联 skills 表 UUID |
| `embedding`   | FloatVector(1024) | BGE-M3 向量         |
| `name`        | VarChar           | 技能名称            |
| `category`    | VarChar           | 分类                |
| `source_type` | VarChar           | 来源类型            |

**索引类型**：IVF_FLAT 或 HNSW，根据数据量选择。

### 4.5 对象存储（MinIO / OSS）

**Bucket 设计**：
- `skills-source/`：存储 SKILL.md 源文件，路径为 `{skill_id}/SKILL.md`。
- `user-uploads/`：用户通过 API 上传的临时文件（路由执行上下文）。


## 5. 接口设计

### 5.1 API 设计规范

- **协议**：HTTPS
- **数据格式**：JSON
- **认证方式**：
  - 公开接口：无需认证
  - 用户接口：Bearer Token (JWT)
  - 管理接口：API Key + HMAC 签名
- **错误码规范**：
  - `0`：成功
  - `1001`：参数错误
  - `1002`：认证失败
  - `1003`：权限不足
  - `2001`：资源不存在
  - `5000`：服务器内部错误

### 5.2 核心接口列表

（详见 PRD 细化文档，此处仅列概要）

| 分组 | 接口                           | 方法     | 说明         |
| ---- | ------------------------------ | -------- | ------------ |
| 技能 | `/skills`                      | GET      | 列表         |
| 技能 | `/skills/{id}`                 | GET      | 详情         |
| 技能 | `/skills/search`               | POST     | 搜索         |
| 路由 | `/router/match`                | POST     | 匹配技能     |
| 路由 | `/router/execute`              | POST     | 执行技能     |
| 插件 | `/plugin/skills/popular`       | GET      | 热门技能     |
| 插件 | `/plugin/skills/{id}/download` | GET      | 下载信息     |
| 插件 | `/plugin/user/installed`       | GET/POST | 用户安装同步 |
| 插件 | `/plugin/recommend`            | POST     | 上下文推荐   |
| 用户 | `/user/login`                  | POST     | 登录         |
| 用户 | `/user/apikeys`                | GET/POST | API Key 管理 |
| 管理 | `/admin/sync/trigger`          | POST     | 触发同步     |
| 管理 | `/admin/stats`                 | GET      | 统计数据     |


## 6. VS Code 插件设计

### 6.1 插件结构

```
vscode-skillhub/
├── src/
│   ├── extension.ts          # 入口，激活事件
│   ├── views/
│   │   ├── skillTreeView.ts  # 侧边栏树视图
│   │   └── webview/          # 详情页 Webview 内容
│   ├── services/
│   │   ├── apiClient.ts      # 后端 API 调用
│   │   ├── skillManager.ts   # 本地技能安装、管理
│   │   └── authService.ts    # 登录状态管理
│   ├── utils/
│   └── types/
├── package.json
└── tsconfig.json
```

### 6.2 核心类设计

**SkillManager 类**：
- `async install(skill: Skill): Promise<void>`
- `getInstalledSkills(): LocalSkill[]`
- `toggleSkill(id: string, enabled: boolean): void`
- `checkForUpdates(): Promise<UpdateInfo[]>`

**APIClient 类**：
- 封装 Axios，自动添加 JWT 认证头。
- 处理 401 自动跳转登录。

**AuthService 类**：
- 使用 `vscode.authentication.getSession('github', ...)` 获取 GitHub Token。
- 与后端交换 JWT，存入 `SecretStorage`。

### 6.3 激活事件与视图

- 激活事件：`onView:skillhub-sidebar`，即用户点击侧边栏图标时激活。
- 侧边栏视图：使用 `vscode.window.createTreeView` 构建树形列表。
- 详情 Webview：使用 `vscode.window.createWebviewPanel` 渲染技能详情，支持 Vue 组件（通过 iframe 或编译为独立 JS）。


## 7. 部署与运维设计

### 7.1 部署架构

**开发环境**：
- Docker Compose 一键启动所有服务（PostgreSQL、Redis、Meilisearch、Milvus、MinIO、后端、前端）。

**生产环境**：
- Kubernetes 集群部署。
- 各服务独立 Deployment，通过 Service 暴露。
- Ingress Nginx 作为统一入口，配置 SSL 证书。
- 使用 HPA 根据 CPU/Memory 自动扩缩容。

### 7.2 容器化清单

| 服务        | 镜像                        | 端口  |
| ----------- | --------------------------- | ----- |
| Backend     | `skillhub-backend:latest`   | 8080  |
| Frontend    | `nginx:alpine` 打包静态文件 | 80    |
| PostgreSQL  | `postgres:15`               | 5432  |
| Redis       | `redis:7-alpine`            | 6379  |
| Meilisearch | `getmeili/meilisearch:v1.5` | 7700  |
| Milvus      | `milvusdb/milvus:latest`    | 19530 |
| MinIO       | `minio/minio:latest`        | 9000  |

### 7.3 监控与告警

- **Prometheus**：采集各服务指标（QPS、延迟、错误率）。
- **Grafana**：展示业务和系统监控大盘。
- **AlertManager**：配置 API 错误率 > 1%、服务宕机等告警规则。
- **日志**：使用 Loki + Grafana 聚合查询。

### 7.4 数据备份策略

- PostgreSQL：每日全量备份 + WAL 归档。
- MinIO：每日同步至异地 OSS。
- Milvus：定期导出元数据快照。


## 8. 安全设计

### 8.1 网络安全

- 所有服务间通信走内网，仅 API Gateway 暴露公网。
- 启用 WAF 防护常见 Web 攻击。

### 8.2 应用安全

- API 请求强制 HTTPS。
- API Key 使用 HMAC-SHA256 签名，防重放攻击（timestamp + nonce）。
- JWT 有效期较短（7天），Refresh Token 机制。
- 输入参数严格校验，防 SQL 注入（ORM 参数化）。
- 路由执行 API 的 `context` 字段做严格过滤，禁止注入恶意指令。

### 8.3 技能安全

- 技能入库前进行 Semgrep 扫描，检测：
  - `os.system`、`subprocess` 等危险调用。
  - 硬编码密钥。
  - Prompt 注入模式（如 “忽略之前指令”）。
- 标记为 `warning` 的技能默认不在热门推荐中展示，需用户确认后安装。


## 9. 风险与技术应对

| 风险点          | 技术应对                                             |
| --------------- | ---------------------------------------------------- |
| GitHub API 限流 | 多 Token 轮换池 + 请求队列 + 增量同步                |
| 向量检索延迟高  | 使用 Milvus GPU 加速版或 Qdrant；热门查询缓存        |
| 路由准确率不足  | 建立人工评测集，持续迭代 Embedding/Reranker 模型     |
| LLM 调用成本高  | 对重复查询缓存；使用小模型做预筛选；按套餐限制调用量 |
| 插件兼容性      | CI 中加入 Windows/macOS/Linux 测试矩阵               |


## 10. 附录

### 10.1 开发环境搭建指南（简要）

```bash
# 1. 克隆代码
git clone https://github.com/wjames2000/skillhub-pro.git

# 2. 启动依赖服务
docker-compose -f docker-compose.dev.yml up -d

# 3. 后端
cd backend
go mod download
go run cmd/server/main.go

# 4. 前端
cd web
pnpm install
pnpm dev

# 5. 插件
cd vscode-extension
npm install
npm run watch
# 按 F5 启动调试
```

### 10.2 关键配置项示例

**后端配置（config.yaml）**：
```yaml
server:
  port: 8080
database:
  host: postgres
  port: 5432
  name: skillhub
  user: skillhub
  password: ${DB_PASSWORD}
redis:
  addr: redis:6379
meilisearch:
  host: http://meilisearch:7700
  api_key: ${MEILI_KEY}
milvus:
  host: milvus
  port: 19530
github:
  tokens:
    - ${GITHUB_TOKEN_1}
    - ${GITHUB_TOKEN_2}
embedding:
  model: BAAI/bge-m3
  endpoint: https://api.openai.com/v1/embeddings  # 或自建服务
llm:
  default_model: claude-3.5-sonnet
  api_key: ${ANTHROPIC_API_KEY}
```

---

以上为 SkillHub Pro 系统的概要设计文档。该文档涵盖了系统架构、模块划分、数据存储、接口规范、插件设计及部署运维等方面，可作为详细设计和编码实现的基线。