package usecase

import notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"

type ListNotificationsOutput struct {
	Views         []*notificationgw.NotificationView
	NextPageToken string
}
