# SkillHub Pro 测试用例文档

## 文档信息

| 项目 | 内容 |
|------|------|
| 项目名称 | SkillHub Pro（技能宝库） |
| 文档版本 | V1.0 |
| 编制日期 | 2026-04-24 |
| 参考文档 | PRD V3.0、API Specification V1.0、WBS V1.0 |
| 测试策略 | 功能测试 + 集成测试 + 性能测试 + 安全测试 + 验收测试 |

---

## 1. 测试范围与策略

### 1.1 测试层次

| 层次 | 范围 | 工具/框架 | 责任人 |
|------|------|-----------|--------|
| 单元测试 | 后端 handler/service/model/pkg | Go testing + testify | 后端 |
| 单元测试 | 前端 component/page/store | Vitest + Testing Library | 前端 |
| 集成测试 | API 完整链路（含数据库/外部服务 Mock） | go test + httptest | 测试 |
| E2E 测试 | Web 端核心流程 | Playwright/Cypress | 测试 |
| 性能测试 | API 并发、爬虫任务 | k6/wrk | 测试 |
| 安全测试 | SQL注入/XSS/认证 | 手工 + 自动化扫描 | 测试 |

### 1.2 PRD 验收标准映射

| PRD 编号 | 需求 | 验收标准 | 测试类型 | 优先级 |
|----------|------|----------|----------|--------|
| 4.1 | 技能同步 | 每日增量 + 每周全量同步成功 | 集成 | P0 |
| 4.2 | 技能浏览与检索 | 列表/详情/搜索功能正常 | 功能 | P0 |
| 4.3 | RESTful API | 所有端点按规范响应 | 集成 | P0 |
| 4.4 | VS Code 插件 | 发现/安装/管理/同步 | E2E | P1 |
| 4.5 | 智能路由器 | Top-1 准确率 >= 85% | 评测 | P0 |
| 5.1 | 性能 | 路由匹配 P95 <= 500ms | 性能 | P0 |
| 5.2 | 安全 | SQL注入/XSS/认证防护 | 安全 | P1 |

---

## 2. 功能测试用例

### 2.1 技能浏览模块

| ID | 用例名称 | 前置条件 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|----------|--------|
| TC-FUNC-001 | 技能列表分页 | 数据库有 >= 20 条技能 | GET /api/v1/skills?page=1&size=20 | 返回 20 条记录，含分页信息 | P0 |
| TC-FUNC-002 | 技能列表按分类筛选 | 有分类 data-analysis | GET /api/v1/skills?category=data-analysis | 仅返回该分类技能 | P0 |
| TC-FUNC-003 | 技能列表排序 | 有多条技能 | GET /api/v1/skills?sort=stars&order=desc | 按 stars 降序排列 | P0 |
| TC-FUNC-004 | 技能详情 | 技能 ID 存在 | GET /api/v1/skills/{id} | 返回完整技能信息 | P0 |
| TC-FUNC-005 | 技能详情 - 不存在 | ID 不存在 | GET /api/v1/skills/999999 | 返回 404 | P0 |
| TC-FUNC-006 | 技能详情 - 参数错误 | ID 非数字 | GET /api/v1/skills/abc | 返回参数错误 | P0 |

### 2.2 搜索模块

| ID | 用例名称 | 前置条件 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|----------|--------|
| TC-FUNC-010 | 全文搜索 | Meilisearch 可用 | POST /api/v1/skills/search {query:"test"} | 返回匹配结果 | P0 |
| TC-FUNC-011 | 全文搜索降级 | Meilisearch 不可用 | 同上（模拟 meiliCli==nil） | 降级为 LIKE 查询 | P1 |
| TC-FUNC-012 | 搜索空结果 | 无匹配 | POST /api/v1/skills/search {query:"zzzznotexist"} | 返回空列表 | P0 |
| TC-FUNC-013 | 搜索参数校验 | 空 query | POST /api/v1/skills/search {} | 返回参数错误 | P0 |

### 2.3 用户认证模块

| ID | 用例名称 | 前置条件 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|----------|--------|
| TC-AUTH-001 | 注册成功 | 用户不存在 | POST /api/v1/auth/register {username,email,password} | 返回 JWT + 用户信息 | P0 |
| TC-AUTH-002 | 注册 - 用户已存在 | 用户名已存在 | 同上 | 返回用户已存在 | P0 |
| TC-AUTH-003 | 登录成功 | 用户已注册 | POST /api/v1/auth/login {username,password} | 返回 JWT | P0 |
| TC-AUTH-004 | 登录 - 密码错误 | 用户存在 | POST /api/v1/auth/login {username,password:"wrong"} | 返回密码错误 | P0 |
| TC-AUTH-005 | 登录 - 用户不存在 | 用户不存在 | POST /api/v1/auth/login {username:"nobody"} | 返回用户不存在 | P0 |
| TC-AUTH-006 | JWT 认证成功 | 有有效 Token | GET /api/v1/user/profile (Authorization: Bearer token) | 返回用户信息 | P0 |
| TC-AUTH-007 | JWT 认证失败 | Token 无效/过期 | 同上（伪造 Token） | 返回未授权 | P0 |

### 2.4 用户功能模块

| ID | 用例名称 | 前置条件 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|----------|--------|
| TC-USER-001 | 添加收藏 | 技能存在、用户已登录 | POST /api/v1/user/favorites {skill_id:1} | 返回收藏成功 | P1 |
| TC-USER-002 | 重复收藏 | 已收藏 | 同上 | 返回已收藏 | P1 |
| TC-USER-003 | 取消收藏 | 已收藏 | DELETE /api/v1/user/favorites/{skill_id} | 返回取消成功 | P1 |
| TC-USER-004 | 提交评分 | 技能存在 | POST /api/v1/user/reviews {skill_id,score:5,content} | 返回评分+平均分 | P1 |
| TC-USER-005 | 评分越界 | score=0 | 同上 | 返回参数错误 | P1 |
| TC-USER-006 | 创建 API Key | 已登录 | POST /api/v1/user/api-keys {name:"test"} | 返回 Key | P1 |
| TC-USER-007 | 吊销 API Key | Key 存在 | DELETE /api/v1/user/api-keys/{id} | 返回吊销成功 | P1 |

### 2.5 智能路由模块

| ID | 用例名称 | 前置条件 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|----------|--------|
| TC-ROUTER-001 | 路由匹配 - hybrid | 向量+关键词引擎可用 | POST /api/v1/router/match {query:"analyze excel"} | 返回匹配技能列表 | P0 |
| TC-ROUTER-002 | 路由匹配 - vector | 仅向量引擎 | POST /api/v1/router/match {query:"test", strategy:"vector"} | 返回向量匹配结果 | P0 |
| TC-ROUTER-003 | 路由匹配 - keyword | 仅关键词引擎 | POST /api/v1/router/match {query:"test", strategy:"keyword"} | 返回关键词匹配结果 | P0 |
| TC-ROUTER-004 | 路由匹配 - 空查询 | query 为空 | POST /api/v1/router/match {} | 返回参数错误 | P0 |
| TC-ROUTER-005 | 路由执行 | 技能存在 + LLM 可用 | POST /api/v1/router/execute {query:"analyze", skill_id:1} | 返回执行结果 | P0 |
| TC-ROUTER-006 | 路由执行 - 技能不存在 | skill_id 无效 | POST /api/v1/router/execute {query:"test", skill_id:9999} | 返回内部错误 | P0 |
| TC-ROUTER-007 | 路由反馈 | log_id 存在 | POST /api/v1/router/feedback {log_id:1, score:4} | 返回反馈成功 | P1 |
| TC-ROUTER-008 | RRF 融合排序 | 向量+关键词均有结果 | hybridSearch | 融合后结果数量 <= topK | P0 |

### 2.6 管理后台模块

| ID | 用例名称 | 前置条件 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|----------|--------|
| TC-ADMIN-001 | 管理员获取技能列表 | 管理员登录 | GET /api/v1/admin/skills?status=0 | 返回待审核技能 | P1 |
| TC-ADMIN-002 | 更新技能状态 | 技能存在 | PUT /api/v1/admin/skills/{id}/status {status:1} | 返回更新成功 | P1 |
| TC-ADMIN-003 | 状态越界 | status=5 | 同上 | 返回参数错误 | P1 |
| TC-ADMIN-004 | 触发全量同步 | 管理员 | POST /api/v1/admin/sync/trigger {type:"full"} | 返回同步任务 | P1 |
| TC-ADMIN-005 | 分类管理 | 管理员 | POST /api/v1/admin/categories {name,slug} | 返回创建的分类 | P1 |
| TC-ADMIN-006 | Dashboard 统计 | 有数据 | GET /api/v1/admin/stats/dashboard | 返回统计信息 | P1 |

### 2.7 插件 API 模块

| ID | 用例名称 | 前置条件 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|----------|--------|
| TC-PLUGIN-001 | 热门技能 | 有活跃技能 | GET /api/v1/plugin/hot?limit=10 | 返回按 installs 降序技能 | P1 |
| TC-PLUGIN-002 | 下载技能 | 技能存在 | GET /api/v1/plugin/download?id=1 | 返回技能+下载URL | P1 |
| TC-PLUGIN-003 | 下载技能 - 不存在 | ID 无效 | GET /api/v1/plugin/download?id=9999 | 返回 404 | P1 |
| TC-PLUGIN-004 | 同步状态 | - | POST /api/v1/plugin/sync/status {installed_ids:[1,2]} | 返回技能同步状态 | P1 |

### 2.8 分类模块

| ID | 用例名称 | 前置条件 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|----------|--------|
| TC-CAT-001 | 分类列表 | 有分类数据 | GET /api/v1/categories | 返回分类列表 | P0 |
| TC-CAT-002 | 分类详情 | 分类存在 | GET /api/v1/categories/{id} | 返回分类信息 | P1 |
| TC-CAT-003 | 创建分类 - 重复 slug | slug 已存在 | POST /api/v1/categories {slug:"test"} | 返回分类已存在 | P1 |
| TC-CAT-004 | 创建分类 - slug 自动小写 | 输入大写 | POST /api/v1/categories {slug:"Data-Analysis"} | slug 变为 "data-analysis" | P1 |
| TC-CAT-005 | 更新分类 | 分类存在 | PUT /api/v1/categories/{id} {name:"new"} | 返回更新后的分类 | P1 |

### 2.9 统计模块

| ID | 用例名称 | 前置条件 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|----------|--------|
| TC-STAT-001 | 平台统计 | 有数据 | GET /api/v1/stats | 返回 total_skills, active_skills 等 | P0 |
| TC-STAT-002 | 热门技能 | 有活跃技能 | GET /api/v1/stats/top-skills?sort=stars&limit=10 | 返回 Top 10 | P1 |

---

## 3. 集成测试用例

| ID | 用例名称 | 测试步骤 | 预期结果 | 优先级 |
|----|----------|----------|----------|--------|
| TC-INT-001 | 注册 → 登录 → 获取 Profile | POST register → POST login → GET profile | 完整流程通过 | P0 |
| TC-INT-002 | 搜索 → 查看详情 → 收藏 | POST search → GET detail → POST favorite | 完整流程通过 | P0 |
| TC-INT-003 | 路由匹配 → 路由执行 → 提交反馈 | POST match → POST execute → POST feedback | 完整流程通过 | P0 |
| TC-INT-004 | 管理员触发同步 → 查看任务状态 | POST trigger → GET tasks | 完整流程通过 | P1 |
| TC-INT-005 | 创建 API Key → 用 Key 调用路由 | POST create → POST router/match (X-API-Key) | Key 认证通过 | P1 |

---

## 4. 非功能测试用例

### 4.1 性能测试

| ID | 用例名称 | 测试条件 | 目标值 | 优先级 |
|----|----------|----------|--------|--------|
| TC-PERF-001 | 路由匹配 API 延迟 | P95，模拟 50 QPS | <= 500ms | P0 |
| TC-PERF-002 | 路由执行 API 延迟 | P95，模拟 10 QPS | <= 3s（含 LLM） | P0 |
| TC-PERF-003 | 全文搜索 P95 延迟 | 模拟 100 QPS | <= 300ms | P0 |
| TC-PERF-004 | 并发支持 | 系统整体 | >= 1000 QPS | P1 |
| TC-PERF-005 | 技能同步 - 增量 | - | <= 10 分钟 | P1 |
| TC-PERF-006 | 技能同步 - 全量 | - | <= 2 小时 | P2 |

### 4.2 安全测试

| ID | 用例名称 | 测试条件 | 预期结果 | 优先级 |
|----|----------|----------|----------|--------|
| TC-SEC-001 | SQL 注入 - 参数注入 | 在 query 参数中注入 SQL | 返回错误而非 SQL 执行 | P0 |
| TC-SEC-002 | SQL 注入 - ORM 安全 | 在 JSON 字段注入 SQL | ORM 参数化查询，注入失败 | P0 |
| TC-SEC-003 | XSS - 输出转义 | 在评论中注入 XSS | 前端转义，不执行脚本 | P1 |
| TC-SEC-004 | JWT 伪造 | 使用伪造 JWT | 认证失败 | P0 |
| TC-SEC-005 | API Key 未授权访问 | 无 Key 调用路由 API | 返回未授权 | P0 |
| TC-SEC-006 | 越权 - 非管理员调用管理 API | 普通用户 Token | 返回权限不足 | P0 |
| TC-SEC-007 | 密码加密存储 | 检查数据库 | 密码以 bcrypt hash 存储 | P0 |
| TC-SEC-008 | 速率限制 | 超限请求 | 返回 429 | P1 |

### 4.3 路由准确率评测

| ID | 用例名称 | 测试条件 | 目标值 | 优先级 |
|----|----------|----------|--------|--------|
| TC-EVAL-001 | Top-1 匹配准确率 | 人工评测集（>= 200 条） | >= 85% | P0 |
| TC-EVAL-002 | Top-3 匹配准确率 | 人工评测集 | >= 95% | P1 |
| TC-EVAL-003 | RRF 融合效果 | 对比纯向量 vs 混合 | 混合优于纯向量 | P1 |
| TC-EVAL-004 | 重排序提升 | 对比粗排 vs 精排 | 精排 Top-1 提升 >= 5% | P1 |

---

## 5. 测试数据需求

### 5.1 种子数据

- 至少 50 条 Skill（覆盖各种分类、状态、星级）
- 至少 5 个分类
- 至少 3 个用户（含 1 个管理员）
- 至少 10 条收藏记录
- 至少 10 条评论记录
- 至少 1 个 API Key
- 至少 3 条路由日志
- 至少 3 条同步任务记录

### 5.2 测试账号

| 账号 | 角色 | 用途 |
|------|------|------|
| testuser / test123 | 普通用户 | 功能测试 |
| admin / admin123 | 管理员 | 管理功能测试 |
| apiuser / apikey123 | API 用户 | API Key 测试 |

---

## 6. 回归测试策略

### 6.1 回归范围

- P0 用例：每次提交必须全部通过
- P1 用例：每日构建必须通过
- P2+ 用例：每个迭代结束前验证

### 6.2 触发条件

- 后端代码变更后自动触发单元测试 + 集成测试
- PR 合入前必须 CI 全部通过
- 发布前必须完成全量回归

---

## 7. 缺陷管理

| 严重级别 | 定义 | 响应时间 | 修复时限 |
|----------|------|----------|----------|
| P0 (致命) | 系统崩溃、核心功能不可用 | 立即 | 4 小时 |
| P1 (严重) | 主要功能异常、数据丢失 | 2 小时 | 24 小时 |
| P2 (一般) | 非核心功能异常、UI 问题 | 8 小时 | 3 天 |
| P3 (轻微) | 文案错误、样式细节 | 24 小时 | 下一迭代 |

---

*以上测试用例覆盖 PRD V3.0 全部功能需求和非功能需求，作为测试执行的标准依据。*
