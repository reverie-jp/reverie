# Reverie プロジェクトガイド

コードを読めば分かることは書かない。**暗黙の前提・命名規則・運用ルール**のみ。

## プロダクトスコープ

- **グループ通話 SNS**、2 人体制で通話特化の最小構成を本番前提で進める
- 実装スコープは **ユーザー / フォロー / 通話 / 通知 / プレゼンス**。投稿・チャット・検索は先送り
- スタック: Go + Connect-go + sqlc + pgx/v5 + PostgreSQL + Redis + LiveKit
- 通話はアカウント無しでも URL 共有で参加可能、ホストが可視性 (非公開 / フォロー中のみ / アカウント必須 / 匿名可) を制御

新機能提案は「user / follow / call / notification / presence に収まるか」をまず確認。

## 開発コマンドは devcontainer 内で実行

ホストを汚さない方針。`buf` / `sqlc` / `migrate` / `go` はすべて devcontainer 経由:

```bash
docker compose -p reverie_devcontainer -f /path/to/reverie/.devcontainer/compose.yaml exec -T reverie-api <command>
```

compose project `reverie_devcontainer`, service `reverie-api`, DB コンテナ `reverie-db`, Redis コンテナ `reverie-redis`。代表的コマンドは `Makefile` (`make proto` / `make sqlc` / `make migrate-up` / `make dev-up`)。env 変更時は `up -d --force-recreate` が必要。

**API サーバーは `make dev-up` (= `air`) で起動**。Go ファイル変更を air が自動 rebuild + 再起動してくれるので、`go build` → `kill` → `run` の手動ループは不要。手動でバイナリを起動しないこと (ポート衝突 / 挙動の不一致を招く)。

## 命名規約 (Google AIP)

将来公開 API なので新規 proto は AIP 準拠:

- **AIP-142**: タイムスタンプは `_time` 末尾 (`create_time`, `last_seen_time`)。`_at` 禁止
- **AIP-122/131**: リソースは `name` (`collection/id`) を第一識別子。Get/Update/Delete は `name` パラメータ、HTTP は `/v1/{name=collection/*}`。name の format/parse は `internal/platform/resourcename` (Go) と `web/app/lib/resource-name.ts` (TS) に集約
- リソース参照は ID 文字列でなく相手の resource message を埋める (`Call.host: User`、生 ULID は持たない)
- **AIP-158**: `List*` は `page_size` / `page_token` / `next_page_token`。カーソルは ULID を opaque 使用、「DB 返却数 == pageSize なら次あり」判定
- **AIP-136**: カスタム動詞は `:verb` (`/v1/calls/{...}:join`, `/v1/users/me:heartbeat`)

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

**Repository**: `*entity.Xxx` (DB 行に 1:1) か `[]ulid.ULID` / `bool` のみ。View を返さない。複数 table を JOIN する SQL は OK だが、返すのは単一 entity か IDs まで。List API は **SQL 側で ID だけフィルタ** → Go に件数を載せない。

**Gateway**: View 型 (`UserView` / `CallView` / `CallParticipantView` / `CallBanView` / `NotificationView`) の定義と組み立てを担当。List API の典型は 3 クエリ: (1) repo が ID 配列 (2) `gateway.BuildListXxxViews` が entity を fetch (3) 関連 gateway に委譲して合成。

**View ビルダー命名**: 単数 `BuildXxxView` → `*XxxView` / 複数 `BuildListXxxViews` → `[]*XxxView`。`BuildXxxListViews` や `BuildXxxViews` は使わない。

**Gateway 間の直接依存 OK**: user gateway は follow gateway を import して `RelationshipsByUserIDs` を呼ぶ。Provider interface のような**先回り抽象は作らない**。循環が実際に起きてから対処。

**entity と View の分離**: DB 行由来のものは entity (非正規化カラムも含む)、requester 依存フラグ (`IsMe` / `IsFollowing` など) は View。entity に文脈依存を載せると「同じ id で値が変わる」矛盾になる。

**domain ロジックは entity メソッドに**: `CallParticipant.IsCurrentlyConnected()` や `User.IsCurrentlyOnline()` のように heartbeat 判定等は entity へ。関連定数 (`entity.ParticipantStaleSeconds` / `entity.UserPresenceStaleSeconds`) も entity パッケージ。

## 非正規化カウンタは DB トリガで

`users.following_count` / `follower_count` は `user_follows` の INSERT/DELETE トリガ (`user_follow_counts_sync`) で自動更新。アプリは follow/unfollow でカウンタを意識しない。ドリフト時は `UPDATE users SET following_count = (SELECT COUNT(*) FROM user_follows WHERE follower_id = users.id)` で再集計。CASCADE DELETE とも整合。新規集計カラムも同方針。

## 認証ルーティング

`interceptor.AuthInterceptor` は **Unary と Streaming の両方**を wrap する (`WrapUnary` / `WrapStreamingHandler`)。3 カテゴリ:

- `publicProcedures`: 認証スキップ。`SocialLogin` / `RefreshToken` のみ
- `optionalAuthProcedures`: ヘッダなしはゲスト通過、あれば検証必須。ゲストも叩ける RPC (`GetUser` / `GetCall` / `JoinCall` / `HeartbeatCall` / `LeaveCall` / `ListPublicCalls` / `GetUserParticipatingCall` / `ListFollowingUsers` / `ListUserFollowers` 等)
- どちらにも無い: 認証必須 (デフォルト、`StreamEvents` / `Heartbeat` 等)

新規 RPC はまず「ゲストも叩けるか」を判断し適切な map に追加。デフォルト認証必須に**落とさない**。

## ゲストユーザー

- **DB に保存しない** (`users` 等には入れない)
- LiveKit identity はサーバー生成の `guest:<ULID>` (クライアント採番 = なりすまし可能なので禁止)
- 表示名などはクライアント `localStorage` のみ
- 単一通話制約は認証ユーザーだけに適用

## Heartbeat / Stale Window は 3 倍ルール

接続生存判定は全て **heartbeat 間隔 × 3 = stale window 60s** で統一:

| 用途 | 間隔 | Stale Window | 判定 |
|------|------|--------------|------|
| 通話参加者 | 20s (client → `HeartbeatCall`) | `entity.ParticipantStaleSeconds` = 60s | `CallParticipant.IsCurrentlyConnected()` |
| ユーザー presence | 20s (client → `Heartbeat`) | `entity.UserPresenceStaleSeconds` = 60s | `User.IsCurrentlyOnline()` |

2 回 miss しても切り替わらない余裕を持たせている。片方だけ変える時は必ず両方見ること (entity コメントで相互参照を残してある)。

## 通話アクティブ状態

LiveKit の管理 API (ListRooms / ListParticipants) を**使わない**。`call_participants.last_seen_time` / `disconnected_time` が真実源:

- クライアントは接続中 20s ごとに `HeartbeatCall`
- 意図的退出 / page unload (`sendBeacon` / fetch keepalive) は `LeaveCall`
- LiveKit API で呼ぶのは `CreateJoinToken` と mutation 系 (`DeleteRoom` / `MutePublishedTrack` / `RemoveParticipant`) のみ。query 系の admin API は呼ばない

**ホスト離脱 = 通話終了**: `LeaveCall` usecase は離脱者がホストのとき `MarkCallEnded` + `DeleteRoom` + `MarkAllCallParticipantsDisconnected` を実行。リロード / 明示退出 / beacon いずれでもホスト不在の通話は残らない。

**フォロー中タイムライン (`ListFollowingCalls`)**: ホストまたは現在接続中の参加者のいずれかがフォロー中ユーザーの通話を返す。

## ユーザー presence

- `users.last_seen_time` を真実源とする (`user_presence` のような専用 table は作らない)。`users` は read-heavy だが `last_seen_time` は非 index なので HOT update で index 影響なし
- `PresenceService.Heartbeat` RPC が `UPDATE users SET last_seen_time = NOW()` を叩くだけ
- `UserView.IsOnline` は `BuildListUserViews` 内で `u.IsCurrentlyOnline()` から導出 — **追加クエリなし**
- proto の `OnlineStatus` enum は `ONLINE` / `OFFLINE` のみ出し分け (`IDLE` は未使用)。`last_seen_time` も `User` 本体に同梱 (「最終アクティブ N 分前」UI 用)
- フロントは root で `usePresenceHeartbeat()` を 1 回起動。Page Visibility API で hidden 時は停止

## リアルタイムイベント (EventService.StreamEvents)

**DB が正、stream は hint**: 取りこぼしを許容する設計。クライアントは接続時に `ListNotifications` 等で state を DB fetch → 以降 stream で差分受信。再接続時も同じ。

### 配線

- `internal/platform/events` に `Publisher` / `Subscriber` / `Bus` インターフェースと Redis pub/sub 実装
- `NotificationService.Create` は DB insert → `events.UserTopic(recipient)` に publish。publish 失敗は log のみ (rollback しない)
- `EventService.StreamEvents` は server-streaming。`events.Subscriber` で自分宛 topic (`event:user:{me}`) を購読
- Redis topic prefix `event:` で他用途と名前空間分離
- `EventService` は **grpc-gateway 非対応** (streaming を HTTP/JSON にブリッジできないため)。`service.RegisterGatewayHandler = nil` を許可する分岐が `runner.go` にある

### Keepalive と server HTTP 設定

- Server は 15s ごとに `KeepAliveEvent` を流し、chunked response を生存させる (LB / proxy / NAT の idle timeout 対策)
- Client は `STREAM_STALE_MS = 35s` の watchdog — 任意のフレーム到着でリセット、超過すると abort → 指数バックオフ再接続
- `http.Server.WriteTimeout` は **0 (無効)** 必須。有限値は streaming response を途中切断するため `ERR_INCOMPLETE_CHUNKED_ENCODING` になる
- `http2.Server.ReadIdleTimeout` / `PingTimeout` も併用して HTTP/2 経路の生存も担保

### イベント型の追加

`StreamEventsResponse.oneof payload` に枝を足すだけ。独立 RPC を増やさない。既存 envelope:
- `KeepAliveEvent` (tag 9) — 空、keepalive 専用
- `NotificationCreatedEvent` (tag 10) — 新規通知

## 通知

- 専用モジュール `internal/modules/notification/`、`notifications` テーブル
- 通知作成は **発火元 usecase から notification gateway の `Create()` を呼ぶ** (例: follow usecase が `user_followed` 通知を作成)。gateway 内で DB insert + stream publish が 1 箇所に集約
- dedup: `(recipient_user_id, type, COALESCE(actor_user_id::text, ''), resource_name)` の unique 制約
- **`CreateNotification` の upsert は ON CONFLICT で行を "refresh"** (新 id + `create_time = NOW()` + `read_time = NULL`)。cooldown が通った後に到達する UPDATE パスは「新イベント」として扱う意味なので、client には新しい event_id が届き toast が再発火する
- 大量配信 (call 作成 → フォロワー全員) は `FanOutCreate` で 1 回の `INSERT ... unnest(...)` と concurrency 20 の並列 publish (detached goroutine で caller をブロックしない)
- **トグルスパム対策は cooldown** (`internal/platform/ratelimit` + `gateway/cooldown.go`): Redis SETNX で `(recipient, type, actor, resource)` の再通知を抑止。`user_followed` は 1h、`following_user_call_started` は 0 (resource_name が毎回 unique なので不要)。Redis 障害時は fail open (通知を止めない)
- **逆操作 (unfollow 等) で通知を削除しない**: 歴史的事実として残す。cooldown 期限内の再 follow は silently skip、期限切れ後は行が refresh されて新イベント扱いになる
- **client の論理 dedup**: サーバーの refresh で同じ (type, actor, resource) に新 id が飛んでくるので、`NotificationProvider` は stream 受信時に既存の論理重複を除去してから先頭挿入する
- フロントは `NotificationProvider` が state と stream を管理。`ListNotifications` で初期化 → stream で差分適用 → sonner toast

## ULID / sqlc の細かい話

- ULID は `internal/platform/ulid` の自前型。`Value()` は zero ULID のとき nil を返す (NULL バインド)。SQL 側で `sqlc.arg(requester_id)::ulid` が NULL → `=` / `EXISTS` が false 評価なので、ゲスト分岐を Go で書かずに済む
- ドメイン型 (`ulid` 等) を使う WHERE 句は `sqlc.arg(x)::ulid` と**キャストを明示**。`$1` だけだと Go 側が `string` に落ちる
- 配列引数は `sqlc.arg(ids)::text[]` で受け、Go 側で `[]ulid.ULID` → `[]string` に変換して渡す
- タイムゾーンは DB `TIMESTAMPTZ` / Go `time.Time`

## コメントポリシー

原則書かない。コードで意図が伝わる名前にする。非自明な why が必要なときだけ短く。
