package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type CreateDirectRoom struct {
	gateway chatgw.Gateway
}

func NewCreateDirectRoom(gateway chatgw.Gateway) *CreateDirectRoom {
	return &CreateDirectRoom{gateway: gateway}
}

func (uc *CreateDirectRoom) Execute(ctx context.Context, myID ulid.ULID, otherIDStr string) (*RoomOutput, error) {
	otherID, err := ulid.Parse(otherIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	if myID == otherID {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.GetOrCreateDirectRoom(ctx, myID, otherID)
	if err != nil {
		return nil, err
	}
	return toRoomOutput(view, myID), nil
}
