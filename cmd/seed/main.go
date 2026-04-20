// cmd/seed: 開発用テストユーザーを作成し、アクセストークンを出力するツール
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"reverie.jp/reverie/internal/gen/sqlc"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/ulid"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://reverie:reverie@localhost:5433/reverie_db"
	}
	jwtSecret := os.Getenv("AUTH_JWT_SECRET_KEY")
	if jwtSecret == "" {
		jwtSecret = "dev_secret_key_change_in_production"
	}

	ctx := context.Background()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	q := sqlc.New(db)
	gw := usergw.New(q)
	jwtManager := jwt.NewManager(jwtSecret, 24*time.Hour, 7*24*time.Hour)

	// ── ユーザー定義 ──────────────────────────────────────────────
	type userDef struct {
		customID    string
		displayName string
		biography   string
		posts       []string
	}

	users := []userDef{
		{
			customID:    "keppi",
			displayName: "けっぴ",
			biography:   "元気いっぱい✌️ 趣味はゲームとアニメ鑑賞。深夜に突然ラーメン食べたくなるタイプ",
			posts: []string{
				"今日はゲームで徹夜してしまった😇 つらい",
				"新しいアニメ観たけど神すぎる…続き気になって眠れない",
				"ラーメン食べたい気持ちが止まらない 🍜",
				"好きなゲームの新作来た！！！テンション上がりすぎて仕事手につかない笑",
				"眠いのに眠れない夜ってなんなんだろ",
			},
		},
		{
			customID:    "kaoru",
			displayName: "かおる",
			biography:   "のんびり生きてます 🌙 読書と音楽が好き。夜の散歩が日課",
			posts: []string{
				"今日読み終わった本、久しぶりに泣いた 📖",
				"夜の散歩って最高だな。静かで気持ちいい 🌙",
				"コーヒー飲みながら音楽聴く時間が一番幸せかもしれない ☕",
				"最近ちょっと疲れてたけど、今日は元気 小さな回復",
				"好きな曲を繰り返し聴いてたら気づいたら2時間経ってた",
			},
		},
		{
			customID:    "testuser",
			displayName: "テストユーザー",
			biography:   "",
			posts:       []string{},
		},
		{
			customID:    "nyamada",
			displayName: "にゃまだ",
			biography:   "猫好きエンジニア🐱 コードと猫に囲まれた生活を送ってます。たまに飯テロ投稿します",
			posts: []string{
				"うちの猫がキーボードの上に乗ってきてコードが全部にゃにゃにゃになった",
				"猫の寝顔見てたら仕事やる気なくなってきた これはしょうがない",
				"新しいライブラリ試したら思ったより簡単で拍子抜け 🐾",
				"猫缶の種類を増やしたら気に入ってくれた 飼い主冥利に尽きる",
				"深夜のコーディングお供は猫と緑茶です",
			},
		},
		{
			customID:    "rewrite",
			displayName: "リライト",
			biography:   "文章書くのが好き ✍️ 小説・エッセイ・たまに詩も。言葉を探す日々",
			posts: []string{
				"書きかけの小説が止まってる。主人公が何をしたいのか自分でもわからなくなってきた",
				"いい文章読んだあとってしばらく何も書けなくなるの自分だけ？",
				"カフェで書いてたら隣の人の会話が気になりすぎて全然進まなかった ☕",
				"タイトルだけ思い浮かんで中身が全然ない作品が5本溜まってる",
				"今日は500文字だけ書けた。それでもゼロよりましだと言い聞かせる",
			},
		},
		{
			customID:    "homagusa",
			displayName: "ほま草",
			biography:   "植物と料理と昼寝 🌿 自給自足の夢を見ながら都会に住んでいます",
			posts: []string{
				"ベランダのバジルが育ちすぎてパスタに追いつかない 🌿",
				"今日の昼ごはんは自家製味噌で作った味噌汁。最高すぎて泣きそう",
				"スーパーで買った野菜の根っこ植えたら芽が出た 感動",
				"昼寝しすぎて夜眠れなくなるの繰り返しをもう3年やってる",
				"料理してる時間が一番頭空っぽになれる気がする 🍳",
			},
		},
		{
			customID:    "markun",
			displayName: "まぁくん",
			biography:   "スポーツ観戦とドライブが好き 🚗 週末はだいたい外にいます",
			posts: []string{
				"昨日の試合熱すぎた！！あの逆転劇は一生語り継ぐ",
				"久々に遠出したら知らない道に迷い込んでそれがすごくよかった",
				"スタジアムで食べるチキンってなんであんなにうまいんだろ 🍗",
				"早起きしてドライブしたら朝日がやばかった。写真撮れなかったのが悔やまれる",
				"好きなチームが負けて落ち込んでるので慰めて",
			},
		},
	}

	// ── 各ユーザーを作成 or 取得 ──────────────────────────────────
	type createdUser struct {
		def userDef
		id  ulid.ULID
	}

	created := make([]createdUser, 0, len(users))

	for _, u := range users {
		existing, err := gw.GetUserByCustomID(ctx, u.customID)
		if err != nil {
			log.Fatalf("failed to check existing user (%s): %v", u.customID, err)
		}

		var userID ulid.ULID
		if existing != nil {
			userID = existing.ID
			fmt.Fprintf(os.Stderr, "既存ユーザーを使用: %s (@%s)\n", existing.DisplayName, existing.CustomID)
		} else {
			userID = ulid.New()
			if err := gw.CreateUser(ctx, usergw.CreateUserParams{
				ID:          userID,
				CustomID:    u.customID,
				DisplayName: u.displayName,
			}); err != nil {
				log.Fatalf("failed to create user (%s): %v", u.customID, err)
			}
			fmt.Fprintf(os.Stderr, "新規ユーザーを作成: %s (@%s)\n", u.displayName, u.customID)
		}

		// biography を更新（初回だけでなく毎回上書きして最新にする）
		if u.biography != "" {
			if _, err := gw.UpdateUser(ctx, usergw.UpdateUserParams{
				ID:          userID,
				DisplayName: u.displayName,
				Biography:   u.biography,
				IsPrivate:   false,
			}); err != nil {
				log.Fatalf("failed to update user biography (%s): %v", u.customID, err)
			}
		}

		created = append(created, createdUser{def: u, id: userID})
	}

	// ── 投稿を作成（既存ユーザーはスキップ、初回のみ） ──────────────
	for _, cu := range created {
		if len(cu.def.posts) == 0 {
			continue
		}

		// 既存の投稿があるかチェック（1件でもあればスキップ）
		existingPosts, err := q.ListUserPosts(ctx, sqlc.ListUserPostsParams{
			AuthorID: cu.id,
			Column2:  time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
			Limit:    1,
		})
		if err != nil {
			log.Fatalf("failed to check existing posts for %s: %v", cu.def.customID, err)
		}
		if len(existingPosts) > 0 {
			fmt.Fprintf(os.Stderr, "投稿スキップ（既存あり）: @%s\n", cu.def.customID)
			continue
		}

		// 投稿を時系列順で作成（古い順に入れる）
		baseTime := time.Now().Add(-time.Duration(len(cu.def.posts)) * time.Hour)
		for i, text := range cu.def.posts {
			postID := ulid.New()
			postTime := baseTime.Add(time.Duration(i) * time.Hour)
			if _, err := db.Exec(ctx,
				`INSERT INTO posts (id, author_id, text, create_time, update_time) VALUES ($1, $2, $3, $4, $4)`,
				postID.String(), cu.id.String(), text, postTime,
			); err != nil {
				log.Fatalf("failed to create post for %s: %v", cu.def.customID, err)
			}
		}
		fmt.Fprintf(os.Stderr, "投稿を作成: @%s (%d件)\n", cu.def.customID, len(cu.def.posts))
	}

	// ── 全ユーザーのトークンを出力 ──────────────────────────────────
	for _, cu := range created {
		accessToken, err := jwtManager.GenerateAccessToken(cu.id)
		if err != nil {
			log.Fatalf("failed to generate access token for %s: %v", cu.def.customID, err)
		}
		fmt.Printf("# @%s (%s)\n", cu.def.customID, cu.def.displayName)
		fmt.Printf("USER_ID=%s\n", cu.id.String())
		fmt.Printf("ACCESS_TOKEN=%s\n\n", accessToken)
	}
}
