package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type ListRooms struct {
	gateway chatgw.Gateway
}

func NewListRooms(gateway chatgw.Gateway) *ListRooms {
	return &ListRooms{gateway: gateway}
}

func (uc *ListRooms) Execute(ctx context.Context, requestorID ulid.ULID) ([]*RoomOutput, error) {
	views, err := uc.gateway.ListRooms(ctx, requestorID)
	if err != nil {
		return nil, err
	}
	out := make([]*RoomOutput, len(views))
	for i, v := range views {
		out[i] = toRoomOutput(v, requestorID)
	}
	return out, nil
}
