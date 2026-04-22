package event

import (
	"reverie.jp/reverie/internal/gen/pb/event/v1/eventv1connect"
	"reverie.jp/reverie/internal/modules/event/handler"
	"reverie.jp/reverie/internal/modules/event/usecase"
	"reverie.jp/reverie/internal/platform/events"
)

func InitModule(subscriber events.Subscriber) eventv1connect.EventServiceHandler {
	streamEvents := usecase.NewStreamEvents(subscriber)
	return handler.New(streamEvents)
}
