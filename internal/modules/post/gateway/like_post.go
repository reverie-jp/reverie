package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (g *gatewayImpl) LikePost(ctx context.Context, postID ulid.ULID, userID ulid.ULID) (*PostView, error) {
	post, err := g.repo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, xerrors.ErrNotFound
	}
	if err := g.repo.CreatePostFavorite(ctx, userID, postID); err != nil {
		return nil, err
	}
	views, err := g.buildViews(ctx, []*entity.Post{post}, userID)
	if err != nil {
		return nil, err
	}
	return views[0], nil
}

func (g *gatewayImpl) UnlikePost(ctx context.Context, postID ulid.ULID, userID ulid.ULID) (*PostView, error) {
	post, err := g.repo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, xerrors.ErrNotFound
	}
	if err := g.repo.DeletePostFavorite(ctx, userID, postID); err != nil {
		return nil, err
	}
	views, err := g.buildViews(ctx, []*entity.Post{post}, userID)
	if err != nil {
		return nil, err
	}
	return views[0], nil
}
