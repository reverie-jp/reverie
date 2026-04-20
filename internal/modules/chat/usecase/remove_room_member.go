package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type RemoveRoomMember struct {
	gateway chatgw.Gateway
}

func NewRemoveRoomMember(gateway chatgw.Gateway) *RemoveRoomMember {
	return &RemoveRoomMember{gateway: gateway}
}

func (uc *RemoveRoomMember) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr, userIDStr string) (*RoomOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.RemoveRoomMember(ctx, roomID, userID)
	if err != nil {
		return nil, err
	}
	return toRoomOutput(view, requestorID), nil
}
