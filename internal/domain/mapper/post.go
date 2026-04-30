package mapper

import (
	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func ToPost(row *sqlc.Post) *entity.Post {
	if row == nil {
		return nil
	}
	return &entity.Post{
		ID:         row.ID,
		AuthorID:   row.AuthorID,
		ShortID:    row.ShortID,
		ReplyToPostID:  row.ReplyToPostID,
		RepostPostID:   row.RepostPostID,
		Text:       row.Text,
		CreateTime: row.CreateTime,
		UpdateTime: row.UpdateTime,
	}
}
