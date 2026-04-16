package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) DeletePost(ctx context.Context, postID ulid.ULID, authorID ulid.ULID) error {
	return g.repo.DeletePost(ctx, postID, authorID)
}
