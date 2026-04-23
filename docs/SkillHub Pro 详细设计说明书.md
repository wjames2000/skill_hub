# SkillHub Pro 详细设计说明书

## 文档信息

| 项目     | 内容                                                         |
| -------- | ------------------------------------------------------------ |
| 项目名称 | SkillHub Pro（技能宝库）                                     |
| 文档类型 | 详细设计说明书                                               |
| 版本     | V1.1                                                         |
| 更新日期 | 2026-04-23                                                   |
| 目标读者 | 开发工程师、测试工程师                                       |
| 参考文档 | 《SkillHub Pro PRD V3.0终稿》《SkillHub Pro 概要设计 V1.0》  |
| 主要变更 | ORM 由 GORM 切换为 XORM；说明 Milvus 依赖协调服务的原因并支持 Consul 替代 ETCD |


## 1. 引言

### 1.1 编写目的

本文档在概要设计的基础上，对系统各模块进行详细设计，明确每个模块的内部结构、类设计、接口详细定义、数据库表结构、核心算法实现及异常处理策略，为编码实现提供精确的指导。

### 1.2 设计范围

- 后端服务：技能管理服务、技能同步服务、智能路由服务、用户认证服务、插件服务、统计服务
- 数据层：PostgreSQL表结构详细定义、索引设计、Redis数据结构、Milvus Collection Schema、Meilisearch索引配置
- VS Code插件：核心类详细设计、API调用时序、本地存储结构
- 部署细节：Dockerfile、Kubernetes资源配置示例

### 1.3 关键设计决策说明

**1. ORM 选型：XORM**

本项目选用 XORM 作为 ORM 框架，原因如下：
- **轻量级**：XORM 比 GORM 更轻量，性能和内存占用略优。
- **结构体标签兼容性好**：XORM 的 tag 设计与 Go 原生结构体标签风格一致，便于模型定义。
- **读写分离支持完善**：XORM 内置支持多数据源和读写分离，适合本项目读多写少的场景。
- **丰富的扩展能力**：支持自定义缓存、事件钩子，便于后续扩展。

**2. 服务协调与配置中心：Consul**

在部署架构中，Milvus 向量数据库默认依赖 ETCD 作为元数据存储和协调服务。考虑到团队已有 Consul 运维经验，且 Consul 同时提供服务发现和 KV 存储能力，本项目统一采用 Consul 替代 ETCD，理由如下：
- **统一技术栈**：后端微服务已计划使用 Consul 进行服务注册与发现，复用可降低运维复杂度。
- **兼容性**：Milvus 2.3+ 支持通过配置将元数据存储指向 Consul（通过 `metastore.type=consul` 和相应地址配置）。
- **功能满足**：Consul 提供的 KV 存储完全满足 Milvus 元数据管理需求。


## 2. 后端模块详细设计

### 2.1 整体代码结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # 程序入口
├── internal/
│   ├── domain/                     # 领域模型（实体、值对象）
│   │   ├── skill.go
│   │   ├── category.go
│   │   ├── user.go
│   │   └── router_log.go
│   ├── repository/                 # 数据仓储接口及实现
│   │   ├── skill_repo.go
│   │   ├── user_repo.go
│   │   └── ...
│   ├── service/                    # 业务逻辑层
│   │   ├── skill_service.go
│   │   ├── sync_service.go
│   │   ├── router_service.go
│   │   ├── auth_service.go
│   │   └── plugin_service.go
│   ├── handler/                    # HTTP 处理器（Controller）
│   │   ├── skill_handler.go
│   │   ├── router_handler.go
│   │   ├── auth_handler.go
│   │   └── plugin_handler.go
│   ├── middleware/                 # 中间件
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── ratelimit.go
│   │   └── logging.go
│   ├── crawler/                    # 爬虫模块
│   │   ├── github_crawler.go
│   │   ├── parser.go
│   │   └── discoverer.go
│   ├── embedding/                  # 向量化模块
│   │   ├── client.go
│   │   └── worker.go
│   ├── llm/                        # LLM 调用网关
│   │   └── gateway.go
│   ├── config/                     # 配置管理
│   │   └── config.go
│   └── pkg/                        # 内部公共库
│       ├── logger/
│       ├── errors/
│       └── utils/
├── pkg/                            # 可对外暴露的公共库
├── migrations/                     # 数据库迁移文件
├── scripts/                        # 脚本
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

### 2.2 领域模型详细设计（基于 XORM）

XORM 使用 `xorm` 标签定义字段映射，通过 `engine.Sync2()` 自动同步表结构。以下是核心实体的定义。

#### 2.2.1 Skill 实体

```go
// internal/domain/skill.go
package domain

import (
    "time"
    "encoding/json"
    "github.com/google/uuid"
    "xorm.io/xorm"
)

type Skill struct {
    ID               uuid.UUID       `xorm:"pk 'id'" json:"id"`
    Name             string          `xorm:"varchar(255) notnull" json:"name"`
    Description      string          `xorm:"text" json:"description"`
    SourceType       string          `xorm:"varchar(50)" json:"source_type"`        // official, github, community
    SourceURL        string          `xorm:"varchar(500)" json:"source_url"`
    GitHubStars      int             `xorm:"default 0" json:"github_stars"`
    GitHubForks      int             `xorm:"default 0" json:"github_forks"`
    CategoryID       uuid.UUID       `xorm:"category_id" json:"category_id"`
    Category         *Category       `xorm:"-" json:"category,omitempty"`           // 非数据库字段
    Author           string          `xorm:"varchar(255)" json:"author"`
    Version          string          `xorm:"varchar(50)" json:"version"`
    DownloadCount    int             `xorm:"default 0" json:"download_count"`
    ViewCount        int             `xorm:"default 0" json:"view_count"`
    RatingAvg        float32         `xorm:"decimal(2,1)" json:"rating_avg"`
    RatingCount      int             `xorm:"default 0" json:"rating_count"`
    SecurityStatus   string          `xorm:"varchar(20) default 'pending'" json:"security_status"` // pending, scanning, safe, warning
    SkillMDContent   string          `xorm:"text" json:"skill_md_content"`
    SkillMDURL       string          `xorm:"varchar(500)" json:"skill_md_url"`
    Metadata         JSONMap         `xorm:"jsonb" json:"metadata"`
    EnhancedDesc     string          `xorm:"text" json:"enhanced_description"`
    EmbeddingVersion string          `xorm:"varchar(20)" json:"embedding_version"`
    LastEmbeddedAt   *time.Time      `json:"last_embedded_at"`
    CreatedAt        time.Time       `xorm:"created" json:"created_at"`
    UpdatedAt        time.Time       `xorm:"updated" json:"updated_at"`
    DeletedAt        time.Time       `xorm:"deleted" json:"-"`
}

// TableName 自定义表名
func (Skill) TableName() string {
    return "skills"
}

// BeforeInsert 钩子，自动生成 UUID
func (s *Skill) BeforeInsert() {
    if s.ID == uuid.Nil {
        s.ID = uuid.New()
    }
}

// JSONMap 用于存储 JSONB 字段（XORM 会自动处理 json.RawMessage 或自定义类型）
type JSONMap map[string]interface{}

func (j JSONMap) FromDB(data []byte) error {
    return json.Unmarshal(data, &j)
}

func (j JSONMap) ToDB() ([]byte, error) {
    return json.Marshal(j)
}
```

#### 2.2.2 User 与 APIKey 实体

```go
// internal/domain/user.go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `xorm:"pk 'id'" json:"id"`
    Email        string    `xorm:"varchar(255) unique" json:"email"`
    PasswordHash string    `xorm:"varchar(255)" json:"-"`
    GitHubID     string    `xorm:"varchar(100) unique" json:"github_id"`
    AvatarURL    string    `xorm:"varchar(500)" json:"avatar_url"`
    Name         string    `xorm:"varchar(100)" json:"name"`
    IsActive     bool      `xorm:"default true" json:"is_active"`
    IsAdmin      bool      `xorm:"default false" json:"is_admin"`
    CreatedAt    time.Time `xorm:"created" json:"created_at"`
    UpdatedAt    time.Time `xorm:"updated" json:"updated_at"`
    
    APIKeys      []APIKey  `xorm:"-" json:"api_keys,omitempty"`
}

func (User) TableName() string {
    return "users"
}

func (u *User) BeforeInsert() {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
}

type APIKey struct {
    ID         uuid.UUID  `xorm:"pk 'id'" json:"id"`
    UserID     uuid.UUID  `xorm:"index" json:"user_id"`
    KeyPrefix  string     `xorm:"varchar(20)" json:"key_prefix"`
    KeyHash    string     `xorm:"varchar(255)" json:"-"`
    Name       string     `xorm:"varchar(100)" json:"name"`
    LastUsedAt *time.Time `json:"last_used_at"`
    ExpiresAt  *time.Time `json:"expires_at"`
    IsActive   bool       `xorm:"default true" json:"is_active"`
    CreatedAt  time.Time  `xorm:"created" json:"created_at"`
}

func (APIKey) TableName() string {
    return "api_keys"
}

func (k *APIKey) BeforeInsert() {
    if k.ID == uuid.Nil {
        k.ID = uuid.New()
    }
}
```

#### 2.2.3 RouterLog 实体

```go
type RouterLog struct {
    ID               uuid.UUID  `xorm:"pk 'id'" json:"id"`
    UserID           *uuid.UUID `xorm:"index" json:"user_id"`
    APIKeyID         *uuid.UUID `xorm:"index" json:"api_key_id"`
    Query            string     `xorm:"text notnull" json:"query"`
    MatchedSkillIDs  JSONArray  `xorm:"jsonb" json:"matched_skill_ids"`
    SelectedSkillID  *uuid.UUID `json:"selected_skill_id"`
    ExecutionSuccess *bool      `json:"execution_success"`
    UserFeedback     *int       `xorm:"smallint" json:"user_feedback"`
    LatencyMs        int        `json:"latency_ms"`
    TokensUsed       int        `json:"tokens_used"`
    ErrorMessage     string     `xorm:"text" json:"error_message"`
    CreatedAt        time.Time  `xorm:"created index" json:"created_at"`
}
```

### 2.3 数据仓储层详细设计（基于 XORM）

#### 2.3.1 仓储接口定义

```go
// internal/repository/skill_repo.go
package repository

import (
    "context"
    "github.com/google/uuid"
    "skillhub/internal/domain"
)

type SkillRepository interface {
    Create(ctx context.Context, skill *domain.Skill) error
    Update(ctx context.Context, skill *domain.Skill) (int64, error)
    Delete(ctx context.Context, id uuid.UUID) (int64, error)
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Skill, error)
    GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Skill, error)
    GetBySourceURL(ctx context.Context, sourceURL string) (*domain.Skill, error)
    List(ctx context.Context, filter *SkillFilter) ([]*domain.Skill, int64, error)
    IncrementViewCount(ctx context.Context, id uuid.UUID) error
    IncrementDownloadCount(ctx context.Context, id uuid.UUID) error
    BatchCreate(ctx context.Context, skills []*domain.Skill) error
    UpdateEmbeddingInfo(ctx context.Context, id uuid.UUID, version string) error
}

type SkillFilter struct {
    CategoryID     *uuid.UUID
    SourceType     *string
    SearchQuery    string
    SortBy         string
    SortOrder      string
    Page           int
    PageSize       int
    MinStars       int
    SecurityStatus *string
}
```

#### 2.3.2 XORM 实现

```go
// internal/repository/skill_repo_xorm.go
package repository

import (
    "context"
    "github.com/google/uuid"
    "xorm.io/xorm"
    "skillhub/internal/domain"
)

type skillRepoXorm struct {
    engine *xorm.Engine
    // 可选：读写分离引擎组
    // eg *xorm.EngineGroup
}

func NewSkillRepository(engine *xorm.Engine) SkillRepository {
    return &skillRepoXorm{engine: engine}
}

func (r *skillRepoXorm) Create(ctx context.Context, skill *domain.Skill) error {
    _, err := r.engine.Context(ctx).Insert(skill)
    return err
}

func (r *skillRepoXorm) Update(ctx context.Context, skill *domain.Skill) (int64, error) {
    return r.engine.Context(ctx).ID(skill.ID).Update(skill)
}

func (r *skillRepoXorm) GetByID(ctx context.Context, id uuid.UUID) (*domain.Skill, error) {
    skill := &domain.Skill{}
    has, err := r.engine.Context(ctx).ID(id).Get(skill)
    if err != nil {
        return nil, err
    }
    if !has {
        return nil, nil
    }
    return skill, nil
}

func (r *skillRepoXorm) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Skill, error) {
    var skills []*domain.Skill
    err := r.engine.Context(ctx).In("id", ids).Find(&skills)
    return skills, err
}

func (r *skillRepoXorm) List(ctx context.Context, filter *SkillFilter) ([]*domain.Skill, int64, error) {
    session := r.engine.Context(ctx).Table(&domain.Skill{})
    
    if filter.CategoryID != nil && *filter.CategoryID != uuid.Nil {
        session = session.Where("category_id = ?", filter.CategoryID)
    }
    if filter.SourceType != nil && *filter.SourceType != "" {
        session = session.Where("source_type = ?", *filter.SourceType)
    }
    if filter.MinStars > 0 {
        session = session.Where("github_stars >= ?", filter.MinStars)
    }
    if filter.SecurityStatus != nil && *filter.SecurityStatus != "" {
        session = session.Where("security_status = ?", *filter.SecurityStatus)
    }

    // 总数
    total, err := session.Clone().Count(&domain.Skill{})
    if err != nil {
        return nil, 0, err
    }

    // 排序
    orderField := "created_at"
    switch filter.SortBy {
    case "stars":
        orderField = "github_stars"
    case "downloads":
        orderField = "download_count"
    case "updated":
        orderField = "updated_at"
    }
    if filter.SortOrder == "asc" {
        session = session.Asc(orderField)
    } else {
        session = session.Desc(orderField)
    }

    // 分页
    offset := (filter.Page - 1) * filter.PageSize
    var skills []*domain.Skill
    err = session.Limit(filter.PageSize, offset).Find(&skills)
    return skills, total, err
}

func (r *skillRepoXorm) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
    _, err := r.engine.Context(ctx).ID(id).Incr("view_count").Update(&domain.Skill{})
    return err
}

func (r *skillRepoXorm) BatchCreate(ctx context.Context, skills []*domain.Skill) error {
    // XORM 支持批量插入
    _, err := r.engine.Context(ctx).Insert(skills)
    return err
}

func (r *skillRepoXorm) UpdateEmbeddingInfo(ctx context.Context, id uuid.UUID, version string) error {
    skill := &domain.Skill{EmbeddingVersion: version}
    _, err := r.engine.Context(ctx).ID(id).Cols("embedding_version", "last_embedded_at").Update(skill)
    return err
}
```

#### 2.3.3 数据库引擎初始化（含读写分离）

```go
// internal/config/database.go
package config

import (
    "time"
    _ "github.com/lib/pq"
    "xorm.io/xorm"
    "xorm.io/xorm/log"
)

func InitDB(cfg *Config) (*xorm.Engine, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
    
    engine, err := xorm.NewEngine("postgres", dsn)
    if err != nil {
        return nil, err
    }
    
    // 设置连接池
    engine.SetMaxIdleConns(10)
    engine.SetMaxOpenConns(100)
    engine.SetConnMaxLifetime(time.Hour)
    
    // 开发环境显示SQL
    if cfg.Env == "development" {
        engine.ShowSQL(true)
        engine.Logger().SetLevel(log.LOG_DEBUG)
    }
    
    // 自动同步表结构（开发环境）
    if cfg.Env == "development" {
        err = engine.Sync2(
            &domain.Skill{},
            &domain.Category{},
            &domain.User{},
            &domain.APIKey{},
            &domain.RouterLog{},
        )
        if err != nil {
            return nil, err
        }
    }
    
    return engine, nil
}

// 读写分离配置（可选）
func InitDBGroup(cfg *Config) (*xorm.EngineGroup, error) {
    masters := []string{cfg.DB.MasterDSN}
    slaves := cfg.DB.SlaveDSNs
    return xorm.NewEngineGroup("postgres", masters, slaves, xorm.RoundRobinPolicy())
}
```

### 2.4 服务层设计

服务层设计保持与之前一致，只需将依赖的 `repository.SkillRepository` 接口的具体实现切换为 XORM 版本即可，无需修改业务逻辑代码。

```golang
// internal/service/router_service.go
package service

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/google/uuid"
    "skillhub/internal/domain"
    "skillhub/internal/repository"
    "skillhub/internal/embedding"
    "skillhub/internal/llm"
    "skillhub/pkg/logger"
)

type RouterService interface {
    Match(ctx context.Context, req *MatchRequest) (*MatchResponse, error)
    Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error)
}

type MatchRequest struct {
    Query     string                 `json:"query"`
    TopK      int                    `json:"top_k"`
    Filters   map[string]interface{} `json:"filters"`
    IncludeContent bool              `json:"include_content"`
}

type MatchResponse struct {
    Matches []SkillMatch `json:"matches"`
    Meta    MatchMeta    `json:"meta"`
}

type SkillMatch struct {
    SkillID     string  `json:"skill_id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Score       float32 `json:"match_score"`
    SkillMDURL  string  `json:"skill_md_url,omitempty"`
}

type MatchMeta struct {
    EmbeddingTimeMs int `json:"embedding_time_ms"`
    VectorSearchMs  int `json:"vector_search_time_ms"`
    RerankTimeMs    int `json:"rerank_time_ms"`
    TotalTimeMs     int `json:"total_time_ms"`
}

type ExecuteRequest struct {
    Query           string                 `json:"query"`
    TopK            int                    `json:"top_k"`
    ExecutionConfig ExecutionConfig        `json:"execution_config"`
    Context         map[string]interface{} `json:"context"`
    Files           []AttachedFile         `json:"files,omitempty"`
}

type ExecutionConfig struct {
    Model       string  `json:"model"`
    MaxTokens   int     `json:"max_tokens"`
    Temperature float32 `json:"temperature"`
    Stream      bool    `json:"stream"`
}

type AttachedFile struct {
    Name string `json:"name"`
    URL  string `json:"url"`
}

type ExecuteResponse struct {
    Result         string          `json:"result"`
    SelectedSkills []SelectedSkill `json:"selected_skills"`
    Meta           ExecutionMeta   `json:"execution_meta"`
}

type SelectedSkill struct {
    SkillID      string  `json:"skill_id"`
    Name         string  `json:"name"`
    Contribution string  `json:"contribution"` // primary, secondary
    MatchScore   float32 `json:"match_score"`
}

type ExecutionMeta struct {
    MatchTimeMs  int    `json:"match_time_ms"`
    LLMTimeMs    int    `json:"llm_time_ms"`
    TokensUsed   int    `json:"total_tokens_used"`
    RequestID    string `json:"request_id"`
}

// 实现
type routerService struct {
    skillRepo      repository.SkillRepository
    vectorRepo     repository.VectorRepository      // Milvus 操作封装
    searchRepo     repository.SearchRepository
    embeddingClient embedding.Client
    rerankerClient  embedding.RerankerClient
    llmGateway      llm.Gateway
    cache           Cache
}

func NewRouterService(
    skillRepo repository.SkillRepository,
    vectorRepo repository.VectorRepository,
    searchRepo repository.SearchRepository,
    embeddingClient embedding.Client,
    rerankerClient embedding.RerankerClient,
    llmGateway llm.Gateway,
    cache Cache,
) RouterService {
    return &routerService{
        skillRepo:       skillRepo,
        vectorRepo:      vectorRepo,
        searchRepo:      searchRepo,
        embeddingClient: embeddingClient,
        rerankerClient:  rerankerClient,
        llmGateway:      llmGateway,
        cache:           cache,
    }
}

func (s *routerService) Match(ctx context.Context, req *MatchRequest) (*MatchResponse, error) {
    start := time.Now()
    meta := MatchMeta{}
    
    // 1. 缓存检查
    cacheKey := fmt.Sprintf("router:match:%s:%d", req.Query, req.TopK)
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        var resp MatchResponse
        json.Unmarshal([]byte(cached), &resp)
        return &resp, nil
    }

    // 2. 查询向量化
    embedStart := time.Now()
    queryVector, err := s.embeddingClient.Embed(ctx, req.Query)
    meta.EmbeddingTimeMs = int(time.Since(embedStart).Milliseconds())
    if err != nil {
        return nil, err
    }

    // 3. 向量检索（Milvus）
    vectorStart := time.Now()
    vectorResults, err := s.vectorRepo.Search(ctx, queryVector, 50, req.Filters)
    meta.VectorSearchMs = int(time.Since(vectorStart).Milliseconds())
    if err != nil {
        logger.Warn("vector search failed, fallback to keyword only", "error", err)
    }

    // 4. 关键词检索（Meilisearch）
    keywordResults, err := s.searchRepo.Search(ctx, req.Query, &repository.SearchFilter{Limit: 50})
    if err != nil {
        logger.Warn("keyword search failed", "error", err)
    }

    // 5. RRF 融合
    fusedResults := s.rrfFusion(vectorResults, keywordResults, 50)

    // 6. 重排序
    rerankStart := time.Now()
    reranked, err := s.rerankerClient.Rerank(ctx, req.Query, fusedResults, req.TopK)
    meta.RerankTimeMs = int(time.Since(rerankStart).Milliseconds())
    if err != nil {
        // 降级：直接取融合后的前 TopK
        reranked = fusedResults[:min(req.TopK, len(fusedResults))]
    }

    // 7. 组装响应
    matches := make([]SkillMatch, len(reranked))
    for i, r := range reranked {
        matches[i] = SkillMatch{
            SkillID:     r.SkillID,
            Name:        r.Name,
            Description: r.Description,
            Score:       r.Score,
        }
        if req.IncludeContent {
            matches[i].SkillMDURL = r.SkillMDURL
        }
    }

    meta.TotalTimeMs = int(time.Since(start).Milliseconds())
    resp := &MatchResponse{Matches: matches, Meta: meta}

    // 缓存结果（TTL 30分钟）
    if data, _ := json.Marshal(resp); err == nil {
        s.cache.Set(ctx, cacheKey, string(data), 30*time.Minute)
    }

    return resp, nil
}

// RRF 融合算法
func (s *routerService) rrfFusion(vectorResults, keywordResults []SearchResult, topK int) []SearchResult {
    scores := make(map[string]float64)
    const k = 60.0
    
    for rank, item := range vectorResults {
        scores[item.SkillID] += 1.0 / (k + float64(rank))
    }
    for rank, item := range keywordResults {
        scores[item.SkillID] += 1.0 / (k + float64(rank))
    }
    
    // 按融合分数排序
    var items []SearchResult
    for id, score := range scores {
        // 需要从原结果中获取完整信息
        items = append(items, SearchResult{SkillID: id, Score: float32(score)})
    }
    sort.Slice(items, func(i, j int) bool { return items[i].Score > items[j].Score })
    
    if len(items) > topK {
        items = items[:topK]
    }
    return items
}

func (s *routerService) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
    // 1. 先调用 Match 获取最佳技能
    matchReq := &MatchRequest{
        Query:     req.Query,
        TopK:      req.TopK,
        Filters:   nil,
        IncludeContent: true,
    }
    matchResp, err := s.Match(ctx, matchReq)
    if err != nil {
        return nil, err
    }
    if len(matchResp.Matches) == 0 {
        return nil, errors.New("no matching skill found")
    }

    // 2. 获取技能的完整内容
    primarySkill := matchResp.Matches[0]
    skillID, _ := uuid.Parse(primarySkill.SkillID)
    skill, err := s.skillRepo.GetByID(ctx, skillID)
    if err != nil {
        return nil, err
    }

    // 3. 组装提示词
    prompt := s.buildPrompt(skill, req)

    // 4. 调用 LLM
    llmStart := time.Now()
    llmResp, err := s.llmGateway.Chat(ctx, &llm.ChatRequest{
        Model:       req.ExecutionConfig.Model,
        Messages:    []llm.Message{{Role: "user", Content: prompt}},
        MaxTokens:   req.ExecutionConfig.MaxTokens,
        Temperature: req.ExecutionConfig.Temperature,
        Stream:      req.ExecutionConfig.Stream,
    })
    llmTimeMs := int(time.Since(llmStart).Milliseconds())

    if err != nil {
        return nil, err
    }

    // 5. 记录日志
    go s.logExecution(ctx, req, primarySkill.SkillID, llmResp, llmTimeMs, matchResp.Meta.TotalTimeMs)

    return &ExecuteResponse{
        Result: llmResp.Content,
        SelectedSkills: []SelectedSkill{{
            SkillID:      primarySkill.SkillID,
            Name:         primarySkill.Name,
            Contribution: "primary",
            MatchScore:   primarySkill.Score,
        }},
        Meta: ExecutionMeta{
            MatchTimeMs: matchResp.Meta.TotalTimeMs,
            LLMTimeMs:   llmTimeMs,
            TokensUsed:  llmResp.Usage.TotalTokens,
            RequestID:   uuid.New().String(),
        },
    }, nil
}

func (s *routerService) buildPrompt(skill *domain.Skill, req *ExecuteRequest) string {
    // 使用模板构建
    return fmt.Sprintf(`你是一个AI助手，可以使用以下技能来完成任务。请严格遵循技能中的指令。

【技能名称】%s
【技能描述】%s
【技能指令】
%s

【用户请求】
%s

请根据技能指令处理用户请求，并返回结果。`, 
        skill.Name, 
        skill.Description, 
        skill.SkillMDContent,
        req.Query)
}
```



### 2.5 智能路由服务

#### 2.5.1 GitHub 爬虫

```golang
// internal/crawler/github_crawler.go
package crawler

import (
    "context"
    "encoding/base64"
    "net/http"
    "strings"
    "sync"
    "time"
    "github.com/google/go-github/v58/github"
    "golang.org/x/oauth2"
    "skillhub/internal/domain"
)

type GitHubCrawler struct {
    clients      []*github.Client // 多 Token 轮换
    currentIdx   int
    mu           sync.Mutex
    rateLimit    *RateLimiter
}

func NewGitHubCrawler(tokens []string) *GitHubCrawler {
    clients := make([]*github.Client, len(tokens))
    for i, token := range tokens {
        ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
        tc := oauth2.NewClient(context.Background(), ts)
        clients[i] = github.NewClient(tc)
    }
    return &GitHubCrawler{
        clients:   clients,
        rateLimit: NewRateLimiter(5000, time.Hour), // GitHub 限制 5000/小时
    }
}

func (c *GitHubCrawler) getClient() *github.Client {
    c.mu.Lock()
    defer c.mu.Unlock()
    client := c.clients[c.currentIdx]
    c.currentIdx = (c.currentIdx + 1) % len(c.clients)
    return client
}

// 发现技能仓库
func (c *GitHubCrawler) DiscoverRepositories(ctx context.Context, minStars int) ([]*Repository, error) {
    var allRepos []*Repository
    
    // 策略1：Topic 搜索
    topics := []string{"claude-skills", "agent-skills", "ai-skills", "mcp-skills"}
    for _, topic := range topics {
        repos, err := c.searchByTopic(ctx, topic, minStars)
        if err != nil {
            continue
        }
        allRepos = append(allRepos, repos...)
    }
    
    // 策略2：路径搜索
    paths := []string{".claude/skills", "skills", ".github/skills"}
    for _, path := range paths {
        repos, err := c.searchByPath(ctx, path, minStars)
        if err != nil {
            continue
        }
        allRepos = append(allRepos, repos...)
    }
    
    // 去重
    seen := make(map[string]bool)
    uniqueRepos := make([]*Repository, 0)
    for _, r := range allRepos {
        if !seen[r.FullName] {
            seen[r.FullName] = true
            uniqueRepos = append(uniqueRepos, r)
        }
    }
    return uniqueRepos, nil
}

func (c *GitHubCrawler) searchByTopic(ctx context.Context, topic string, minStars int) ([]*Repository, error) {
    client := c.getClient()
    query := fmt.Sprintf("topic:%s stars:>=%d", topic, minStars)
    return c.executeSearch(ctx, client, query)
}

func (c *GitHubCrawler) executeSearch(ctx context.Context, client *github.Client, query string) ([]*Repository, error) {
    opts := &github.SearchOptions{
        Sort:  "stars",
        Order: "desc",
        ListOptions: github.ListOptions{PerPage: 100},
    }
    
    var allRepos []*Repository
    for {
        c.rateLimit.Wait()
        result, resp, err := client.Search.Repositories(ctx, query, opts)
        if err != nil {
            return nil, err
        }
        for _, repo := range result.Repositories {
            allRepos = append(allRepos, &Repository{
                FullName:    repo.GetFullName(),
                CloneURL:    repo.GetCloneURL(),
                Stars:       repo.GetStargazersCount(),
                Description: repo.GetDescription(),
                UpdatedAt:   repo.GetUpdatedAt().Time,
            })
        }
        if resp.NextPage == 0 {
            break
        }
        opts.Page = resp.NextPage
    }
    return allRepos, nil
}

// 获取 SKILL.md 内容
func (c *GitHubCrawler) FetchSkillMD(ctx context.Context, repoFullName string) (string, string, error) {
    client := c.getClient()
    owner, repoName, _ := strings.Cut(repoFullName, "/")
    
    // 尝试多个可能路径
    possiblePaths := []string{
        "SKILL.md",
        "skills/SKILL.md",
        ".claude/skills/SKILL.md",
        ".github/skills/SKILL.md",
    }
    
    for _, path := range possiblePaths {
        c.rateLimit.Wait()
        content, _, resp, err := client.Repositories.GetContents(ctx, owner, repoName, path, nil)
        if err != nil {
            if resp != nil && resp.StatusCode == http.StatusNotFound {
                continue
            }
            return "", "", err
        }
        decoded, _ := base64.StdEncoding.DecodeString(*content.Content)
        return string(decoded), path, nil
    }
    return "", "", fmt.Errorf("SKILL.md not found in %s", repoFullName)
}
```

#### 2.5.2 解析器与同步调度

```golang
// internal/crawler/parser.go
package crawler

import (
    "regexp"
    "strings"
    "gopkg.in/yaml.v3"
    "skillhub/internal/domain"
)

type SkillMDFrontmatter struct {
    Name        string   `yaml:"name"`
    Description string   `yaml:"description"`
    Version     string   `yaml:"version"`
    Author      string   `yaml:"author"`
    Tags        []string `yaml:"tags"`
    Category    string   `yaml:"category"`
    License     string   `yaml:"license"`
}

func ParseSkillMD(content string) (*SkillMDFrontmatter, string, error) {
    // 提取 YAML frontmatter
    re := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)
    matches := re.FindStringSubmatch(content)
    if len(matches) != 3 {
        return nil, content, nil // 无 frontmatter
    }
    
    var fm SkillMDFrontmatter
    if err := yaml.Unmarshal([]byte(matches[1]), &fm); err != nil {
        return nil, "", err
    }
    return &fm, strings.TrimSpace(matches[2]), nil
}

// internal/service/sync_service.go
type SyncService interface {
    SyncOfficialRepo(ctx context.Context) error
    SyncGitHubDiscoveries(ctx context.Context) error
    SyncSingleRepository(ctx context.Context, repoURL string) (*domain.Skill, error)
}

type syncService struct {
    crawler       *crawler.GitHubCrawler
    skillRepo     repository.SkillRepository
    categoryRepo  repository.CategoryRepository
    embeddingWorker *embedding.Worker
    securityScanner *SecurityScanner
    messageQueue  MessageQueue
}

func (s *syncService) SyncSingleRepository(ctx context.Context, repoURL string) (*domain.Skill, error) {
    // 1. 获取 SKILL.md
    content, path, err := s.crawler.FetchSkillMD(ctx, repoURL)
    if err != nil {
        return nil, err
    }
    
    // 2. 解析
    fm, body, err := crawler.ParseSkillMD(content)
    if err != nil {
        return nil, err
    }
    
    // 3. 构建 Skill 实体
    skill := &domain.Skill{
        Name:           fm.Name,
        Description:    fm.Description,
        SourceType:     domain.SourceGitHub,
        SourceURL:      repoURL,
        Author:         fm.Author,
        Version:        fm.Version,
        SkillMDContent: content,
        Metadata: domain.JSONMap{
            "path":  path,
            "tags":  fm.Tags,
            "raw_yaml": fm,
        },
        SecurityStatus: domain.SecurityPending,
    }
    
    // 4. 保存到数据库
    existing, _ := s.skillRepo.GetBySourceURL(ctx, repoURL)
    if existing != nil {
        skill.ID = existing.ID
        skill.CreatedAt = existing.CreatedAt
        if err := s.skillRepo.Update(ctx, skill); err != nil {
            return nil, err
        }
    } else {
        if err := s.skillRepo.Create(ctx, skill); err != nil {
            return nil, err
        }
    }
    
    // 5. 存储文件到 MinIO
    skill.SkillMDURL, _ = s.uploadToMinIO(ctx, skill.ID.String(), content)
    s.skillRepo.Update(ctx, skill)
    
    // 6. 异步任务：安全扫描、向量化
    s.messageQueue.Publish("skill.created", skill.ID)
    
    return skill, nil
}
```



### 2.6 技能同步服务

```golang
// internal/embedding/client.go
package embedding

import "context"

type Client interface {
    Embed(ctx context.Context, text string) ([]float32, error)
    BatchEmbed(ctx context.Context, texts []string) ([][]float32, error)
}

type RerankerClient interface {
    Rerank(ctx context.Context, query string, documents []SearchResult, topK int) ([]SearchResult, error)
}

type SearchResult struct {
    SkillID     string
    Name        string
    Description string
    SkillMDURL  string
    Score       float32
}

// OpenAI 兼容实现
type openAIEmbeddingClient struct {
    endpoint   string
    apiKey     string
    model      string
    httpClient *http.Client
}

func (c *openAIEmbeddingClient) Embed(ctx context.Context, text string) ([]float32, error) {
    reqBody := map[string]interface{}{
        "model": c.model,
        "input": text,
    }
    // 发送 HTTP 请求...
    // 返回向量
}

// internal/embedding/worker.go
type Worker struct {
    embeddingClient Client
    llmClient       llm.Client  // 用于生成增强描述
    skillRepo       repository.SkillRepository
    vectorRepo      repository.VectorRepository
}

func (w *Worker) ProcessSkill(ctx context.Context, skillID uuid.UUID) error {
    skill, err := w.skillRepo.GetByID(ctx, skillID)
    if err != nil {
        return err
    }
    
    // 1. 生成增强描述（如果没有）
    if skill.EnhancedDesc == "" {
        enhanced, err := w.generateEnhancedDesc(ctx, skill)
        if err != nil {
            return err
        }
        skill.EnhancedDesc = enhanced
    }
    
    // 2. 生成向量
    textToEmbed := fmt.Sprintf("%s\n%s\n%s", skill.Name, skill.Description, skill.EnhancedDesc)
    vector, err := w.embeddingClient.Embed(ctx, textToEmbed)
    if err != nil {
        return err
    }
    
    // 3. 存入 Milvus
    err = w.vectorRepo.Upsert(ctx, &VectorRecord{
        SkillID:     skill.ID.String(),
        Embedding:   vector,
        Name:        skill.Name,
        Category:    skill.CategoryID.String(),
        SourceType:  string(skill.SourceType),
    })
    if err != nil {
        return err
    }
    
    // 4. 更新数据库
    skill.EmbeddingVersion = "bge-m3-v1"
    now := time.Now()
    skill.LastEmbeddedAt = &now
    return w.skillRepo.UpdateEmbeddingInfo(ctx, skill.ID, skill.EmbeddingVersion)
}

func (w *Worker) generateEnhancedDesc(ctx context.Context, skill *domain.Skill) (string, error) {
    prompt := fmt.Sprintf(`请阅读以下技能的SKILL.md内容，并提取出该技能的核心能力、主要步骤、适用场景和前置条件。输出一段简洁的描述（不超过200字）。

技能内容：
%s

描述：`, skill.SkillMDContent)
    
    resp, err := w.llmClient.Chat(ctx, &llm.ChatRequest{
        Model: "gpt-4o-mini",
        Messages: []llm.Message{{Role: "user", Content: prompt}},
        MaxTokens: 300,
    })
    return resp.Content, err
}
```



### 2.7 认证中间件

```golang
// internal/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "skillhub/internal/service"
)

func JWTAuth(authService service.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "msg": "missing authorization header"})
            c.Abort()
            return
        }
        
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "msg": "invalid authorization format"})
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        claims, err := authService.ValidateJWT(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "msg": "invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}

func APIKeyAuth(authService service.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "msg": "missing api key"})
            c.Abort()
            return
        }
        
        userID, err := authService.ValidateAPIKey(apiKey)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "msg": "invalid api key"})
            c.Abort()
            return
        }
        
        c.Set("user_id", userID)
        c.Set("api_key", apiKey)
        c.Next()
    }
}

// 签名验证中间件（用于高安全接口）
func HMACSignature(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        timestamp := c.GetHeader("X-Timestamp")
        nonce := c.GetHeader("X-Nonce")
        signature := c.GetHeader("X-Signature")
        
        // 验证时间戳（5分钟内有效）
        // 验证 nonce 是否重复（Redis）
        // 计算签名并对比
        // ...
        c.Next()
    }
}
```






## 3. 数据库详细设计

### 3.1 PostgreSQL 完整表结构

与前述设计相同，XORM 可通过 `Sync2` 自动创建或更新表结构。以下为建表 SQL 参考（与前述一致，此处仅列出核心表）：

```sql
-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 分类表
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 技能表
CREATE TABLE skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    source_type VARCHAR(50),
    source_url VARCHAR(500),
    github_stars INT DEFAULT 0,
    github_forks INT DEFAULT 0,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    author VARCHAR(255),
    version VARCHAR(50),
    download_count INT DEFAULT 0,
    view_count INT DEFAULT 0,
    rating_avg DECIMAL(2,1),
    rating_count INT DEFAULT 0,
    security_status VARCHAR(20) DEFAULT 'pending',
    skill_md_content TEXT,
    skill_md_url VARCHAR(500),
    metadata JSONB,
    enhanced_description TEXT,
    embedding_version VARCHAR(20),
    last_embedded_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_skills_category ON skills(category_id);
CREATE INDEX idx_skills_stars ON skills(github_stars DESC);
CREATE INDEX idx_skills_source ON skills(source_type);
CREATE INDEX idx_skills_security ON skills(security_status);
CREATE INDEX idx_skills_created ON skills(created_at DESC);
```

### 3.2 Milvus Collection Schema

```python
# milvus_schema.py
from pymilvus import Collection, FieldSchema, CollectionSchema, DataType

fields = [
    FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
    FieldSchema(name="skill_id", dtype=DataType.VARCHAR, max_length=36),
    FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1024),
    FieldSchema(name="name", dtype=DataType.VARCHAR, max_length=255),
    FieldSchema(name="category", dtype=DataType.VARCHAR, max_length=36),
    FieldSchema(name="source_type", dtype=DataType.VARCHAR, max_length=20),
]

schema = CollectionSchema(fields, description="Skill embeddings for semantic search")
collection = Collection("skills", schema)

# 创建索引
index_params = {
    "metric_type": "IP",  # Inner Product
    "index_type": "IVF_FLAT",
    "params": {"nlist": 1024}
}
collection.create_index("embedding", index_params)
collection.load()
```



### 3.3 Meilisearch 索引配置

```json
{
  "uid": "skills",
  "primaryKey": "id",
  "searchableAttributes": ["name", "description", "enhanced_description", "tags"],
  "filterableAttributes": ["category", "source_type", "github_stars", "security_status"],
  "sortableAttributes": ["github_stars", "download_count", "created_at", "updated_at"],
  "typoTolerance": {
    "enabled": true,
    "minWordSizeForTypos": {
      "oneTypo": 5,
      "twoTypos": 9
    }
  },
  "pagination": {
    "maxTotalHits": 10000
  }
}
```



### 3.4 Redis 数据结构

| Key 模式                    | 类型       | 用途                 | TTL          |
| :-------------------------- | :--------- | :------------------- | :----------- |
| `skill:detail:{id}`         | String     | 技能详情缓存         | 10分钟       |
| `skills:trending`           | Sorted Set | 热门技能（按下载量） | 无，每日更新 |
| `router:match:{query_hash}` | String     | 路由匹配结果缓存     | 30分钟       |
| `ratelimit:api:{key}`       | String     | API 限流计数器       | 窗口期       |
| `ratelimit:ip:{ip}`         | String     | IP 限流              | 窗口期       |
| `session:{token}`           | String     | 用户会话             | 7天          |
| `nonce:{nonce}`             | String     | 防重放 nonce         | 5分钟        |

## 4. VS Code 插件详细设计

### 4.1 插件激活与生命周期

```typescript
// src/extension.ts
import * as vscode from 'vscode';
import { SkillTreeDataProvider } from './views/skillTreeView';
import { SkillHubService } from './services/skillHubService';
import { AuthService } from './services/authService';

export async function activate(context: vscode.ExtensionContext) {
    console.log('SkillHub Pro extension is now active');
    
    const authService = new AuthService(context);
    const skillHubService = new SkillHubService(authService);
    
    // 初始化树视图
    const treeDataProvider = new SkillTreeDataProvider(skillHubService, authService);
    const treeView = vscode.window.createTreeView('skillhub-sidebar', {
        treeDataProvider,
        showCollapseAll: true
    });
    context.subscriptions.push(treeView);
    
    // 注册命令
    context.subscriptions.push(
        vscode.commands.registerCommand('skillhub.search', () => {
            vscode.commands.executeCommand('skillhub-sidebar.focus');
        }),
        vscode.commands.registerCommand('skillhub.install', (skillId: string) => {
            skillHubService.installSkill(skillId);
        }),
        vscode.commands.registerCommand('skillhub.login', () => {
            authService.login();
        }),
        vscode.commands.registerCommand('skillhub.logout', () => {
            authService.logout();
        }),
        vscode.commands.registerCommand('skillhub.refresh', () => {
            treeDataProvider.refresh();
        }),
        vscode.commands.registerCommand('skillhub.showDetail', (skill: Skill) => {
            showSkillDetailPanel(context, skill, skillHubService);
        }),
        vscode.commands.registerCommand('skillhub.openSkillFolder', (skill: LocalSkill) => {
            vscode.commands.executeCommand('revealFileInOS', vscode.Uri.file(skill.path));
        }),
        vscode.commands.registerCommand('skillhub.toggleSkill', (skill: LocalSkill) => {
            skillHubService.toggleSkill(skill);
        })
    );
    
    // 监听文件变化，提供场景化推荐
    context.subscriptions.push(
        vscode.window.onDidChangeActiveTextEditor(async (editor) => {
            if (editor && authService.isLoggedIn()) {
                const recommendations = await skillHubService.getRecommendations(editor);
                if (recommendations.length > 0) {
                    showRecommendationNotification(recommendations);
                }
            }
        })
    );
}

export function deactivate() {}
```

### 4.2 核心服务类

```typescript
// src/services/skillHubService.ts
import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs/promises';
import { exec } from 'child_process';
import { promisify } from 'util';
import axios, { AxiosInstance } from 'axios';
import { AuthService } from './authService';

const execAsync = promisify(exec);

export interface Skill {
    id: string;
    name: string;
    description: string;
    version: string;
    author: string;
    github_stars: number;
    download_count: number;
    source_url: string;
    skill_md_url?: string;
}

export interface LocalSkill extends Skill {
    path: string;
    enabled: boolean;
    installedVersion: string;
    hasUpdate: boolean;
}

export class SkillHubService {
    private apiClient: AxiosInstance;
    private localSkillsPath: string;
    
    constructor(private authService: AuthService) {
        this.apiClient = axios.create({
            baseURL: 'https://api.skillhub.pro/v1',
            timeout: 10000
        });
        
        // 请求拦截器添加认证头
        this.apiClient.interceptors.request.use(async (config) => {
            const token = await this.authService.getToken();
            if (token) {
                config.headers.Authorization = `Bearer ${token}`;
            }
            return config;
        });
        
        // 本地技能存储目录（可配置）
        this.localSkillsPath = path.join(
            vscode.workspace.workspaceFolders?.[0]?.uri.fsPath || process.cwd(),
            '.claude', 'skills'
        );
    }
    
    async getPopularSkills(): Promise<Skill[]> {
        const response = await this.apiClient.get('/plugin/skills/popular');
        return response.data.data.items;
    }
    
    async searchSkills(query: string): Promise<Skill[]> {
        const response = await this.apiClient.post('/skills/search', { query });
        return response.data.data.items;
    }
    
    async getSkillDetail(id: string): Promise<Skill & { content: string }> {
        const response = await this.apiClient.get(`/skills/${id}`);
        return response.data.data;
    }
    
    async installSkill(skill: Skill): Promise<void> {
        return vscode.window.withProgress({
            location: vscode.ProgressLocation.Notification,
            title: `Installing ${skill.name}...`,
            cancellable: false
        }, async (progress) => {
            try {
                // 获取下载信息
                const resp = await this.apiClient.get(`/plugin/skills/${skill.id}/download`);
                const { download_url, type } = resp.data.data;
                
                const targetDir = path.join(this.localSkillsPath, skill.id);
                await fs.mkdir(targetDir, { recursive: true });
                
                if (type === 'git') {
                    // Git clone
                    await execAsync(`git clone ${download_url} ${targetDir}`);
                } else {
                    // 下载并解压 ZIP
                    const zipPath = path.join(targetDir, 'temp.zip');
                    const response = await axios.get(download_url, { responseType: 'arraybuffer' });
                    await fs.writeFile(zipPath, response.data);
                    // 解压逻辑...
                    await fs.unlink(zipPath);
                }
                
                // 记录本地安装状态
                this.saveLocalSkillMeta(skill, targetDir);
                
                // 同步到云端
                if (this.authService.isLoggedIn()) {
                    await this.apiClient.post('/plugin/user/installed', {
                        skill_id: skill.id,
                        action: 'install',
                        local_version: skill.version
                    });
                }
                
                vscode.window.showInformationMessage(`${skill.name} installed successfully!`);
            } catch (error) {
                vscode.window.showErrorMessage(`Failed to install ${skill.name}: ${error.message}`);
            }
        });
    }
    
    getLocalSkills(): LocalSkill[] {
        // 扫描本地技能目录，读取元数据文件
        // 返回 LocalSkill 数组
        return [];
    }
    
    async toggleSkill(skill: LocalSkill): Promise<void> {
        const newEnabled = !skill.enabled;
        const newPath = skill.enabled 
            ? skill.path + '.disabled' 
            : skill.path.replace(/\.disabled$/, '');
        
        await fs.rename(skill.path, newPath);
        skill.enabled = newEnabled;
        skill.path = newPath;
        this.saveLocalSkillMeta(skill, newPath);
    }
    
    async getRecommendations(editor: vscode.TextEditor): Promise<Skill[]> {
        const document = editor.document;
        const cursorLine = document.lineAt(editor.selection.active.line);
        const context = cursorLine.text;
        
        const response = await this.apiClient.post('/plugin/recommend', {
            file_path: document.fileName,
            file_language: document.languageId,
            cursor_context: context,
            top_k: 3
        });
        return response.data.data.recommendations;
    }
    
    private saveLocalSkillMeta(skill: Skill, localPath: string) {
        const metaPath = path.join(localPath, '.skillhub.json');
        const meta = {
            id: skill.id,
            name: skill.name,
            version: skill.version,
            installedAt: new Date().toISOString()
        };
        fs.writeFile(metaPath, JSON.stringify(meta, null, 2));
    }
}
```

### 4.3 树视图数据提供者

```typescript
// src/views/skillTreeView.ts
import * as vscode from 'vscode';
import { SkillHubService, Skill, LocalSkill } from '../services/skillHubService';
import { AuthService } from '../services/authService';

export class SkillTreeDataProvider implements vscode.TreeDataProvider<TreeItem> {
    private _onDidChangeTreeData = new vscode.EventEmitter<TreeItem | undefined>();
    readonly onDidChangeTreeData = this._onDidChangeTreeData.event;
    
    constructor(
        private service: SkillHubService,
        private authService: AuthService
    ) {}
    
    refresh(): void {
        this._onDidChangeTreeData.fire(undefined);
    }
    
    getTreeItem(element: TreeItem): vscode.TreeItem {
        return element;
    }
    
    async getChildren(element?: TreeItem): Promise<TreeItem[]> {
        if (!element) {
            // 根节点
            return [
                new TreeItem('🔥 Popular', 'popular', vscode.TreeItemCollapsibleState.Expanded),
                new TreeItem('📂 Categories', 'categories', vscode.TreeItemCollapsibleState.Collapsed),
                new TreeItem('📦 My Skills', 'local', vscode.TreeItemCollapsibleState.Expanded),
                new TreeItem('🔍 Search...', 'search', vscode.TreeItemCollapsibleState.None, {
                    command: 'skillhub.search',
                    title: 'Search Skills'
                })
            ];
        }
        
        switch (element.contextValue) {
            case 'popular':
                const skills = await this.service.getPopularSkills();
                return skills.map(s => new SkillItem(s));
            case 'categories':
                return [
                    new TreeItem('Document Processing', 'category'),
                    new TreeItem('Code Development', 'category'),
                    // ... 更多分类
                ];
            case 'local':
                if (!this.authService.isLoggedIn()) {
                    return [new TreeItem('Sign in to sync your skills', 'login-prompt', vscode.TreeItemCollapsibleState.None, {
                        command: 'skillhub.login',
                        title: 'Sign In'
                    })];
                }
                const localSkills = this.service.getLocalSkills();
                return localSkills.map(s => new LocalSkillItem(s));
            default:
                return [];
        }
    }
}

class TreeItem extends vscode.TreeItem {
    constructor(
        public readonly label: string,
        public readonly contextValue: string,
        public readonly collapsibleState: vscode.TreeItemCollapsibleState = vscode.TreeItemCollapsibleState.None,
        public readonly command?: vscode.Command
    ) {
        super(label, collapsibleState);
    }
}

class SkillItem extends TreeItem {
    constructor(skill: Skill) {
        super(skill.name, 'skill', vscode.TreeItemCollapsibleState.None);
        this.description = `${skill.author} ⭐ ${skill.github_stars}`;
        this.tooltip = skill.description;
        this.iconPath = new vscode.ThemeIcon('package');
        this.command = {
            command: 'skillhub.showDetail',
            title: 'Show Skill Detail',
            arguments: [skill]
        };
        this.contextValue = 'skill';
    }
}

class LocalSkillItem extends TreeItem {
    constructor(skill: LocalSkill) {
        super(skill.name, 'localSkill', vscode.TreeItemCollapsibleState.None);
        this.description = skill.enabled ? 'Enabled' : 'Disabled';
        this.iconPath = new vscode.ThemeIcon(skill.enabled ? 'check' : 'circle-slash');
        this.contextValue = 'localSkill';
    }
}
```

## 5. API 接口详细定义

（此处以表格形式列出完整接口，部分已在 PRD 细化中列出，现补充详细请求/响应格式）

### 5.1 技能管理

#### `GET /api/v1/skills`

**响应**：

```json
{
  "code": 0,
  "data": {
    "items": [ { ... } ],
    "pagination": { "page": 1, "size": 20, "total": 10234 }
  }
}
```



#### `GET /api/v1/skills/{id}`

**响应**：

```json
{
  "code": 0,
  "data": {
    "id": "uuid",
    "name": "excel-trend-analyzer",
    "description": "...",
    "skill_md_content": "---\nname: ...\n---\n\n# Instructions...",
    "skill_md_url": "https://oss...",
    "author": "anthropic",
    "version": "1.2.0",
    "github_stars": 12500,
    "download_count": 8430,
    "rating_avg": 4.5,
    "rating_count": 128,
    "category": { "id": "...", "name": "Data Analysis", "slug": "data-analysis" },
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-04-20T08:30:00Z"
  }
}
```



### 5.2 智能路由

#### `POST /api/v1/router/match`

**请求**：

```json
{
  "query": "分析销售Excel趋势",
  "top_k": 3,
  "include_content": false,
  "filters": { "source_type": ["official", "github"], "min_stars": 10 }
}
```



**响应**：

```json
{
  "code": 0,
  "data": {
    "matches": [
      {
        "skill_id": "uuid",
        "name": "excel-trend-analyzer",
        "description": "分析Excel数据趋势并生成图表",
        "match_score": 0.952,
        "skill_md_url": "https://..."
      }
    ],
    "meta": {
      "embedding_time_ms": 120,
      "vector_search_time_ms": 45,
      "rerank_time_ms": 38,
      "total_time_ms": 203
    }
  }
}
```



#### `POST /api/v1/router/execute`

**请求**：

```json
{
  "query": "分析附件销售数据并生成周报",
  "top_k": 1,
  "execution_config": {
    "model": "claude-3.5-sonnet",
    "max_tokens": 4096,
    "temperature": 0.2,
    "stream": false
  },
  "context": { "user_id": "user123" },
  "files": [ { "name": "sales.xlsx", "url": "https://..." } ]
}
```



**响应**：

```json
{
  "code": 0,
  "data": {
    "result": "根据分析，3月销售额为...",
    "selected_skills": [{
      "skill_id": "uuid",
      "name": "excel-trend-analyzer",
      "contribution": "primary",
      "match_score": 0.952
    }],
    "execution_meta": {
      "match_time_ms": 203,
      "llm_time_ms": 2340,
      "total_tokens_used": 1850,
      "request_id": "req_abc123"
    }
  }
}
```



### 5.3 插件专用接口

#### `GET /api/v1/plugin/skills/popular`

**响应**：

```json
{
  "code": 0,
  "data": {
    "items": [ { ... } ],
    "updated_at": "2026-04-23T10:00:00Z"
  }
}
```



#### `GET /api/v1/plugin/skills/{id}/download`

**响应**：

```json
{
  "code": 0,
  "data": {
    "download_url": "https://github.com/.../repo.git",
    "type": "git",
    "version": "1.2.0"
  }
}
```



#### `POST /api/v1/plugin/user/installed`

**请求**：

```json
{
  "skill_id": "uuid",
  "action": "install",  // install, uninstall, update
  "local_version": "1.2.0"
}
```



#### `POST /api/v1/plugin/recommend`

**请求**：

```json
{
  "file_path": "src/analysis/sales.py",
  "file_language": "python",
  "cursor_context": "df = pd.read_excel('sales.xlsx')",
  "top_k": 5
}
```


## 6. 部署配置示例

### 6.1 Docker Compose（开发环境，使用 Consul 替代 ETCD）

以下配置中，Milvus 的元数据存储由默认的 ETCD 改为 Consul。

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: skillhub
      POSTGRES_USER: skillhub
      POSTGRES_PASSWORD: devpass
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  meilisearch:
    image: getmeili/meilisearch:v1.5
    ports:
      - "7700:7700"
    environment:
      MEILI_MASTER_KEY: devMasterKey
    volumes:
      - meili_data:/meili_data

  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio_data:/data

  # Consul 作为服务发现和 KV 存储（同时供 Milvus 元数据存储）
  consul:
    image: consul:1.15
    ports:
      - "8500:8500"
    command: agent -server -bootstrap -ui -client=0.0.0.0
    volumes:
      - consul_data:/consul/data

  milvus:
    image: milvusdb/milvus:v2.3.4
    ports:
      - "19530:19530"
      - "9091:9091"
    environment:
      # 元数据存储使用 Consul
      METASTORE_TYPE: consul
      CONSUL_ADDRESS: consul:8500
      # 对象存储使用 MinIO
      MINIO_ADDRESS: minio:9000
      MINIO_ACCESS_KEY_ID: minioadmin
      MINIO_SECRET_ACCESS_KEY: minioadmin
    depends_on:
      - consul
      - minio

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      REDIS_ADDR: redis:6379
      MEILISEARCH_HOST: http://meilisearch:7700
      MILVUS_HOST: milvus
      MILVUS_PORT: 19530
      MINIO_ENDPOINT: minio:9000
      CONSUL_ADDR: consul:8500   # 用于服务注册
    depends_on:
      - postgres
      - redis
      - meilisearch
      - milvus
      - minio
      - consul

  web:
    build: ./web
    ports:
      - "3000:80"
    depends_on:
      - backend

volumes:
  postgres_data:
  meili_data:
  minio_data:
  consul_data:
```

### 6.2 Milvus 使用 Consul 的配置说明

- Milvus 2.3 及以上版本支持通过环境变量 `METASTORE_TYPE=consul` 和 `CONSUL_ADDRESS` 指定使用 Consul 作为元数据存储。
- 若使用 Helm 部署，可在 `values.yaml` 中设置：
  ```yaml
  externalEtcd:
    enabled: false
  externalConsul:
    enabled: true
    host: consul.consul.svc.cluster.local
    port: 8500
  ```

### 6.3 Kubernetes Deployment 示例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: skillhub-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: skillhub-backend
  template:
    metadata:
      labels:
        app: skillhub-backend
    spec:
      containers:
      - name: backend
        image: skillhub/backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: skillhub-secrets
              key: db-host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: skillhub-secrets
              key: db-password
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: skillhub-backend
spec:
  selector:
    app: skillhub-backend
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```




## 7. 异常处理与日志规范

### 7.1 错误码定义

| 错误码 | HTTP状态码 | 说明                         |
| :----- | :--------- | :--------------------------- |
| 0      | 200        | 成功                         |
| 1001   | 400        | 参数错误                     |
| 1002   | 401        | 认证失败                     |
| 1003   | 403        | 权限不足                     |
| 1004   | 404        | 资源不存在                   |
| 1005   | 409        | 资源冲突                     |
| 2001   | 500        | 数据库错误                   |
| 2002   | 500        | 外部服务错误（GitHub API等） |
| 2003   | 503        | 服务暂时不可用               |
| 3001   | 400        | 路由匹配失败（无相关技能）   |
| 3002   | 500        | LLM调用失败                  |

### 7.2 日志格式

采用结构化日志（JSON格式），包含以下字段：

- `level`: debug, info, warn, error
- `ts`: 时间戳
- `caller`: 调用位置
- `msg`: 消息
- `trace_id`: 请求追踪ID
- `user_id`: 用户ID（如有）
- `duration_ms`: 耗时
- 其他自定义字段

**示例**：

```json
{"level":"info","ts":"2026-04-23T10:15:30Z","caller":"router/service.go:123","msg":"route matched","trace_id":"abc123","user_id":"uuid","query":"分析Excel","matched_count":3,"duration_ms":203}
```


## 8. 测试策略



### 8.1 单元测试

- 使用 `testify` 框架。

- 对 XORM 仓储层使用 SQLite 内存模式进行测试，避免依赖真实 PostgreSQL。

  ```go
  engine, _ := xorm.NewEngine("sqlite3", ":memory:")
  ```

- 对每个 service 层方法编写测试，mock 依赖的 repository 和外部 client。

- 覆盖率目标：核心业务逻辑 > 80%。

### 8.2 集成测试

- 使用 `testcontainers-go` 启动真实依赖（PostgreSQL, Redis, Milvus）。
- 测试完整的 API 调用链路。

### 8.3 路由准确率评测

- 建立包含 500 条人工标注的查询-技能对的数据集。
- 定期运行评测脚本，计算 Top-1, Top-3 准确率。
- 结果纳入 CI 报告。


## 9. 附录

### 9.1 配置项完整列表

| 配置项                  | 类型     | 默认值                                          | 说明                                       |
| :---------------------- | :------- | :---------------------------------------------- | :----------------------------------------- |
| `SERVER_PORT`           | int      | 8080                                            | 服务端口                                   |
| `DB_HOST`               | string   | localhost                                       | PostgreSQL主机                             |
| `DB_PORT`               | int      | 5432                                            | PostgreSQL端口                             |
| `DB_NAME`               | string   | skillhub                                        | 数据库名                                   |
| `DB_USER`               | string   | skillhub                                        | 用户名                                     |
| `DB_PASSWORD`           | string   |                                                 | 密码                                       |
| `REDIS_ADDR`            | string   | localhost:6379                                  | Redis地址                                  |
| `MEILISEARCH_HOST`      | string   | [http://localhost:7700](http://localhost:7700/) | Meilisearch地址                            |
| `MEILISEARCH_KEY`       | string   |                                                 | Master Key                                 |
| `MILVUS_HOST`           | string   | localhost                                       | Milvus主机                                 |
| `MILVUS_PORT`           | int      | 19530                                           | Milvus端口                                 |
| `MINIO_ENDPOINT`        | string   | localhost:9000                                  | MinIO端点                                  |
| `MINIO_ACCESS_KEY`      | string   |                                                 | 访问密钥                                   |
| `MINIO_SECRET_KEY`      | string   |                                                 | 私有密钥                                   |
| `GITHUB_TOKENS`         | []string |                                                 | GitHub Token 列表                          |
| `EMBEDDING_MODEL`       | string   | BAAI/bge-m3                                     | 向量模型                                   |
| `EMBEDDING_ENDPOINT`    | string   |                                                 | 向量服务地址                               |
| `RERANKER_MODEL`        | string   | BAAI/bge-reranker-v2-m3                         | 重排序模型                                 |
| `LLM_DEFAULT_MODEL`     | string   | claude-3.5-sonnet                               | 默认LLM                                    |
| `ANTHROPIC_API_KEY`     | string   |                                                 | Anthropic API Key                          |
| `OPENAI_API_KEY`        | string   |                                                 | OpenAI API Key                             |
| `JWT_SECRET`            | string   |                                                 | JWT签名密钥                                |
| `API_RATE_LIMIT_FREE`   | int      | 100                                             | 免费版日限额                               |
| `CONSUL_ADDR`           | string   | localhost:8500                                  | Consul 地址，用于服务注册与配置中心        |
| `MILVUS_METASTORE_TYPE` | string   | consul                                          | Milvus 元数据存储类型（etcd/consul/mysql） |
| `MILVUS_CONSUL_ADDR`    | string   | consul:8500                                     | Milvus 使用的 Consul 地址                  |



### 9.2 ORM 切换注意事项

1. **Tag 差异**：XORM 使用 `pk` 标识主键，`default` 设置默认值，`notnull` 设置非空，而 GORM 使用 `primaryKey`、`default`、`not null`。
2. **软删除**：XORM 通过 `deleted` tag 支持软删除，删除时自动填充时间。
3. **关联查询**：XORM 不默认加载关联，需要显式使用 `engine.Join()` 或手动查询。
4. **钩子方法**：XORM 支持 `BeforeInsert()`、`AfterInsert()` 等钩子，用法与 GORM 类似。

### 9.3 Consul 替代 ETCD 的收益

- 减少运维组件数量，降低系统复杂度。
- 利用 Consul 的服务健康检查能力，增强微服务可观测性。
- 对于 Milvus 而言，元数据存储 QPS 要求不高，Consul 完全胜任。

---

以上为基于 XORM 和 Consul 调整后的 SkillHub Pro 详细设计说明书。开发团队可依据此文档进行编码实现。