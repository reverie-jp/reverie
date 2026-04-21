package mapper

import (
	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func ToCall(row *sqlc.Call) *entity.Call {
	if row == nil {
		return nil
	}
	return &entity.Call{
		ID:         row.ID,
		HostUserID: row.HostUserID,
		Visibility: entity.CallVisibility(row.Visibility),
		CreateTime: row.CreateTime,
		UpdateTime: row.UpdateTime,
	}
}

func ToCallParticipant(row *sqlc.CallParticipant) *entity.CallParticipant {
	if row == nil {
		return nil
	}
	return &entity.CallParticipant{
		CallID:              row.CallID,
		ParticipantIdentity: row.ParticipantIdentity,
		UserID:              row.UserID,
		DisplayName:         row.DisplayName,
		FirstJoinTime:       row.FirstJoinTime,
		LastSeenTime:        row.LastSeenTime,
		DisconnectedTime:    row.DisconnectedTime,
	}
}
