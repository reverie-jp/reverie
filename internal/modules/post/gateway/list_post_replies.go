package gateway

import (
	"context"

	"reverie.jp/reverie/internal/modules/post/repository"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListPostReplies(ctx context.Context, params ListPostRepliesParams, requestorID ulid.ULID) ([]*PostView, error) {
	posts, err := g.repo.ListPostReplies(ctx, repository.ListPostRepliesParams{
		PostID: params.PostID,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		return nil, err
	}
	return g.buildViews(ctx, posts, requestorID)
}
