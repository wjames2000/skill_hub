# SkillHub Pro 上线前 Checklist

## 基础设施

- [ ] K8s 集群健康检查通过（所有 Node Ready）
- [ ] Namespace `skill-hub` 已创建并配置正确
- [ ] ConfigMap `skill-hub-config` 配置项已按生产环境调整
- [ ] Secret `skill-hub-secret` 中敏感信息已替换（JWT_SECRET、GITHUB_TOKENS 等）
- [ ] PostgreSQL 数据库已创建并可访问
- [ ] Redis 服务已启动并可访问
- [ ] Meilisearch 服务已启动并创建索引
- [ ] MinIO 服务已启动并创建 bucket
- [ ] Consul 服务已注册
- [ ] Milvus 向量数据库已就绪
- [ ] Ollama 模型服务可用（Embedding、LLM、Reranker）

## 数据与迁移

- [ ] 数据库迁移脚本已全部执行（`001_init_sync_tables.sql`, `002_vector_router_tables.sql`）
- [ ] 迁移执行记录已确认，无报错
- [ ] Meilisearch 索引已创建（skills）
- [ ] Milvus Collection 已创建
- [ ] MinIO bucket `skills` 已创建

## CI/CD

- [ ] GitHub Actions CI 流水线通过（lint + build）
- [ ] Docker 镜像已构建并推送到 GitHub Container Registry
- [ ] CI/CD 部署流水线配置完成
- [ ] KUBE_CONFIG 已配置为 GitHub Secret
- [ ] 首次手动触发部署验证成功

## 服务部署

- [ ] router-api Deployment 已部署并 Running（至少 2 副本）
- [ ] sync-worker Deployment 已部署并 Running
- [ ] admin-api Deployment 已部署并 Running
- [ ] frontend Deployment 已部署并 Running
- [ ] 所有 Service 已创建（ClusterIP 正确）
- [ ] Ingress 已配置并指向正确 Service
- [ ] Ingress TLS 证书已申请（Let's Encrypt）
- [ ] 域名 DNS 解析已配置（指向 Ingress Controller IP）

## 健康检查与监控

- [ ] router-api /health 端点返回 200
- [ ] admin-api /health 端点返回 200
- [ ] Prometheus 已部署并能抓取指标
- [ ] Grafana 已部署并配置 Prometheus 数据源
- [ ] Grafana Dashboard 已导入（SkillHub Pro Overview）
- [ ] 告警规则已配置（ServiceDown、HighErrorRate、HighLatency 等）
- [ ] 告警通知渠道已配置（邮件/钉钉/企业微信等）

## 功能验证

- [ ] 前端页面可正常访问（HTTP 200）
- [ ] API 接口正常响应（`GET /skills`、`GET /skills/{id}`）
- [ ] 搜索功能正常（全文搜索 + 语义搜索）
- [ ] 用户登录/注册流程正常
- [ ] GitHub OAuth 登录正常
- [ ] 智能路由接口（`/router/match`、`/router/execute`）正常
- [ ] 技能同步任务可触发并完成
- [ ] 管理后台可访问（`/admin`）
- [ ] 插件 API 正常响应

## 安全

- [ ] JWT Secret 已替换为强随机字符串
- [ ] GitHub Token 已配置
- [ ] 数据库密码已替换
- [ ] MinIO 访问密钥已替换
- [ ] API Key 管理功能正常
- [ ] HTTPS 已开启（TLS 证书配置完成）
- [ ] CORS 配置正确（仅允许前端域名）

## 性能与容量

- [ ] 负载测试通过（API 并发 ≥ 100 QPS）
- [ ] 数据库连接池配置合理（max_open: 25）
- [ ] 容器资源 requests/limits 配置合理
- [ ] HPA（水平自动扩缩容）已配置（可选）
- [ ] 日志收集与轮转已配置

## 回滚准备

- [ ] 上一版本 Docker 镜像已保留（`latest` 外保留版本标签）
- [ ] 数据库迁移脚本可回滚（如需要，准备反向 SQL）
- [ ] 回滚方案文档已准备

## 上线执行

- [ ] 上线时间窗口已确认（非高峰期）
- [ ] 通知相关干系人
- [ ] 手动触发 Deploy Workflow（`workflow_dispatch` → `production`）
- [ ] 观察监控面板 30 分钟，确认无异常
- [ ] 执行验证脚本：`deployments/scripts/verify-deployment.sh`

## 上线后

- [ ] 验证全链路端到端流程
- [ ] 执行一次技能全量同步
- [ ] 确认数据正确写入
- [ ] 确认日志无异常
- [ ] 发布上线通知
- [ ] 更新项目状态文档

---

**签字确认**

| 角色 | 姓名 | 日期 | 签字 |
| --- | --- | --- | --- |
| 后端负责人 | | | |
| 前端负责人 | | | |
| 测试负责人 | | | |
| 运维负责人 | | | |
| 产品经理 | | | |
