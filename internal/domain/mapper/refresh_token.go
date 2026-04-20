package mapper

import (
	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func ToRefreshToken(row *sqlc.RefreshToken) *entity.RefreshToken {
	if row == nil {
		return nil
	}

	return &entity.RefreshToken{
		ID:         row.ID,
		UserID:     row.UserID,
		TokenHash:  row.TokenHash,
		ExpireTime: row.ExpireTime,
		CreateTime: row.CreateTime,
	}
}
