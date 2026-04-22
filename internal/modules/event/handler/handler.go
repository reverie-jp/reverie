package handler

import (
	"reverie.jp/reverie/internal/gen/pb/event/v1/eventv1connect"
	"reverie.jp/reverie/internal/modules/event/usecase"
)

type Handler struct {
	eventv1connect.UnimplementedEventServiceHandler
	streamEvents *usecase.StreamEvents
}

func New(streamEvents *usecase.StreamEvents) *Handler {
	return &Handler{streamEvents: streamEvents}
}
