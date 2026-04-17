package gateway

import (
	"context"

	"reverie.jp/reverie/internal/modules/post/repository"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListPostReposts(ctx context.Context, params ListPostRepostsParams, requestorID ulid.ULID) ([]*PostView, error) {
	posts, err := g.repo.ListPostReposts(ctx, repository.ListPostRepostsParams{
		PostID: params.PostID,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		return nil, err
	}
	return g.buildViews(ctx, posts, requestorID)
}
