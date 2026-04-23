# SkillHub Pro（技能宝库）

> 一个面向开发者的 AI 技能搜索与安装平台，聚合 GitHub 开源 AI 技能（cursorrules、clinerules、agent skills 等），支持语义搜索、一键安装和 VS Code 插件管理。

## 项目结构

```
skill-hub/
├── backend/                 # Go 后端服务
│   ├── cmd/                 # 入口
│   │   ├── router-api/      # 智能路由 API 服务
│   │   ├── sync-worker/     # 技能同步 Worker
│   │   └── admin-api/       # 管理后台 API
│   ├── internal/            # 内部逻辑
│   ├── migrations/          # 数据库迁移
│   └── deployments/         # 部署配置
├── frontend/                # Vue3 + Vite 前端
├── docs/                    # 项目文档
└── deployments/             # K8s 部署清单
```

## 快速开始

### 后端

```bash
cd backend
make dev-router
```

### 前端

```bash
cd frontend
npm install
npm run dev
```

## 技术栈

- **后端**: Go + Gin + XORM
- **前端**: Vue3 + Vite + TypeScript
- **数据库**: PostgreSQL + Redis + Meilisearch + Milvus
- **基础设施**: Docker Compose / Kubernetes
