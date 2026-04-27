package gateway

import (
	"context"

	"reverie.jp/reverie/internal/modules/post/repository"
)

func (g *gatewayImpl) CreatePost(ctx context.Context, params CreatePostParams) (*PostView, error) {
	post, err := g.repo.CreatePost(ctx, repository.CreatePostParams{
		ID:        params.ID,
		AuthorID:  params.AuthorID,
		ReplyToID: params.ReplyToID,
		RepostID:  params.RepostID,
		Text:      params.Text,
	})
	if err != nil {
		return nil, err
	}
	author, err := g.userGateway.BuildUserView(ctx, params.AuthorID, params.AuthorID)
	if err != nil {
		return nil, err
	}
	return &PostView{Post: post, Author: author}, nil
}
