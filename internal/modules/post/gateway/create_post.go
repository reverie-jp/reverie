package gateway

import (
	"context"
	"strings"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/modules/post/repository"
	"reverie.jp/reverie/internal/platform/shortid"
)

func (g *gatewayImpl) CreatePost(ctx context.Context, params CreatePostParams) (*PostView, error) {
	post, err := g.createWithShortID(ctx, params)
	if err != nil {
		return nil, err
	}

	author, err := g.userGateway.BuildUserView(ctx, params.AuthorID, params.AuthorID)
	if err != nil {
		return nil, err
	}

	return &PostView{
		Post:   post,
		Author: author,
	}, nil
}

func (g *gatewayImpl) createWithShortID(ctx context.Context, params CreatePostParams) (*entity.Post, error) {
	for {
		sid, err := shortid.Generate()
		if err != nil {
			return nil, err
		}
		post, err := g.repo.CreatePost(ctx, repository.CreatePostParams{
			ID:        params.ID,
			AuthorID:  params.AuthorID,
			ShortID:   sid,
			ReplyToPostID: params.ReplyToPostID,
			RepostPostID:  params.RepostPostID,
			Text:      params.Text,
		})
		if err != nil {
			if isUniqueViolation(err) {
				continue
			}
			return nil, err
		}
		return post, nil
	}
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "23505")
}
