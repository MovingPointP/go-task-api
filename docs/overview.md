# サービス概要

## 目的と責務

go-task-api はユーザーごとにタスクを管理する REST API サービスです。
ユーザー登録・ログインによる認証基盤と、認証済みユーザー向けのタスク CRUD 機能を提供します。

## 技術スタック

| カテゴリ | 技術 |
|---|---|
| 言語 | Go 1.25 |
| Web フレームワーク | Gin v1.12 |
| ORM | GORM v1.31 |
| データベース | PostgreSQL（本番）/ SQLite（テスト） |
| 認証 | JWT (HS256) — golang-jwt/jwt v5 |
| パスワードハッシュ | bcrypt |
| API ドキュメント | Swagger (swaggo v1.16) |
| 設定管理 | godotenv |
| コンテナ | Docker / docker-compose |
