package mapper

import (
	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func ToUser(row *sqlc.User) *entity.User {
	if row == nil {
		return nil
	}

	return &entity.User{
		ID:                row.ID,
		CustomID:          row.CustomID,
		CustomIDChangedAt: row.CustomIDChangedAt,
		DisplayName:       row.DisplayName,
		Biography:         row.Biography,
		AvatarMediaID:     row.AvatarMediaID,
		BannerMediaID:     row.BannerMediaID,
		IsPrivate:         row.IsPrivate,
		Birthdate:         row.Birthdate,
		CreateTime:        row.CreateTime,
		UpdateTime:        row.UpdateTime,
	}
}
