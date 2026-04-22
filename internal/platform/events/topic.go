package events

import "reverie.jp/reverie/internal/platform/ulid"

// Topic namespace for real-time events. Prefixed so other Redis usages (cache
// etc.) don't collide.
const topicPrefix = "event:"

// UserTopic is the inbox channel for a specific user: notifications, DMs, etc.
func UserTopic(userID ulid.ULID) string {
	return topicPrefix + "user:" + userID.String()
}

// TimelineTopic receives call / post updates for a user's timeline. Fan-out
// is done at publish time (recipients are expanded from followers).
func TimelineTopic(userID ulid.ULID) string {
	return topicPrefix + "timeline:" + userID.String()
}
