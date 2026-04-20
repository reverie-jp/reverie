package handler

import (
	"reverie.jp/reverie/internal/gen/pb/user/v1/userv1connect"
	"reverie.jp/reverie/internal/modules/user/usecase"
)

type Handler struct {
	userv1connect.UnimplementedUserServiceHandler
	getUser *usecase.GetUser
}

func New(getUser *usecase.GetUser) *Handler {
	return &Handler{getUser: getUser}
}
