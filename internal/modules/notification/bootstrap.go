package notification

import (
	"reverie.jp/reverie/internal/gen/pb/notification/v1/notificationv1connect"
	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	"reverie.jp/reverie/internal/modules/notification/handler"
	"reverie.jp/reverie/internal/modules/notification/usecase"
)

func InitModule(notificationGateway notificationgw.Gateway) notificationv1connect.NotificationServiceHandler {
	listNotifications := usecase.NewListNotifications(notificationGateway)
	markNotificationsRead := usecase.NewMarkNotificationsRead(notificationGateway)
	getUnreadNotificationCount := usecase.NewGetUnreadNotificationCount(notificationGateway)
	return handler.New(listNotifications, markNotificationsRead, getUnreadNotificationCount)
}
