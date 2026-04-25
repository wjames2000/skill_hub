# SkillHub Pro

<p align="center">
  <img src="https://via.placeholder.com/150?text=SkillHub" alt="SkillHub Pro Logo" width="150"/>
</p>

<p align="center">
  <strong>一站式 AI 技能聚合与智能路由平台</strong>
</p>

<p align="center">
  <a href="#功能特性">功能特性</a> •
  <a href="#技术栈">技术栈</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#系统架构">系统架构</a> •
  <a href="#api-概览">API 概览</a> •
  <a href="#vs-code-插件">VS Code 插件</a> •
  <a href="#参与贡献">参与贡献</a> •
  <a href="#开源协议">开源协议</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version"/>
  <img src="https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat&logo=vue.js" alt="Vue Version"/>
  <img src="https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql" alt="PostgreSQL"/>
  <img src="https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis" alt="Redis"/>
  <img src="https://img.shields.io/badge/Milvus-2.3+-00A3E0?style=flat&logo=milvus" alt="Milvus"/>
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"/>
</p>

[English](https://github.com/wjames2000/skill_hub/blob/main/README.md) | 中文

---

## 📖 项目简介

**SkillHub Pro** 是一个集发现、管理和智能执行 AI 技能（Skills）于一体的综合平台。它能够自动汇聚官方技能库（如 Anthropic 的 Claude Skills）及 GitHub 高星项目中的技能资源，提供功能强大的 Web 界面供浏览，并推出 VS Code 插件实现无缝 IDE 集成。

SkillHub Pro 的核心亮点是 **智能技能路由器（Intelligent Skill Router）**——一个可理解自然语言查询并精准匹配最合适技能，进而通过 LLM 端到端执行任务的语义引擎。

> *“让 AI 能力的发现与使用，像安装 IDE 插件一样简单。”*

## ✨ 功能特性

### 🔍 多源技能聚合
- **官方源**：自动同步 `anthropics/skills` 及 LUNARTECH Superpowers Library。
- **GitHub 发现**：综合运用 GitHub Topics、路径搜索、Awesome 列表等多种策略爬取仓库，支持配置 Star 数阈值。
- **质量把控**：自动安全扫描（Semgrep）与元数据提取。

### 🧠 智能技能路由器
- **语义匹配**：融合稠密向量检索与关键词匹配，结合 Cross-Encoder 重排序，确保高准确率。
- **端到端执行**：`/router/execute` API 可自动匹配技能并调用 LLM 完成任务。
- **上下文组装**：动态构建提示词，支持技能指令压缩，降低 Token 消耗。

### 🌐 Web 平台（Vue3）
- 按分类、热门度、最新收录浏览技能。
- 支持全文搜索与语义搜索。
- 技能详情页带有语法高亮的 `SKILL.md` 预览。
- 用户账户系统：收藏、评分、评论、API Key 管理。

### 🔌 VS Code 插件
- 侧边栏浏览与安装技能，无需离开编辑器。
- 一键安装至本地 `.claude/skills/`、`.cursor/skills/` 等目录。
- 本地技能管理（启用/禁用、更新检测）。
- 已安装技能的云端同步，跨设备使用。
- 基于当前文件/语言上下文的情景化推荐。

### 🚀 RESTful API
- 提供技能列表、搜索、详情等全面接口。
- API Key 认证与速率限制。
- 专为插件与路由执行设计的端点。

## 🛠️ 技术栈

| 层级         | 技术选型                                     |
| ------------ | -------------------------------------------- |
| **后端**     | Go 1.21+、Gin Web 框架、XORM                 |
| **前端**     | Vue 3、TypeScript、Vite、Pinia、Element Plus |
| **数据库**   | PostgreSQL 15+、Redis 7+                     |
| **检索引擎** | Meilisearch（全文搜索）、Milvus（向量检索）  |
| **消息队列** | RabbitMQ / Asynq                             |
| **对象存储** | MinIO / 阿里云 OSS                           |
| **LLM 网关** | LiteLLM（支持 Claude、GPT、DeepSeek 等）     |
| **插件开发** | VS Code Extension API、TypeScript            |
| **部署运维** | Docker、Kubernetes、GitHub Actions           |

## 🚀 快速开始

### 环境要求
- Docker 及 Docker Compose
- Go 1.21+（本地开发）
- Node.js 18+ 及 pnpm

### 使用 Docker Compose 启动
```bash
# 克隆仓库
git clone https://github.com/your-org/skillhub-pro.git
cd skillhub-pro

# 启动所有服务（PostgreSQL、Redis、Meilisearch、Milvus、MinIO、后端、前端）
docker-compose -f docker-compose.dev.yml up -d

# 访问 Web 界面
open http://localhost:3000
```

### 本地开发环境

**后端**
```bash
cd backend
cp config/config.example.yaml config/config.yaml
# 编辑 config.yaml 填入本地配置
go mod download
go run cmd/server/main.go
```

**前端**
```bash
cd web
pnpm install
pnpm dev
```

**VS Code 插件**
```bash
cd vscode-extension
npm install
npm run watch
# 按 F5 启动扩展开发宿主窗口开始调试
```

## 🏛️ 系统架构

![架构图](https://via.placeholder.com/800x400?text=架构图)

系统采用清晰的分层架构：

- **接入层**：Nginx/Kong 负责路由转发、SSL 终结、限流熔断。
- **业务服务层**：技能管理、同步爬虫、智能路由、用户认证、插件服务等。
- **基础设施层**：PostgreSQL（主数据）、Redis（缓存与会话）、Meilisearch（关键词检索）、Milvus（向量存储）、MinIO（文件存储）。

智能路由器分为三个阶段：
1. **粗排召回**：融合 Milvus 向量检索与 Meilisearch 关键词检索。
2. **精排重排序**：Cross-Encoder 模型对候选集重新打分。
3. **执行**：上下文组装与 LLM 调用。

## 📡 API 概览

| 端点                            | 方法 | 说明                           |
| ------------------------------- | ---- | ------------------------------ |
| `/api/v1/skills`                | GET  | 获取技能列表（支持过滤、分页） |
| `/api/v1/skills/{id}`           | GET  | 获取技能详情                   |
| `/api/v1/skills/search`         | POST | 全文/语义搜索                  |
| `/api/v1/router/match`          | POST | 根据自然语言查询匹配最佳技能   |
| `/api/v1/router/execute`        | POST | 匹配并执行（需配置 LLM）       |
| `/api/v1/plugin/skills/popular` | GET  | 获取热门技能（供插件使用）     |
| `/api/v1/plugin/recommend`      | POST | 获取上下文情景推荐             |

后端运行时，完整的 Swagger 文档可通过 `/swagger/index.html` 访问。

## 🧩 VS Code 插件

插件可在 [VS Code Marketplace](https://marketplace.visualstudio.com/)（链接即将上线）获取。

**主要命令：**
- `SkillHub: Search Skills` —— 打开技能浏览器。
- `SkillHub: Install Skill` —— 安装当前查看的技能。
- `SkillHub: Refresh` —— 同步本地技能与云端。

登录 SkillHub Pro 账户后，插件会自动同步已安装的技能状态。

## 📦 项目结构

```
skillhub-pro/
├── backend/                # Go 后端（Gin + XORM）
│   ├── cmd/server/
│   ├── internal/
│   │   ├── domain/         # 领域模型
│   │   ├── repository/     # 数据访问层（XORM）
│   │   ├── service/        # 业务逻辑
│   │   ├── handler/        # HTTP 控制器
│   │   ├── crawler/        # GitHub 技能爬虫
│   │   ├── embedding/      # 向量嵌入与 Milvus 客户端
│   │   └── llm/            # LLM 网关
│   └── migrations/         # 数据库迁移文件
├── web/                    # Vue3 前端
│   ├── src/
│   │   ├── views/
│   │   ├── components/
│   │   └── stores/
│   └── public/
├── vscode-extension/       # VS Code 插件
│   ├── src/
│   │   ├── views/
│   │   └── services/
│   └── package.json
├── docker-compose.yml
├── docker-compose.dev.yml
└── README.md
```

## 🤝 参与贡献

我们非常欢迎社区的贡献！请参阅[贡献指南](CONTRIBUTING.md)了解如何：

- 通过 Issue 报告 Bug 或提出功能建议
- 提交 Pull Request（Fork 仓库、新建分支、编写测试、更新文档）
- 遵循[行为准则](CODE_OF_CONDUCT.md)

### 开发规范
- 后端：遵循 Go 项目标准布局，使用 `golangci-lint` 进行代码检查。
- 前端：使用 Composition API 和 TypeScript，确保响应式设计。
- 插件：遵守 VS Code UX 指南，同时支持浅色/深色主题。

## 📄 开源协议

本项目基于 Apache License 2.0 协议开源，详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [Anthropic Skills](https://github.com/anthropics/skills) —— 定义了技能标准。
- [Awesome Claude Skills](https://github.com/punkpeye/awesome-claude-skills) —— 社区精选清单。
- [Milvus](https://milvus.io/) —— 向量数据库支持。
- [Meilisearch](https://www.meilisearch.com/) —— 闪电般的全文搜索引擎。

---

<p align="center">
  Made with ❤️ by SkillHub Team
</p>