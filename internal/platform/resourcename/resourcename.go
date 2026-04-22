// Package resourcename translates between internal identifiers and
// AIP-122-style resource names ("collection/id" hierarchical strings).
package resourcename

import (
	"errors"
	"strings"

	"reverie.jp/reverie/internal/platform/ulid"
)

const (
	usersCollection         = "users"
	callsCollection         = "calls"
	participantsCollection  = "participants"
	bansCollection          = "bans"
	notificationsCollection = "notifications"
)

// User

func FormatUser(customID string) string {
	return usersCollection + "/" + customID
}

func ParseUser(name string) (string, error) {
	segments, err := split(name)
	if err != nil {
		return "", err
	}
	if len(segments) != 2 || segments[0] != usersCollection || segments[1] == "" {
		return "", errors.New("invalid user resource name: " + name)
	}
	return segments[1], nil
}

// Call

func FormatCall(id ulid.ULID) string {
	return callsCollection + "/" + id.String()
}

func ParseCall(name string) (ulid.ULID, error) {
	segments, err := split(name)
	if err != nil {
		return ulid.ULID{}, err
	}
	if len(segments) != 2 || segments[0] != callsCollection || segments[1] == "" {
		return ulid.ULID{}, errors.New("invalid call resource name: " + name)
	}
	return ulid.Parse(segments[1])
}

// CallParticipant

func FormatCallParticipant(callID ulid.ULID, identity string) string {
	return FormatCall(callID) + "/" + participantsCollection + "/" + identity
}

func ParseCallParticipant(name string) (callID ulid.ULID, identity string, err error) {
	segments, err := split(name)
	if err != nil {
		return ulid.ULID{}, "", err
	}
	if len(segments) != 4 ||
		segments[0] != callsCollection ||
		segments[2] != participantsCollection ||
		segments[1] == "" || segments[3] == "" {
		return ulid.ULID{}, "", errors.New("invalid call participant resource name: " + name)
	}
	callID, err = ulid.Parse(segments[1])
	if err != nil {
		return ulid.ULID{}, "", err
	}
	return callID, segments[3], nil
}

// CallBan

func FormatCallBan(callID, userID ulid.ULID) string {
	return FormatCall(callID) + "/" + bansCollection + "/" + userID.String()
}

func ParseCallBan(name string) (callID, userID ulid.ULID, err error) {
	segments, err := split(name)
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	if len(segments) != 4 ||
		segments[0] != callsCollection ||
		segments[2] != bansCollection ||
		segments[1] == "" || segments[3] == "" {
		return ulid.ULID{}, ulid.ULID{}, errors.New("invalid call ban resource name: " + name)
	}
	callID, err = ulid.Parse(segments[1])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	userID, err = ulid.Parse(segments[3])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	return callID, userID, nil
}

// Notification

// FormatNotification builds the resource name for a notification scoped to a
// user: "users/{custom_id}/notifications/{id}".
func FormatNotification(userCustomID string, notificationID ulid.ULID) string {
	return FormatUser(userCustomID) + "/" + notificationsCollection + "/" + notificationID.String()
}

// ParseNotification extracts the owning user's custom_id and the notification
// ULID from a resource name.
func ParseNotification(name string) (userCustomID string, notificationID ulid.ULID, err error) {
	segments, err := split(name)
	if err != nil {
		return "", ulid.ULID{}, err
	}
	if len(segments) != 4 ||
		segments[0] != usersCollection ||
		segments[2] != notificationsCollection ||
		segments[1] == "" || segments[3] == "" {
		return "", ulid.ULID{}, errors.New("invalid notification resource name: " + name)
	}
	notificationID, err = ulid.Parse(segments[3])
	if err != nil {
		return "", ulid.ULID{}, err
	}
	return segments[1], notificationID, nil
}

func split(name string) ([]string, error) {
	if name == "" {
		return nil, errors.New("empty resource name")
	}
	return strings.Split(name, "/"), nil
}
