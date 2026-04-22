package entity

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

// UserPresenceStaleSeconds is the heartbeat grace window. A user whose
// LastSeenTime is older than this many seconds is treated as offline.
// Clients heartbeat every 20s (see web/app/lib/use-presence-heartbeat.ts),
// so the 60s window tolerates two missed beats before flipping to offline.
const UserPresenceStaleSeconds = 60

type User struct {
	ID                 ulid.ULID
	CustomID           string
	CustomIDChangeTime *time.Time
	DisplayName        string
	Biography          *string
	Location           *string
	Website            *string
	AvatarURL          *string
	BannerURL          *string
	IsPrivate          bool
	FollowingCount     int32
	FollowerCount      int32
	LastSeenTime       *time.Time
	CreateTime         time.Time
	UpdateTime         time.Time
}

// IsCurrentlyOnline returns true when the user's last heartbeat landed within
// the presence stale window.
func (u *User) IsCurrentlyOnline() bool {
	if u == nil || u.LastSeenTime == nil {
		return false
	}
	return time.Since(*u.LastSeenTime) < UserPresenceStaleSeconds*time.Second
}

type AuthProvider struct {
	UserID ulid.ULID
}
