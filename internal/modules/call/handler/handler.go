package handler

import (
	"reverie.jp/reverie/internal/gen/pb/call/v1/callv1connect"
	"reverie.jp/reverie/internal/modules/call/usecase"
)

type Handler struct {
	callv1connect.UnimplementedCallServiceHandler
	joinCall *usecase.JoinCall
}

func New(joinCall *usecase.JoinCall) *Handler {
	return &Handler{
		joinCall: joinCall,
	}
}
