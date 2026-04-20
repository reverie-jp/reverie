package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type UpdateRoom struct {
	gateway chatgw.Gateway
}

func NewUpdateRoom(gateway chatgw.Gateway) *UpdateRoom {
	return &UpdateRoom{gateway: gateway}
}

func (uc *UpdateRoom) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr, name string) (*RoomOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.UpdateRoom(ctx, roomID, name)
	if err != nil {
		return nil, err
	}
	return toRoomOutput(view, requestorID), nil
}
