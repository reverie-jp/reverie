package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListPostLikes struct {
	postGateway postgw.Gateway
}

func NewListPostLikes(postGateway postgw.Gateway) *ListPostLikes {
	return &ListPostLikes{postGateway: postGateway}
}

func (uc *ListPostLikes) Execute(ctx context.Context, postIDStr string, requestorID ulid.ULID, pageSize int32) ([]*usergw.UserView, error) {
	postID, err := ulid.Parse(postIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid post_id")
	}

	return uc.postGateway.ListPostLikes(ctx, postID, requestorID, normalizePageSize(pageSize))
}
