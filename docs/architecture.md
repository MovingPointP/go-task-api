# アーキテクチャ図

## レイヤー構成

クリーンアーキテクチャ（Clean Architecture）の考え方に基づき、4つの層に分離されています。
各層は内側の層にのみ依存し、外側の層を知りません。

```mermaid
graph TD
    subgraph "外側"
        H[Handler層<br/>internal/handler]
        I[Infrastructure層<br/>internal/infrastructure]
    end
    subgraph "内側"
        U[Usecase層<br/>internal/usecase]
        D[Domain層<br/>internal/domain]
    end
    subgraph "共通"
        P[pkg/jwt]
    end

    H -->|インターフェース経由| U
    I -->|インターフェース実装| D
    U -->|インターフェース経由| D
    H --> P
    U --> P
```

## 依存関係の詳細

```mermaid
graph LR
    Main[cmd/server/main.go] --> H
    Main --> I
    Main --> U

    H[AuthHandler<br/>TaskHandler] --> U_IF[AuthUsecase I/F<br/>TaskUsecase I/F]
    U[authUsecase<br/>taskUsecase] --> R_IF[UserRepository I/F<br/>TaskRepository I/F]
    I[userRepository<br/>taskRepository] --> R_IF
    U_IF -.実装.-> U
    R_IF -.実装.-> I

    U --> E[entity.User<br/>entity.Task]
    I --> E
    H --> E
```

## 各層の責務

| 層 | パッケージ | 責務 |
|---|---|---|
| **Handler 層** | `internal/handler` | HTTP リクエストのパース、バリデーション、レスポンス生成 |
| **Middleware 層** | `internal/handler/middleware` | JWT 検証、コンテキストへのユーザー ID 設定 |
| **Usecase 層** | `internal/usecase` | ビジネスロジック（重複チェック、パスワードハッシュ化、所有権確認）|
| **Domain 層** | `internal/domain` | エンティティ定義、リポジトリインターフェース、ドメインエラー |
| **Infrastructure 層** | `internal/infrastructure` | GORM を使った DB 操作の実装 |
| **pkg 層** | `pkg/jwt` | JWT 生成・検証ユーティリティ |

## リクエスト処理フロー

```mermaid
sequenceDiagram
    participant C as Client
    participant MW as AuthMiddleware
    participant H as Handler
    participant U as Usecase
    participant R as Repository
    participant DB as PostgreSQL

    C->>MW: HTTP Request (Bearer Token)
    MW->>MW: JWT 検証・UserID 取得
    MW->>H: ctx に UserID をセット
    H->>H: リクエストボディバインド・バリデーション
    H->>U: ビジネスロジック呼び出し
    U->>R: データアクセス
    R->>DB: SQL クエリ (GORM)
    DB-->>R: 結果
    R-->>U: エンティティ
    U-->>H: エンティティ or エラー
    H-->>C: JSON レスポンス
```

## DI（依存性注入）

`cmd/server/main.go` でコンストラクタ注入により全依存関係を組み立てています。

```
DB → UserRepository → AuthUsecase → AuthHandler → Router
DB → TaskRepository → TaskUsecase → TaskHandler → Router
```
