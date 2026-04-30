package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (g *gatewayImpl) GetPost(ctx context.Context, authorCustomID string, shortID string, requestorID ulid.ULID) (*PostView, error) {
	post, err := g.repo.GetPostByShortID(ctx, authorCustomID, shortID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, nil
	}

	views, err := g.buildViews(ctx, []*entity.Post{post}, requestorID)
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, xerrors.ErrNotFound
	}
	return views[0], nil
}
