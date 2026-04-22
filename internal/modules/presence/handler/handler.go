package handler

import (
	"reverie.jp/reverie/internal/gen/pb/presence/v1/presencev1connect"
	"reverie.jp/reverie/internal/modules/presence/usecase"
)

type Handler struct {
	presencev1connect.UnimplementedPresenceServiceHandler
	heartbeat *usecase.Heartbeat
}

func New(heartbeat *usecase.Heartbeat) *Handler {
	return &Handler{heartbeat: heartbeat}
}
