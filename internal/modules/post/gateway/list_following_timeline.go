package gateway

import (
	"context"

	"reverie.jp/reverie/internal/modules/post/repository"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListFollowingTimeline(ctx context.Context, params ListFollowingTimelineParams, requestorID ulid.ULID) ([]*PostView, error) {
	posts, err := g.repo.ListFollowingTimeline(ctx, repository.ListFollowingTimelineParams{
		FollowerID: params.FollowerID,
		Cursor:     params.Cursor,
		Limit:      params.Limit,
	})
	if err != nil {
		return nil, err
	}
	return g.buildViews(ctx, posts, requestorID)
}
