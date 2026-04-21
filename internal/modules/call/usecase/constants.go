package usecase

// participantStaleSeconds is the heartbeat grace window. A participant whose
// last_seen_time is older than this many seconds (and has no explicit
// disconnected_time) is considered disconnected. Clients heartbeat every
// 30s, so 60s leaves room for one missed heartbeat.
const participantStaleSeconds = 60
