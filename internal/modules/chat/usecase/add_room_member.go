package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type AddRoomMember struct {
	gateway chatgw.Gateway
}

func NewAddRoomMember(gateway chatgw.Gateway) *AddRoomMember {
	return &AddRoomMember{gateway: gateway}
}

func (uc *AddRoomMember) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr, userIDStr string) (*RoomOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.AddRoomMember(ctx, roomID, userID)
	if err != nil {
		return nil, err
	}
	return toRoomOutput(view, requestorID), nil
}
