package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/modules/post/repository"
	"reverie.jp/reverie/internal/platform/ulid"
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

func (g *gatewayImpl) buildViews(ctx context.Context, posts []*entity.Post, requestorID ulid.ULID) ([]*PostView, error) {
	if len(posts) == 0 {
		return []*PostView{}, nil
	}

	// collect all author IDs (including repost originals)
	authorIDs := make([]ulid.ULID, 0, len(posts))
	seen := map[ulid.ULID]bool{}

	// fetch repost originals
	repostOriginals := map[ulid.ULID]*entity.Post{}
	for _, p := range posts {
		if !seen[p.AuthorID] {
			seen[p.AuthorID] = true
			authorIDs = append(authorIDs, p.AuthorID)
		}
		if p.RepostID != nil {
			if _, ok := repostOriginals[*p.RepostID]; !ok {
				orig, err := g.repo.GetPostByID(ctx, *p.RepostID)
				if err != nil {
					return nil, err
				}
				if orig != nil {
					repostOriginals[*p.RepostID] = orig
					if !seen[orig.AuthorID] {
						seen[orig.AuthorID] = true
						authorIDs = append(authorIDs, orig.AuthorID)
					}
				}
			}
		}
	}

	userViews, err := g.userGateway.BuildListViews(ctx, authorIDs)
	if err != nil {
		return nil, err
	}
	userMap := map[ulid.ULID]*entity.User{}
	for _, uv := range userViews {
		userMap[uv.User.ID] = uv.User
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
			Author:        userMap[p.AuthorID],
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
		if p.RepostID != nil {
			if orig, ok := repostOriginals[*p.RepostID]; ok {
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
