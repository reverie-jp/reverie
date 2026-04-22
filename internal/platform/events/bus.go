package events

import (
	"context"

	eventv1 "reverie.jp/reverie/internal/gen/pb/event/v1"
)

// Publisher fans out a single StreamEventsResponse to everyone listening on the given
// topic. Implementations are expected to be fire-and-forget: transport errors
// are returned but the caller usually only logs. The DB is the source of
// truth; stream is a hint.
type Publisher interface {
	Publish(ctx context.Context, topic string, env *eventv1.StreamEventsResponse) error
}

// Subscriber establishes a subscription for the given topics and returns a
// channel that yields StreamEventsResponse messages. Close the returned context (or
// cancel the parent) to stop the subscription and release the underlying
// connection. The channel is closed when the subscription terminates.
type Subscriber interface {
	Subscribe(ctx context.Context, topics ...string) (<-chan *eventv1.StreamEventsResponse, error)
}

// Bus is the combined Publisher + Subscriber used by bootstrap wiring.
type Bus interface {
	Publisher
	Subscriber
}
