# SkillHub Pro

<p align="center">
  <img src="https://via.placeholder.com/150?text=SkillHub" alt="SkillHub Pro Logo" width="150"/>
</p>

<p align="center">
  <strong>One-stop AI Skills Aggregation & Intelligent Routing Platform</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#tech-stack">Tech Stack</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#api-overview">API</a> •
  <a href="#vs-code-extension">VS Code Extension</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version"/>
  <img src="https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat&logo=vue.js" alt="Vue Version"/>
  <img src="https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql" alt="PostgreSQL"/>
  <img src="https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis" alt="Redis"/>
  <img src="https://img.shields.io/badge/Milvus-2.3+-00A3E0?style=flat&logo=milvus" alt="Milvus"/>
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"/>
</p>

English | [中文](https://github.com/wjames2000/skill_hub/blob/main/README_CN.md)

---

## 📖 Overview

**SkillHub Pro** is an all-in-one platform for discovering, managing, and intelligently executing AI skills. It automatically aggregates skills from official repositories (like Anthropic's Claude skills) and high-star GitHub projects, provides a powerful web interface for browsing, and offers a VS Code extension for seamless IDE integration.

What sets SkillHub Pro apart is its **Intelligent Skill Router** – a semantic engine that understands natural language queries and matches them with the most appropriate skill, then optionally executes the task end-to-end via LLM integration.

> *"Discover AI capabilities as easily as installing IDE extensions."*

## ✨ Features

### 🔍 Multi-Source Skill Aggregation
- **Official Sources**: Automatically sync from `anthropics/skills` and the LUNARTECH Superpowers library.
- **GitHub Discovery**: Crawl repositories using multiple strategies (GitHub Topics, path searches, Awesome lists) with configurable star thresholds.
- **Quality Control**: Automated security scanning (Semgrep) and metadata extraction.

### 🧠 Intelligent Skill Router
- **Semantic Matching**: Hybrid retrieval (dense vectors + keywords) with cross-encoder reranking for high accuracy.
- **End-to-End Execution**: `/router/execute` API matches a skill and invokes an LLM to complete the task.
- **Context Assembly**: Dynamically builds prompts with compressed skill instructions to minimize token usage.

### 🌐 Web Platform (Vue3)
- Browse skills by category, popularity, or recent updates.
- Full-text and semantic search.
- Detailed skill pages with syntax-highlighted `SKILL.md` preview.
- User accounts, favorites, ratings, and API key management.

### 🔌 VS Code Extension
- Sidebar view for discovering and installing skills without leaving the editor.
- One-click installation to local `.claude/skills/`, `.cursor/skills/`, etc.
- Local skill management (enable/disable, update checks).
- Cloud sync for installed skills across devices.
- Context-aware recommendations based on open file/language.

### 🚀 RESTful API
- Comprehensive API for skill listing, search, and details.
- API key authentication with rate limiting.
- Dedicated endpoints for plugin integration and router execution.

## 🛠️ Tech Stack

| Layer               | Technology                                     |
| ------------------- | ---------------------------------------------- |
| **Backend**         | Go 1.21+, Gin Web Framework, XORM              |
| **Frontend**        | Vue 3, TypeScript, Vite, Pinia, Element Plus   |
| **Database**        | PostgreSQL 15+, Redis 7+                       |
| **Search & Vector** | Meilisearch (full-text), Milvus (vector DB)    |
| **Message Queue**   | RabbitMQ / Asynq                               |
| **Object Storage**  | MinIO / AWS S3                                 |
| **LLM Gateway**     | LiteLLM (supports Claude, GPT, DeepSeek, etc.) |
| **Extension**       | VS Code Extension API, TypeScript              |
| **Deployment**      | Docker, Kubernetes, GitHub Actions             |

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ / pnpm

### Run with Docker Compose
```bash
# Clone the repository
git clone https://github.com/your-org/skillhub-pro.git
cd skillhub-pro

# Start all services (PostgreSQL, Redis, Meilisearch, Milvus, MinIO, Backend, Frontend)
docker-compose -f docker-compose.dev.yml up -d

# Visit the web UI
open http://localhost:3000
```

### Local Development

**Backend**
```bash
cd backend
cp config/config.example.yaml config/config.yaml
# Edit config.yaml with your local credentials
go mod download
go run cmd/server/main.go
```

**Frontend**
```bash
cd web
pnpm install
pnpm dev
```

**VS Code Extension**
```bash
cd vscode-extension
npm install
npm run watch
# Press F5 to start debugging in a new Extension Development Host window
```

## 🏛️ Architecture

![Architecture Diagram](https://via.placeholder.com/800x400?text=Architecture+Diagram)

The system follows a clean, layered architecture:

- **Gateway Layer**: Nginx/Kong for routing, SSL termination, and rate limiting.
- **Business Services**: Skill management, sync crawler, intelligent router, authentication, plugin service.
- **Infrastructure Layer**: PostgreSQL (primary data), Redis (cache & sessions), Meilisearch (keyword search), Milvus (vector embeddings), MinIO (file storage).

The Intelligent Router operates in three stages:
1. **Coarse Retrieval**: Hybrid search combining Milvus ANN and Meilisearch keyword matching.
2. **Fine Reranking**: Cross-encoder model to reorder candidates.
3. **Execution**: Context assembly and LLM invocation.

## 📡 API Overview

| Endpoint                        | Method | Description                                    |
| ------------------------------- | ------ | ---------------------------------------------- |
| `/api/v1/skills`                | GET    | List skills with filtering & pagination        |
| `/api/v1/skills/{id}`           | GET    | Get skill detail                               |
| `/api/v1/skills/search`         | POST   | Full-text/semantic search                      |
| `/api/v1/router/match`          | POST   | Match best skills for a natural language query |
| `/api/v1/router/execute`        | POST   | Match and execute (requires LLM config)        |
| `/api/v1/plugin/skills/popular` | GET    | Get trending skills (for extension)            |
| `/api/v1/plugin/recommend`      | POST   | Get context-aware recommendations              |

Full Swagger documentation is available at `/swagger/index.html` when the backend is running.

## 🧩 VS Code Extension

The extension is available on the [VS Code Marketplace](https://marketplace.visualstudio.com/) (link coming soon).

**Key Commands:**
- `SkillHub: Search Skills` - Open the skill browser.
- `SkillHub: Install Skill` - Install the currently viewed skill.
- `SkillHub: Refresh` - Sync local skills with cloud.

The extension automatically syncs installed skills with your SkillHub Pro account when signed in.

## 📦 Project Structure

```
skillhub-pro/
├── backend/                # Go backend (Gin + XORM)
│   ├── cmd/server/
│   ├── internal/
│   │   ├── domain/         # Domain models
│   │   ├── repository/     # Data access layer (XORM)
│   │   ├── service/        # Business logic
│   │   ├── handler/        # HTTP controllers
│   │   ├── crawler/        # GitHub skill crawler
│   │   ├── embedding/      # Vector embedding & Milvus client
│   │   └── llm/            # LLM gateway
│   └── migrations/         # Database migrations
├── web/                    # Vue3 frontend
│   ├── src/
│   │   ├── views/
│   │   ├── components/
│   │   └── stores/
│   └── public/
├── vscode-extension/       # VS Code extension
│   ├── src/
│   │   ├── views/
│   │   └── services/
│   └── package.json
├── docker-compose.yml
├── docker-compose.dev.yml
└── README.md
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to:

- Report bugs or suggest features via issues.
- Submit pull requests (fork & branch, write tests, update docs).
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md).

### Development Guidelines
- Backend: Follow standard Go project layout. Use `golangci-lint` for linting.
- Frontend: Use Composition API and TypeScript. Ensure responsive design.
- Extension: Adhere to VS Code UX guidelines and support both light/dark themes.

## 📄 License

This project is licensed under the Apache License 2.0 – see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [Anthropic Skills](https://github.com/anthropics/skills) for defining the skill standard.
- [Awesome Claude Skills](https://github.com/punkpeye/awesome-claude-skills) for community curation.
- [Milvus](https://milvus.io/) for the vector database.
- [Meilisearch](https://www.meilisearch.com/) for lightning-fast full-text search.

---

<p align="center">
  Made with ❤️ by the SkillHub Team
</p>