package entity

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type NotificationType string

const (
	NotificationTypeUserFollowed              NotificationType = "user_followed"
	NotificationTypeFollowingUserCallStarted  NotificationType = "following_user_call_started"
)

type Notification struct {
	ID              ulid.ULID
	RecipientUserID ulid.ULID
	Type            NotificationType
	ActorUserID     *ulid.ULID
	ResourceName    string
	ReadTime        *time.Time
	CreateTime      time.Time
}
