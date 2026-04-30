package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (g *gatewayImpl) DeletePost(ctx context.Context, authorCustomID string, shortID string, authorID ulid.ULID) error {
	post, err := g.repo.GetPostByShortID(ctx, authorCustomID, shortID)
	if err != nil {
		return err
	}
	if post == nil {
		return xerrors.ErrNotFound
	}
	return g.repo.DeletePost(ctx, post.ID, authorID)
}
