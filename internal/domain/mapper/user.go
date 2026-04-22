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
		ID:                 row.ID,
		CustomID:           row.CustomID,
		CustomIDChangeTime: row.CustomIDChangeTime,
		DisplayName:        row.DisplayName,
		Biography:          row.Biography,
		Location:           row.Location,
		Website:            row.Website,
		AvatarURL:          row.AvatarUrl,
		BannerURL:          row.BannerUrl,
		IsPrivate:          row.IsPrivate,
		FollowingCount:     row.FollowingCount,
		FollowerCount:      row.FollowerCount,
		LastSeenTime:       row.LastSeenTime,
		CreateTime:         row.CreateTime,
		UpdateTime:         row.UpdateTime,
	}
}
