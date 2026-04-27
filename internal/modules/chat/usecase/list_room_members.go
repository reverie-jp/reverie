package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListRoomMembers struct {
	gateway chatgw.Gateway
}

func NewListRoomMembers(gateway chatgw.Gateway) *ListRoomMembers {
	return &ListRoomMembers{gateway: gateway}
}

func (uc *ListRoomMembers) Execute(ctx context.Context, roomIDStr string) ([]*ChatUserOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	members, err := uc.gateway.ListRoomMembers(ctx, roomID)
	if err != nil {
		return nil, err
	}
	outputs := make([]*ChatUserOutput, len(members))
	for i, m := range members {
		outputs[i] = entityUserToOutput(m)
	}
	return outputs, nil
}
