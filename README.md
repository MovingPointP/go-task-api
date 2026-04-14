# go-task-api

Go + Gin + Supabase で構築するタスク管理 REST API。

## 技術スタック

| カテゴリ | 技術 |
|----------|------|
| 言語 | [Go](https://go.dev/) |
| Webフレームワーク | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| データベース | [Supabase](https://supabase.com/) (PostgreSQL) |
| DBドライバ | [pgx](https://github.com/jackc/pgx) |
| 認証 | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) |
| APIドキュメント | [swaggo/swag](https://github.com/swaggo/swag) (OpenAPI) |
| 環境変数 | [godotenv](https://github.com/joho/godotenv) |
| コンテナ | Docker |
| デプロイ | [Render](https://render.com/) |

## アーキテクチャ

クリーンアーキテクチャを採用し、以下の層で構成されています。

```
go-task-api/
├── cmd/server/          # エントリーポイント
├── docs/                # swag生成のOpenAPIドキュメント
├── internal/
│   ├── domain/
│   │   ├── entity/      # エンティティ（ビジネスオブジェクト）
│   │   └── repository/  # リポジトリインターフェース
│   ├── usecase/         # ビジネスロジック
│   ├── infrastructure/  # DB・外部サービスの実装
│   └── handler/         # HTTPハンドラ
└── pkg/                 # 共有ユーティリティ
```

## 主な機能

- ユーザー登録・ログイン（JWT認証）
- タスクのCRUD操作（認証必須）
- OpenAPI (Swagger UI) によるAPIドキュメント
