# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- 初始项目脚手架
- 基础 API 服务框架
- Worker 任务处理框架
- 数据库迁移脚本
- Docker Compose 开发环境
- Repository/Service/Handler 完整主链路实现
- 任务入队与 Worker 状态流转（含重试与死信）
- PDF 解析器（段落定位/页码映射）
- LLM 抽取层（Provider 抽象 + OpenAI 实现 + Schema 校验）
- 规则引擎（主题匹配 + 红黄绿评分）
- 报告生成、Markdown/PDF 导出、历史对比能力
- JWT 鉴权、请求 ID、基础指标与就绪探针
- 前端 MVP 页面骨架与 E2E 冒烟用例
- CI、Issue/PR 模板、pre-commit、运行手册与架构文档
- 规则配置后台（版本发布/回滚/灰度）与管理员页面
- 运营埋点链路（前端 SDK、后端写库、漏斗与周/月看板）
- Python `document-parser` 服务（OCR/报告结构化/条款结构化）
- RAG 模块（条款切片、向量检索、召回重排、离线评估工具）
- i18n 国际化（前端中英文切换、后端错误码国际化、导出多语言模板）
- 发布与回滚交付物（生产配置模板、备份脚本、灰度与回滚演练报告、v1.0.0 release notes）

## [0.1.0] - 2026-02-28

### Added
- 项目初始化
