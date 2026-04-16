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
	"reverie.jp/reverie/internal/modules/user/gateway"
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
	gw := gateway.New(q)

	// テストユーザーを作成（既存の場合は上書きなし）
	userID := ulid.New()
	customID := "testuser"
	displayName := "テストユーザー"

	// カスタムIDが既に存在するか確認
	existing, err := gw.GetUserByCustomID(ctx, customID)
	if err != nil {
		log.Fatalf("failed to check existing user: %v", err)
	}

	if existing != nil {
		userID = existing.ID
		fmt.Fprintf(os.Stderr, "既存ユーザーを使用: %s (%s)\n", existing.DisplayName, existing.CustomID)
	} else {
		err = gw.CreateUser(ctx, gateway.CreateUserParams{
			ID:          userID,
			CustomID:    customID,
			DisplayName: displayName,
		})
		if err != nil {
			log.Fatalf("failed to create user: %v", err)
		}
		fmt.Fprintf(os.Stderr, "新規ユーザーを作成: %s (%s)\n", displayName, customID)
	}

	// JWTトークンを生成
	jwtManager := jwt.NewManager(jwtSecret, 24*time.Hour, 7*24*time.Hour)
	accessToken, err := jwtManager.GenerateAccessToken(userID)
	if err != nil {
		log.Fatalf("failed to generate access token: %v", err)
	}

	fmt.Printf("USER_ID=%s\n", userID.String())
	fmt.Printf("ACCESS_TOKEN=%s\n", accessToken)
}
