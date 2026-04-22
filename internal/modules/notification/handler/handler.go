package handler

import (
	"reverie.jp/reverie/internal/gen/pb/notification/v1/notificationv1connect"
	"reverie.jp/reverie/internal/modules/notification/usecase"
)

type Handler struct {
	notificationv1connect.UnimplementedNotificationServiceHandler
	listNotifications          *usecase.ListNotifications
	markNotificationsRead      *usecase.MarkNotificationsRead
	getUnreadNotificationCount *usecase.GetUnreadNotificationCount
}

func New(
	listNotifications *usecase.ListNotifications,
	markNotificationsRead *usecase.MarkNotificationsRead,
	getUnreadNotificationCount *usecase.GetUnreadNotificationCount,
) *Handler {
	return &Handler{
		listNotifications:          listNotifications,
		markNotificationsRead:      markNotificationsRead,
		getUnreadNotificationCount: getUnreadNotificationCount,
	}
}
