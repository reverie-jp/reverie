package adapter

import (
	"reverie.jp/reverie/internal/domain/entity"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
)

func toProtoVisibility(v entity.CallVisibility) callv1.CallVisibility {
	switch v {
	case entity.CallVisibilityOpen:
		return callv1.CallVisibility_CALL_VISIBILITY_OPEN
	case entity.CallVisibilityUsersOnly:
		return callv1.CallVisibility_CALL_VISIBILITY_USERS_ONLY
	case entity.CallVisibilityLocked:
		return callv1.CallVisibility_CALL_VISIBILITY_LOCKED
	default:
		return callv1.CallVisibility_CALL_VISIBILITY_UNSPECIFIED
	}
}

func fromProtoVisibility(v callv1.CallVisibility) entity.CallVisibility {
	switch v {
	case callv1.CallVisibility_CALL_VISIBILITY_OPEN:
		return entity.CallVisibilityOpen
	case callv1.CallVisibility_CALL_VISIBILITY_USERS_ONLY:
		return entity.CallVisibilityUsersOnly
	case callv1.CallVisibility_CALL_VISIBILITY_LOCKED:
		return entity.CallVisibilityLocked
	default:
		return ""
	}
}
