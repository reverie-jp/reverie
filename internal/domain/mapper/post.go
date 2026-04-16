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
		ReplyToID:  row.ReplyToID,
		RepostID:   row.RepostID,
		Text:       row.Text,
		CreateTime: row.CreateTime,
		UpdateTime: row.UpdateTime,
	}
}
