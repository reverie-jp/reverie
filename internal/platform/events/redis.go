package events

import (
	"context"
	"errors"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"

	eventv1 "reverie.jp/reverie/internal/gen/pb/event/v1"
)

// redisBus is a Redis pub/sub backed Bus. Payloads are protobuf-encoded
// StreamEventsResponses. Pub/sub in Redis is ephemeral — the DB remains the source
// of truth, so lost messages only delay client UI refresh.
type redisBus struct {
	client *redis.Client
}

func NewRedisBus(client *redis.Client) Bus {
	return &redisBus{client: client}
}

func (r *redisBus) Publish(ctx context.Context, topic string, env *eventv1.StreamEventsResponse) error {
	data, err := proto.Marshal(env)
	if err != nil {
		return err
	}
	return r.client.Publish(ctx, topic, data).Err()
}

func (r *redisBus) Subscribe(ctx context.Context, topics ...string) (<-chan *eventv1.StreamEventsResponse, error) {
	if len(topics) == 0 {
		return nil, errors.New("events: Subscribe requires at least one topic")
	}
	ps := r.client.Subscribe(ctx, topics...)
	// Block until SUBSCRIBE is acknowledged so callers know the subscription
	// is live before they return tokens downstream.
	if _, err := ps.Receive(ctx); err != nil {
		_ = ps.Close()
		return nil, err
	}

	out := make(chan *eventv1.StreamEventsResponse, 16)
	go func() {
		defer close(out)
		defer ps.Close()
		msgCh := ps.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				env := &eventv1.StreamEventsResponse{}
				if err := proto.Unmarshal([]byte(msg.Payload), env); err != nil {
					slog.Warn("events: failed to unmarshal envelope",
						slog.String("topic", msg.Channel),
						slog.String("err", err.Error()),
					)
					continue
				}
				select {
				case out <- env:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out, nil
}
