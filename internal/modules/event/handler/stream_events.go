package handler

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"reverie.jp/reverie/internal/application/server/interceptor"
	eventv1 "reverie.jp/reverie/internal/gen/pb/event/v1"
	"reverie.jp/reverie/internal/platform/xerrors"
)

// keepAliveInterval is how often the server sends a KeepAlive envelope when
// the subscriber's topic is idle. Short enough to beat typical LB / browser
// chunked-response idle timeouts (commonly 30–60s), long enough to not be
// chatty over mobile networks.
const keepAliveInterval = 15 * time.Second

func (h *Handler) StreamEvents(ctx context.Context, _ *connect.Request[eventv1.StreamEventsRequest], stream *connect.ServerStream[eventv1.StreamEventsResponse]) error {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return xerrors.ErrUnauthenticated
	}

	ch, err := h.streamEvents.Subscribe(ctx, userID)
	if err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}

	slog.Info("stream opened", slog.String("user", userID.String()))
	defer slog.Info("stream closed", slog.String("user", userID.String()))

	sendKeepAlive := func(t time.Time) error {
		return stream.Send(&eventv1.StreamEventsResponse{
			CreateTime: timestamppb.New(t),
			Payload: &eventv1.StreamEventsResponse_KeepAlive{
				KeepAlive: &eventv1.KeepAliveEvent{},
			},
		})
	}

	// Fire one frame immediately so the client's fetch flips out of
	// "pending" as soon as the handler is reached (nicer DevTools UX, and
	// confirms for the client that auth/transport are healthy).
	if err := sendKeepAlive(time.Now()); err != nil {
		slog.Warn("stream initial keepalive send failed", slog.String("err", err.Error()))
		return err
	}

	ticker := time.NewTicker(keepAliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case env, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(env); err != nil {
				slog.Warn("stream send failed", slog.String("err", err.Error()))
				return err
			}
		case t := <-ticker.C:
			slog.Debug("stream keepalive", slog.String("user", userID.String()))
			if err := sendKeepAlive(t); err != nil {
				slog.Warn("stream keepalive send failed", slog.String("err", err.Error()))
				return err
			}
		}
	}
}
