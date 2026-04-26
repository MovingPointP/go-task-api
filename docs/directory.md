# ディレクトリ構成

```
go-task-api/
├── cmd/
│   └── server/
│       └── main.go          # エントリーポイント・DI組み立て
├── docs/                    # Swagger 自動生成ファイル & 設計書
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── task.go      # Task エンティティ・ドメインエラー定義
│   │   │   └── user.go      # User エンティティ・ドメインエラー定義
│   │   └── repository/
│   │       ├── task_repository.go   # TaskRepository インターフェース
│   │       └── user_repository.go   # UserRepository インターフェース
│   ├── handler/
│   │   ├── auth_handler.go  # 認証ハンドラー (register / login)
│   │   ├── task_handler.go  # タスクハンドラー (CRUD)
│   │   ├── router.go        # ルーティング定義・CORS設定
│   │   └── middleware/
│   │       └── auth.go      # JWT認証ミドルウェア
│   ├── infrastructure/
│   │   ├── database/
│   │   │   └── gorm.go      # DB接続・AutoMigrate
│   │   └── persistence/
│   │       ├── task_repository.go   # TaskRepository 実装
│   │       └── user_repository.go   # UserRepository 実装
│   └── usecase/
│       ├── auth_usecase.go  # 認証ユースケース (register / login)
│       └── task_usecase.go  # タスクユースケース (CRUD)
├── pkg/
│   └── jwt/
│       └── jwt.go           # JWT生成・検証ユーティリティ
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

## 各フォルダの役割

| フォルダ | 役割 |
|---|---|
| `cmd/server` | アプリケーションのエントリーポイント。依存性注入（DI）を行いサーバーを起動する |
| `internal/domain/entity` | ドメインオブジェクト（User, Task）とドメインエラーを定義する |
| `internal/domain/repository` | リポジトリのインターフェースを定義する（実装は infrastructure 層に委ねる） |
| `internal/handler` | HTTP リクエスト/レスポンスの処理。ルーティング定義と CORS 設定を含む |
| `internal/handler/middleware` | ルーターに組み込むミドルウェア（JWT 認証） |
| `internal/usecase` | ビジネスロジック層。インターフェースを通じてリポジトリを利用する |
| `internal/infrastructure/database` | DB 接続と GORM AutoMigrate の設定 |
| `internal/infrastructure/persistence` | リポジトリインターフェースの GORM 実装 |
| `pkg/jwt` | JWT の生成・検証ロジック（内部パッケージに依存しない汎用ユーティリティ） |
