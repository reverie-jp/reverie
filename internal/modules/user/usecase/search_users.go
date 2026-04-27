package usecase

import (
	"context"
	"encoding/base64"
	"time"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type SearchUsersInput struct {
	Query     string
	PageToken string
	PageSize  int32
}

type SearchUsers struct {
	userGateway usergw.Gateway
}

func NewSearchUsers(userGateway usergw.Gateway) *SearchUsers {
	return &SearchUsers{userGateway: userGateway}
}

func (uc *SearchUsers) Execute(ctx context.Context, input SearchUsersInput, requestorID ulid.ULID) ([]*usergw.UserView, error) {
	if input.Query == "" {
		return nil, xerrors.ErrInvalidArgument.WithMessage("query is required")
	}

	var cursor *time.Time
	if input.PageToken != "" {
		raw, err := base64.StdEncoding.DecodeString(input.PageToken)
		if err != nil {
			return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
		}
		t, err := time.Parse("2006-01-02T15:04:05.999999999Z", string(raw))
		if err != nil {
			return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
		}
		cursor = &t
	}

	limit := input.PageSize
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	users, err := uc.userGateway.SearchUsers(ctx, input.Query, cursor, limit)
	if err != nil {
		return nil, err
	}

	ids := make([]ulid.ULID, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}

	return uc.userGateway.BuildListUserViews(ctx, requestorID, ids)
}
