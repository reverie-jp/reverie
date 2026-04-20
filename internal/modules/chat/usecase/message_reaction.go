package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type AddMessageReaction struct {
	gateway chatgw.Gateway
}

func NewAddMessageReaction(gateway chatgw.Gateway) *AddMessageReaction {
	return &AddMessageReaction{gateway: gateway}
}

func (uc *AddMessageReaction) Execute(ctx context.Context, requestorID ulid.ULID, messageIDStr, emoji string) (*MessageOutput, error) {
	if emoji == "" {
		emoji = "❤️"
	}
	messageID, err := ulid.Parse(messageIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.AddReaction(ctx, messageID, requestorID, emoji)
	if err != nil {
		return nil, err
	}
	return toMessageOutput(view, requestorID), nil
}

type RemoveMessageReaction struct {
	gateway chatgw.Gateway
}

func NewRemoveMessageReaction(gateway chatgw.Gateway) *RemoveMessageReaction {
	return &RemoveMessageReaction{gateway: gateway}
}

func (uc *RemoveMessageReaction) Execute(ctx context.Context, requestorID ulid.ULID, messageIDStr, emoji string) (*MessageOutput, error) {
	if emoji == "" {
		emoji = "❤️"
	}
	messageID, err := ulid.Parse(messageIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.RemoveReaction(ctx, messageID, requestorID, emoji)
	if err != nil {
		return nil, err
	}
	return toMessageOutput(view, requestorID), nil
}
