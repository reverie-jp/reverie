package usecase

import (
	"context"
	"errors"

	eventv1 "reverie.jp/reverie/internal/gen/pb/event/v1"
	"reverie.jp/reverie/internal/platform/events"
	"reverie.jp/reverie/internal/platform/ulid"
)

type StreamEvents struct {
	subscriber events.Subscriber
}

func NewStreamEvents(subscriber events.Subscriber) *StreamEvents {
	return &StreamEvents{subscriber: subscriber}
}

// Subscribe opens an events subscription for the authenticated user and
// returns the channel. The caller (handler) is responsible for forwarding
// envelopes to the Connect server stream and terminating on ctx cancellation.
func (uc *StreamEvents) Subscribe(ctx context.Context, userID ulid.ULID) (<-chan *eventv1.StreamEventsResponse, error) {
	if userID.IsZero() {
		return nil, errors.New("user_id is required")
	}
	return uc.subscriber.Subscribe(ctx, events.UserTopic(userID))
}
