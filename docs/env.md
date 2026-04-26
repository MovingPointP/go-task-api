# 環境変数一覧

| 変数名 | 用途 | 必須/任意 | デフォルト値 |
|---|---|---|---|
| `DATABASE_URL` | PostgreSQL 接続 DSN | 必須 | なし |
| `JWT_SECRET` | JWT 署名鍵 (HS256) | 必須 | なし |
| `JWT_EXPIRATION_HOURS` | JWT トークンの有効期限（時間） | 任意 | `72` |
| `PORT` | サーバーのリッスンポート | 任意 | `8080` |
| `HOST` | Swagger UI 表示用ホスト名 | 任意 | `localhost:{PORT}` |

## 設定例 (.env)

```env
DATABASE_URL=postgres://user:password@localhost:5432/taskdb?sslmode=disable
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION_HOURS=72
PORT=8080
```

## 注意事項

- `DATABASE_URL` が未設定の場合、起動時に `fatal: failed to connect to database` でプロセスが終了する
- `JWT_SECRET` が未設定の場合、トークン生成・検証時にエラーが返される
- `.env` ファイルは `godotenv` により起動時に自動読み込みされるが、OS 環境変数が優先される
