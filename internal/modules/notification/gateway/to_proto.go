package gateway

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"reverie.jp/reverie/internal/domain/entity"
	notificationv1 "reverie.jp/reverie/internal/gen/pb/notification/v1"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	"reverie.jp/reverie/internal/platform/resourcename"
)

// ToProto converts a view into the proto Notification shared between the
// ListNotifications response and the event envelope payload. Kept on the
// gateway so both usages produce identical output.
func ToProto(view *NotificationView) *notificationv1.Notification {
	if view == nil || view.Notification == nil {
		return nil
	}
	n := view.Notification
	pb := &notificationv1.Notification{
		Type:         toProtoType(n.Type),
		ResourceName: n.ResourceName,
		CreateTime:   timestamppb.New(n.CreateTime),
	}
	if n.ReadTime != nil {
		pb.ReadTime = timestamppb.New(*n.ReadTime)
	}
	if view.Actor != nil && view.Actor.User != nil {
		pb.Name = resourcename.FormatNotification(view.Actor.User.CustomID, n.ID)
		pb.Actor = useradapter.ToUser(view.Actor)
	} else {
		pb.Name = "notifications/" + n.ID.String()
	}
	return pb
}

func toProtoType(t entity.NotificationType) notificationv1.NotificationType {
	switch t {
	case entity.NotificationTypeUserFollowed:
		return notificationv1.NotificationType_NOTIFICATION_TYPE_USER_FOLLOWED
	case entity.NotificationTypeFollowingUserCallStarted:
		return notificationv1.NotificationType_NOTIFICATION_TYPE_FOLLOWING_USER_CALL_STARTED
	default:
		return notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED
	}
}
