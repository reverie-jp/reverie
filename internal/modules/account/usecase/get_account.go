package usecase

import (
	"context"
	"fmt"
	"time"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type GetAccountOutput struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
	AvatarURL   *string
	CreateTime  time.Time
}

type GetAccount struct {
	q *sqlc.Queries
}

func NewGetAccount(q *sqlc.Queries) *GetAccount {
	return &GetAccount{q: q}
}

func (uc *GetAccount) Execute(ctx context.Context, userID ulid.ULID) (*GetAccountOutput, error) {
	user, err := uc.q.GetUserByID(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &GetAccountOutput{
		ID:          user.ID,
		CustomID:    user.CustomID,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarUrl,
		CreateTime:  user.CreateTime,
	}, nil
}
