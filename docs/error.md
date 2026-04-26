# エラーハンドリング

## エラーレスポンスの形式

全エラーレスポンスは以下の JSON 形式で返されます。

```json
{
  "error": "エラーメッセージ"
}
```

## HTTP ステータスコードと意味

| ステータスコード | 意味 | 発生箇所の例 |
|---|---|---|
| `400 Bad Request` | リクエストボディのバリデーションエラー、パスパラメータ不正 | ShouldBind 失敗、無効なタスク ID |
| `401 Unauthorized` | 認証失敗 | トークン未提供、無効・期限切れトークン、認証情報不一致 |
| `404 Not Found` | リソースが存在しない | 指定 ID のタスクが存在しないまたは他ユーザーのタスク |
| `409 Conflict` | リソース重複 | メールアドレスがすでに登録済み |
| `500 Internal Server Error` | サーバー内部エラー | DB 操作失敗 |

## ドメインエラー定義

| エラー変数 | エラーメッセージ | 用途 |
|---|---|---|
| `entity.ErrEmailAlreadyInUse` | `"email already in use"` | メールアドレス重複チェック |
| `entity.ErrInvalidCredentials` | `"invalid email or password"` | ログイン認証失敗 |
| `entity.ErrTaskNotFound` | `"task not found"` | 存在しないタスクへのアクセス |

## エラー処理方針

- ドメインエラー（上記3種）は `errors.Is()` で判定し、適切な HTTP ステータスコードにマッピングする
- ユースケース層は内部エラーを `fmt.Errorf("...: %w", err)` でラップして返す
- ハンドラー層はドメインエラーのみ識別し、それ以外は 500 として汎用メッセージを返す（内部詳細を隠蔽）
- 認証ミドルウェアはエラー発生時に `ctx.AbortWithStatusJSON()` を使用してリクエストを中断する
