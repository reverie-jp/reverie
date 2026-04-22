# Reverie プロジェクトガイド

コードを読めば分かることは書かない。**暗黙の前提・命名規則・運用ルール**のみ。

## プロダクトスコープ

- **グループ通話 SNS**、2 人体制で通話特化の最小構成を本番前提で進める
- 実装スコープは **ユーザー / フォロー / 通話のみ**。投稿・チャット・通知・検索は先送り
- スタック: Go + Connect-go + sqlc + pgx/v5 + PostgreSQL + LiveKit
- 通話はアカウント無しでも URL 共有で参加可能、ホストが可視性 (非公開 / フォロー中のみ / アカウント必須 / 匿名可) を制御

新機能提案は「user / follow / call に収まるか」をまず確認。

## 開発コマンドは devcontainer 内で実行

ホストを汚さない方針。`buf` / `sqlc` / `migrate` / `go` はすべて devcontainer 経由:

```bash
docker compose -p reverie_devcontainer -f /path/to/reverie/.devcontainer/compose.yaml exec -T reverie-api <command>
```

compose project `reverie_devcontainer`, service `reverie-api`, DB コンテナ `reverie-db`。代表的コマンドは `Makefile` (`make proto` / `make sqlc` / `make migrate-up` / `make dev-up`)。

## 命名規約 (Google AIP)

将来公開 API なので新規 proto は AIP 準拠:

- **AIP-142**: タイムスタンプは `_time` 末尾 (`create_time`, `first_join_time`)。`_at` 禁止
- **AIP-122/131**: リソースは `name` (`collection/id`) を第一識別子。Get/Update/Delete は `name` パラメータ、HTTP は `/v1/{name=collection/*}`。name の format/parse は `internal/platform/resourcename` (Go) と `web/app/lib/resource-name.ts` (TS) に集約
- リソース参照は ID 文字列でなく相手の resource message を埋める (`Call.host: User`、生 ULID は持たない)
- **AIP-158**: `List*` は `page_size` / `page_token` / `next_page_token`。カーソルは ULID を opaque 使用、「DB 返却数 == pageSize なら次あり」判定
- **AIP-136**: カスタム動詞は `:verb` (`/v1/calls/{...}:join`)

## マイグレーション

**本番前なので `migrations/000001_init.{up,down}.sql` を直接編集**。追加ファイルは作らない。up/down 対称に保つ。本番リリース後はルール変更予定 (そのとき本ファイルも更新)。

## モジュール構造

```
internal/modules/<domain>/
  bootstrap.go    InitModule(...) → service handler
  handler/        Connect handler (3 行) — 1 RPC = 1 ファイル
  adapter/        境界変換 — 1 RPC = 1 ファイル
  usecase/        ビジネスロジック — 1 RPC = _input.go / _output.go / .go の 3 ファイル
  gateway/        他モジュール公開の読み取り API + View 組み立て
  repository/     DB 永続化 (entity のみ返す)
  sql/            sqlc クエリ
```

### handler は 3 行

```go
func (h *Handler) GetUser(ctx, req) (resp, err) {
    input, err := adapter.FromGetUserRequest(ctx, req)
    if err != nil { return nil, err }
    output, err := h.getUser.Execute(ctx, input)
    if err != nil { return nil, err }
    return adapter.ToGetUserResponse(output), nil
}
```

認証コンテキスト抽出 (`interceptor.UserIDFromContext`) は adapter 内で行い、handler に漏らさない。

### usecase

- **Input は値型**: エラー時は `XxxInput{}` のゼロ値返却。`Validate()` を実装し usecase 先頭で呼ぶ
- **Output は `*XxxOutput`**: nil で「結果なし」。基本的に gateway の View 型 (`*callgw.CallView` など) を包むだけ。usecase 内に view 型を新設しない
- validation は `internal/platform/validation` + 標準 `errors.New`

### エラー

`internal/platform/xerrors` にドメインエラー定義 (connect code 付き)。`.WithMessage(...)` / `.WithCause(err)` で構築。

## Repository / Gateway / View の境界 (厳守)

崩れると usecase が肥大化する。

**Repository**: `*entity.Xxx` (DB 行に 1:1) か `[]ulid.ULID` / `bool` のみ。View を返さない。複数 table を JOIN する SQL は OK (`ListActiveCallIDsForFollower` が user_follows を JOIN など) だが、返すのは単一 entity か IDs まで。List API は **SQL 側で ID だけフィルタ** → Go に件数を載せない。

**Gateway**: View 型 (`UserView` / `CallView` / `CallParticipantView` / `CallBanView`) の定義と組み立てを担当。List API の典型は 3 クエリ: (1) repo が ID 配列 (2) `gateway.BuildListXxxViews` が entity を fetch (3) 関連 gateway に委譲して合成。

**View ビルダー命名**: 単数 `BuildXxxView` → `*XxxView` / 複数 `BuildListXxxViews` → `[]*XxxView`。`BuildXxxListViews` や `BuildXxxViews` は使わない。

**Gateway 間の直接依存 OK**: user gateway は follow gateway を import して `RelationshipsByUserIDs` を呼ぶ。Provider interface のような**先回り抽象は作らない**。循環が実際に起きてから対処。

**entity と View の分離**: DB 行由来のものは entity (非正規化カラムも含む)、requester 依存フラグ (`IsMe` / `IsFollowing` など) は View。entity に文脈依存を載せると「同じ id で値が変わる」矛盾になる。

**domain ロジックは entity メソッドに**: `CallParticipant.IsCurrentlyConnected()` のように heartbeat 判定等は entity へ。関連定数 (`entity.ParticipantStaleSeconds`) も entity パッケージ。

## 非正規化カウンタは DB トリガで

`users.following_count` / `follower_count` は `user_follows` の INSERT/DELETE トリガ (`user_follow_counts_sync`) で自動更新。アプリは follow/unfollow でカウンタを意識しない。ドリフト時は `UPDATE users SET following_count = (SELECT COUNT(*) FROM user_follows WHERE follower_id = users.id)` で再集計。CASCADE DELETE とも整合。新規集計カラムも同方針。

## 認証ルーティング

`interceptor.AuthInterceptor` の 3 カテゴリ:

- `publicProcedures`: 認証スキップ。`SocialLogin` / `RefreshToken` のみ
- `optionalAuthProcedures`: ヘッダなしはゲスト通過、あれば検証必須。ゲストも叩ける RPC (`GetUser` / `GetCall` / `JoinCall` / `HeartbeatCall` / `LeaveCall` / `ListPublicCalls` / `GetUserParticipatingCall` / `ListFollowingUsers` / `ListUserFollowers` 等)
- どちらにも無い: 認証必須 (デフォルト)

新規 RPC はまず「ゲストも叩けるか」を判断し適切な map に追加。デフォルト認証必須に**落とさない**。

## ゲストユーザー

- **DB に保存しない** (`users` 等には入れない)
- LiveKit identity はサーバー生成の `guest:<ULID>` (クライアント採番 = なりすまし可能なので禁止)
- 表示名などはクライアント `localStorage` のみ
- 単一通話制約は認証ユーザーだけに適用

## 通話アクティブ状態

LiveKit の管理 API (ListRooms / ListParticipants) を**使わない**。`call_participants.last_seen_time` / `disconnected_time` が真実源:

- クライアントは接続中 30 秒ごとに `HeartbeatCall`
- 意図的退出 / page unload (`sendBeacon` / fetch keepalive) は `LeaveCall`
- サーバー判定: `entity.CallParticipant.IsCurrentlyConnected()` (= `last_seen_time > NOW() - ParticipantStaleSeconds 秒 AND disconnected_time IS NULL`、1 heartbeat 欠損許容)
- LiveKit API で呼ぶのは `CreateJoinToken` と mutation 系 (`DeleteRoom` / `MutePublishedTrack` / `RemoveParticipant`) のみ。query 系の admin API は呼ばない

**ホスト離脱 = 通話終了**: `LeaveCall` usecase は離脱者がホストのとき `MarkCallEnded` + `DeleteRoom` + `MarkAllCallParticipantsDisconnected` を実行。リロード / 明示退出 / beacon いずれでもホスト不在の通話は残らない。

**フォロー中タイムライン (`ListFollowingCalls`)**: ホストまたは現在接続中の参加者のいずれかがフォロー中ユーザーの通話を返す。

## ULID / sqlc の細かい話

- ULID は `internal/platform/ulid` の自前型。`Value()` は zero ULID のとき nil を返す (NULL バインド)。SQL 側で `sqlc.arg(requester_id)::ulid` が NULL → `=` / `EXISTS` が false 評価なので、ゲスト分岐を Go で書かずに済む
- ドメイン型 (`ulid` 等) を使う WHERE 句は `sqlc.arg(x)::ulid` と**キャストを明示**。`$1` だけだと Go 側が `string` に落ちる
- 配列引数は `sqlc.arg(ids)::text[]` で受け、Go 側で `[]ulid.ULID` → `[]string` に変換して渡す
- タイムゾーンは DB `TIMESTAMPTZ` / Go `time.Time`

## コメントポリシー

原則書かない。コードで意図が伝わる名前にする。非自明な why が必要なときだけ短く。
