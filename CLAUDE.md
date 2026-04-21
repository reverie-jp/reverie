# Reverie プロジェクトガイド

このファイルは reverie リポジトリで作業する Claude 向けのプロジェクト固有ルール集です。コードを読めば分かることは書かず、**暗黙の前提・命名規則・運用ルール**を中心にまとめています。

## プロダクトスコープ

- reverie は**グループ通話がメインの SNS**
- 2 人体制のため、当初 MVP を縮小し **通話機能に絞った最小構成を本番運用前提** で進める
- 実装スコープ: **ユーザー / フォロー / 通話のみ**（投稿・チャット・通知・検索などは先送り）
- 技術スタック: Go + Connect-go + sqlc + pgx/v5 + PostgreSQL + LiveKit（予定）
- 通話はアカウント無しでも URL 共有で参加可能、ホストが可視性（非公開 / フォロー中のみ / アカウント必須 / 匿名可）を制御する設計

新機能の提案は「user / follow / call スコープに収まるか」をまず確認。投稿・チャット等を前提にした設計提案は対象外。

## 開発コマンドは devcontainer 内で実行

ホストマシンを汚さない方針。`buf` / `sqlc` / `migrate` / `go` などはすべて devcontainer 内で実行する:

```bash
docker compose -p reverie_devcontainer -f /path/to/reverie/.devcontainer/compose.yaml exec -T reverie-api <command>
```

- compose project: `reverie_devcontainer`（VS Code Dev Containers 拡張と同じ名前を維持）
- service: `reverie-api`
- コンテナ内 workdir: `/reverie-api`
- DB コンテナ: `reverie-db`（同 compose 内）

コンテナが停止していたら `up -d` で起動。代表的なコマンドは `Makefile` にあるので（`make proto` / `make sqlc` / `make migrate-up` / `make dev-up`）そちらを `exec` 経由で呼ぶのが楽。

## 命名規約: Google AIP に寄せる

将来公開 API として外部配布するので、新規 proto は AIP 準拠を前提に書く。既存の既知準拠項目:

- **AIP-142 タイムスタンプ**: フィールドは `_time` で終える（`create_time` / `update_time` / `expire_time` / `first_join_time` など）。`_at` は使わない
- **AIP-122 / 131 リソース名**: リソースは `name` フィールド (`collection/id` 階層文字列) を第一識別子にする。例: `users/{custom_id}`, `calls/{ulid}`, `calls/{ulid}/participants/{identity}`
  - 全ての Get / Update / Delete / サブリソース系 RPC は `name` をパラメータに取り、HTTP パスは `/v1/{name=collection/*}` 形式
  - name の format/parse は `internal/platform/resourcename`（Go）と `web/app/lib/resource-name.ts`（TS）に集約。adapter 層で変換し、usecase / entity は内部 ID のみ扱う
  - リソース参照は ID 文字列ではなく相手の resource message を埋め込む（例: `Call.host: User`、`host_user_id` のような生 ULID は持たない）
- **AIP-158 ページネーション**: 全ての `List*` は `page_size` / `page_token` リクエスト + `next_page_token` レスポンスを備える。カーソルは ULID を opaque として使用、LIMIT は `pageSize + 1` ではなく「DB 返却数 == pageSize なら次がある」判定
- **AIP-136 カスタムメソッド**: `:verb` 形式（`/v1/calls/{...}:join` など）

## マイグレーション戦略

**本番運用前なので `migrations/000001_init.{up,down}.sql` を直接編集**する。追加ファイル（`000002_*` など）は作らない。

- 新規テーブル / カラムは `000001_init.up.sql` に追記し、対応する DROP を `000001_init.down.sql` に追記（up/down 対称）
- 本番リリースが決まった時点でこのルールは変える（そのときは本ファイルも更新）

## アーキテクチャパターン

モジュール構造（`internal/modules/<domain>/`）:

```
<domain>/
  bootstrap.go          InitModule(...) → service handler を返す
  handler/              Connect handler (3 行構造)
  adapter/              境界変換: FromXxxRequest / ToXxxResponse
  usecase/              ビジネスロジック + Input/Output 型
  gateway/              他モジュールへ公開する読み取り系 API（UserView 等 View 構築含む）
  repository/           DB 永続化
  sql/                  sqlc クエリ定義
```

### handler は 3 行構造

```go
func (h *Handler) GetUser(ctx, req) (resp, err) {
    input, err := adapter.FromGetUserRequest(ctx, req)
    if err != nil { return nil, err }
    output, err := h.getUser.Execute(ctx, input)
    if err != nil { return nil, err }
    return adapter.ToGetUserResponse(output), nil
}
```

認証コンテキスト抽出（`interceptor.UserIDFromContext`）は adapter 内で行い、handler には漏らさない。

### adapter パッケージの役割

- トランスポート（Connect/proto）⇔ usecase の境界変換を一手に引き受ける
- `FromXxxRequest(ctx, req) (usecase.XxxInput, error)`
- `ToXxxResponse(output) *connect.Response[...]`
- resource 変換は別ファイル（例: `adapter/user.go` に `ToUser(view) *userv1.User`）
- 小さなヘルパー（`toTokenPair`, `toProviderString` 等）は package private で置く

### usecase の Input / Output 型

- **Input は値型**（`struct` / `XxxInput`）: エラー時は `XxxInput{}` のゼロ値を返す（nil は型エラー）
- **Output は `*XxxOutput`**: nil で「結果なし」を表現可能、コピーコスト回避
- Input には `Validate()` を実装、usecase 先頭で呼ぶ
- Input の validation は `reverie.jp/reverie/internal/platform/validation` + 標準 `errors.New`

### エラー規約

- ドメインエラーは `internal/platform/xerrors` に定義（`ErrUserNotFound`, `ErrInvalidRefreshToken` 等）
- 各 error は connect code 付き: `connect.CodeNotFound`, `connect.CodeUnauthenticated` など
- 詳細メッセージは `.WithMessage(...)` / 原因は `.WithCause(err)` で構築

## proto スタイル

- `proto/<service>/v1/*.proto`
- `buf.gen.yaml` で Go 生成 → `internal/gen/pb/` に出力
- Connect (`*connect.go`) / grpc-gateway (`*.pb.gw.go`) 両方を使う
- メソッドに HTTP annotation を必ず付ける（REST 対応のため）
- 再生成: `make proto`（= `cd proto && buf generate`）

## 認証ルーティング

`interceptor.AuthInterceptor` が 3 カテゴリで RPC を振り分ける:

- **`publicProcedures`**: 認証処理自体をスキップ。Authorization ヘッダは無視される。`SocialLogin` / `RefreshToken` のみ
- **`optionalAuthProcedures`**: ヘッダが無ければゲスト扱いで通過、あれば検証必須。ゲストも使える全ての RPC（`GetUser` / `GetCall` / `JoinCall` / `HeartbeatCall` / `LeaveCall` / `ListPublicCalls` / `GetUserParticipatingCall` 等）
- **どちらにも無い procedure**: 常に認証必須（デフォルト）

新しい RPC を追加したらまず「ゲストも叩けるか？」を判断し、適切な map に追加する。デフォルト「認証必須」に落とさない。

## ゲストユーザーモデル

- **ゲストは DB に一切保存しない**（`users` / `user_auth_providers` 等には入れない）
- LiveKit identity はサーバーが生成する `guest:<ULID>`（クライアント採番しない → なりすまし防止）
- ゲストの表示名などはクライアントの `localStorage` のみに置く
- 単一通話制約は **認証ユーザーだけ**（ゲストは毎回別 identity なので被らない）

## 通話アクティブ状態の真実源

LiveKit の管理 API（ListRooms / ListParticipants）は使わず、`call_participants` の `last_seen_time` / `disconnected_time` を真実源にする:

- クライアントは接続中 **30 秒ごとに `HeartbeatCall`** を叩く
- 意図的退出時は `LeaveCall`、ページ unload 時は `navigator.sendBeacon` 相当（fetch keepalive）で `LeaveCall`
- サーバー判定: `last_seen_time > NOW() - 60s AND disconnected_time IS NULL`（missing 1 heartbeat 許容）
- LiveKit API の唯一の呼び出しは `CreateJoinToken`（トークン発行）のみ。admin API を usecase から呼ばない

新しい通話関連機能を足す際、「LiveKit 管理 API で取れる」と思っても DB 側で表現する方針を維持する。

## その他

- コメントは原則書かない（コードで意図が伝わる名前にする）。非自明な why が必要なときだけ短く
- タイムゾーン: DB は `TIMESTAMPTZ`、Go は `time.Time`
- ULID は `internal/platform/ulid` の自前型。sqlc の生成型とやり取りする際は `.String()` で文字列化が必要なケースがある（既存コード参照）
- sqlc でドメイン型（`ulid` など）を扱う WHERE 句は `sqlc.arg(x)::ulid` のようにキャストを明示する。`$1` だけだと Go 側の型が `string` に落ちてしまう
