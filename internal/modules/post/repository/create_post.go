package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func (r *repositoryImpl) CreatePost(ctx context.Context, params CreatePostParams) (*entity.Post, error) {
	row, err := r.q.CreatePost(ctx, sqlc.CreatePostParams{
		ID:        params.ID,
		AuthorID:  params.AuthorID,
		ShortID:   params.ShortID,
		ReplyToPostID: params.ReplyToPostID,
		RepostPostID:  params.RepostPostID,
		Text:      params.Text,
	})
	if err != nil {
		return nil, err
	}
	return mapper.ToPost(&row), nil
}
