package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type SendMessage struct {
	gateway chatgw.Gateway
}

func NewSendMessage(gateway chatgw.Gateway) *SendMessage {
	return &SendMessage{gateway: gateway}
}

func (uc *SendMessage) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr, content string) (*MessageOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	if content == "" {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.SendMessage(ctx, roomID, requestorID, content)
	if err != nil {
		return nil, err
	}
	return toMessageOutput(view, requestorID), nil
}
