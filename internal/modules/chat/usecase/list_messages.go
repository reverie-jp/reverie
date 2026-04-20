package usecase

import (
	"context"
	"encoding/base64"
	"time"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListMessages struct {
	gateway chatgw.Gateway
}

func NewListMessages(gateway chatgw.Gateway) *ListMessages {
	return &ListMessages{gateway: gateway}
}

func (uc *ListMessages) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr, pageToken string, pageSize int32) ([]*MessageOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}

	limit := pageSize
	if limit <= 0 || limit > 50 {
		limit = 30
	}

	var cursor *time.Time
	if pageToken != "" {
		raw, err := base64.StdEncoding.DecodeString(pageToken)
		if err == nil {
			t, err := time.Parse("2006-01-02T15:04:05.999999999Z", string(raw))
			if err == nil {
				cursor = &t
			}
		}
	}

	views, err := uc.gateway.ListMessages(ctx, roomID, cursor, limit, requestorID)
	if err != nil {
		return nil, err
	}
	out := make([]*MessageOutput, len(views))
	for i, v := range views {
		out[i] = toMessageOutput(v, requestorID)
	}
	return out, nil
}
