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
	CreateTime time.Time
	UpdateTime time.Time
}

type CallParticipant struct {
	CallID              ulid.ULID
	ParticipantIdentity string
	UserID              *ulid.ULID
	DisplayName         string
	FirstJoinTime       time.Time
	LastSeenTime        time.Time
	DisconnectedTime    *time.Time
}
