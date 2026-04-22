package entity

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type CallVisibility string

const (
	CallVisibilityOpen      CallVisibility = "open"
	CallVisibilityUsersOnly CallVisibility = "users_only"
	CallVisibilityLocked    CallVisibility = "locked"
)

type Call struct {
	ID         ulid.ULID
	HostUserID ulid.ULID
	Visibility CallVisibility
	EndTime    *time.Time
	CreateTime time.Time
	UpdateTime time.Time
}

// ParticipantStaleSeconds is the heartbeat grace window. A participant whose
// last_seen_time is older than this many seconds (and has no explicit
// disconnected_time) is considered disconnected. Clients heartbeat every 30s,
// so 60s leaves room for one missed heartbeat.
const ParticipantStaleSeconds = 60

type CallParticipant struct {
	CallID              ulid.ULID
	ParticipantIdentity string
	UserID              *ulid.ULID
	DisplayName         string
	FirstJoinTime       time.Time
	LastSeenTime        time.Time
	DisconnectedTime    *time.Time
	MutedByHost         bool
}

func (p *CallParticipant) IsCurrentlyConnected() bool {
	if p == nil || p.DisconnectedTime != nil {
		return false
	}
	return time.Since(p.LastSeenTime) < ParticipantStaleSeconds*time.Second
}

type CallBan struct {
	CallID     ulid.ULID
	UserID     ulid.ULID
	CreateTime time.Time
}
