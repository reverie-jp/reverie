package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type CreateGroupRoom struct {
	gateway chatgw.Gateway
}

func NewCreateGroupRoom(gateway chatgw.Gateway) *CreateGroupRoom {
	return &CreateGroupRoom{gateway: gateway}
}

func (uc *CreateGroupRoom) Execute(ctx context.Context, creatorID ulid.ULID, name string, memberIDStrs []string) (*RoomOutput, error) {
	if name == "" {
		return nil, xerrors.ErrInvalidArgument
	}
	memberIDs := make([]ulid.ULID, 0, len(memberIDStrs))
	for _, s := range memberIDStrs {
		id, err := ulid.Parse(s)
		if err != nil {
			return nil, xerrors.ErrInvalidArgument
		}
		memberIDs = append(memberIDs, id)
	}
	view, err := uc.gateway.CreateGroupRoom(ctx, creatorID, name, memberIDs)
	if err != nil {
		return nil, err
	}
	return toRoomOutput(view, creatorID), nil
}
