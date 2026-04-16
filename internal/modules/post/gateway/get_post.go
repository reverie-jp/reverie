package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) GetPost(ctx context.Context, postID ulid.ULID, requestorID ulid.ULID) (*PostView, error) {
	post, err := g.repo.GetPostByID(ctx, postID)
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
		return nil, nil
	}
	return views[0], nil
}
