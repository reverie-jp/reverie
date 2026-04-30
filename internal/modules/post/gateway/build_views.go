package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/modules/post/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (g *gatewayImpl) ListTimeline(ctx context.Context, params ListTimelineParams, requestorID ulid.ULID) ([]*PostView, error) {
	posts, err := g.repo.ListTimeline(ctx, repository.ListTimelineParams{
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		return nil, err
	}
	return g.buildViews(ctx, posts, requestorID)
}

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

func (g *gatewayImpl) ListUserPosts(ctx context.Context, params ListUserPostsParams, requestorID ulid.ULID) ([]*PostView, error) {
	posts, err := g.repo.ListUserPosts(ctx, repository.ListUserPostsParams{
		AuthorID: params.AuthorID,
		Cursor:   params.Cursor,
		Limit:    params.Limit,
	})
	if err != nil {
		return nil, err
	}
	return g.buildViews(ctx, posts, requestorID)
}

func (g *gatewayImpl) ListPostReplies(ctx context.Context, authorCustomID string, shortID string, params ListPostRepliesParams, requestorID ulid.ULID) ([]*PostView, error) {
	post, err := g.repo.GetPostByShortID(ctx, authorCustomID, shortID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, xerrors.ErrNotFound
	}
	posts, err := g.repo.ListPostReplies(ctx, repository.ListPostRepliesParams{
		PostID: post.ID,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		return nil, err
	}
	return g.buildViews(ctx, posts, requestorID)
}

func (g *gatewayImpl) ListPostReposts(ctx context.Context, authorCustomID string, shortID string, params ListPostRepostsParams, requestorID ulid.ULID) ([]*PostView, error) {
	post, err := g.repo.GetPostByShortID(ctx, authorCustomID, shortID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, xerrors.ErrNotFound
	}
	posts, err := g.repo.ListPostReposts(ctx, repository.ListPostRepostsParams{
		PostID: post.ID,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		return nil, err
	}
	return g.buildViews(ctx, posts, requestorID)
}

func (g *gatewayImpl) ListUserLikedPosts(ctx context.Context, params ListUserLikedPostsParams, requestorID ulid.ULID) ([]*PostView, error) {
	posts, err := g.repo.ListUserLikedPosts(ctx, repository.ListUserLikedPostsParams{
		UserID: params.UserID,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		return nil, err
	}
	return g.buildViews(ctx, posts, requestorID)
}

func (g *gatewayImpl) ListPostLikes(ctx context.Context, authorCustomID string, shortID string, requestorID ulid.ULID, limit int32) ([]*usergw.UserView, error) {
	post, err := g.repo.GetPostByShortID(ctx, authorCustomID, shortID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, xerrors.ErrNotFound
	}
	users, err := g.repo.ListPostLikes(ctx, post.ID, limit)
	if err != nil {
		return nil, err
	}
	ids := make([]ulid.ULID, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}
	return g.userGateway.BuildListUserViews(ctx, requestorID, ids)
}

// buildViews composes PostView slices from post entities.
// Author is enriched via userGateway.BuildListUserViews which resolves
// follow relationships (IsFollowing / IsFollowedBy / IsMe) in a single query.
func (g *gatewayImpl) buildViews(ctx context.Context, posts []*entity.Post, requestorID ulid.ULID) ([]*PostView, error) {
	if len(posts) == 0 {
		return []*PostView{}, nil
	}

	authorIDs := make([]ulid.ULID, 0, len(posts))
	seen := map[ulid.ULID]bool{}
	repostOriginals := map[ulid.ULID]*entity.Post{}

	for _, p := range posts {
		if !seen[p.AuthorID] {
			seen[p.AuthorID] = true
			authorIDs = append(authorIDs, p.AuthorID)
		}
		if p.RepostPostID != nil {
			if _, ok := repostOriginals[*p.RepostPostID]; !ok {
				orig, err := g.repo.GetPostByID(ctx, *p.RepostPostID)
				if err != nil {
					return nil, err
				}
				if orig != nil {
					repostOriginals[*p.RepostPostID] = orig
					if !seen[orig.AuthorID] {
						seen[orig.AuthorID] = true
						authorIDs = append(authorIDs, orig.AuthorID)
					}
				}
			}
		}
	}

	// BuildListUserViews resolves follow relationships relative to requestorID.
	userViews, err := g.userGateway.BuildListUserViews(ctx, requestorID, authorIDs)
	if err != nil {
		return nil, err
	}
	viewByID := map[ulid.ULID]*usergw.UserView{}
	for _, uv := range userViews {
		if uv.User != nil {
			viewByID[uv.User.ID] = uv
		}
	}

	buildOne := func(p *entity.Post) (*PostView, error) {
		replyCount, err := g.repo.CountPostReplies(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		repostCount, err := g.repo.CountPostReposts(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		favCount, err := g.repo.CountPostFavorites(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		isFavorited, err := g.repo.GetPostFavorite(ctx, requestorID, p.ID)
		if err != nil {
			return nil, err
		}
		return &PostView{
			Post:          p,
			Author:        viewByID[p.AuthorID],
			ReplyCount:    replyCount,
			RepostCount:   repostCount,
			FavoriteCount: favCount,
			IsFavorited:   isFavorited,
		}, nil
	}

	views := make([]*PostView, len(posts))
	for i, p := range posts {
		v, err := buildOne(p)
		if err != nil {
			return nil, err
		}
		if p.RepostPostID != nil {
			if orig, ok := repostOriginals[*p.RepostPostID]; ok {
				v.RepostOf, err = buildOne(orig)
				if err != nil {
					return nil, err
				}
			}
		}
		views[i] = v
	}

	return views, nil
}
