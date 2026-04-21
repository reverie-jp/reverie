// Package resourcename translates between internal identifiers and
// AIP-122-style resource names ("collection/id" hierarchical strings).
package resourcename

import (
	"errors"
	"strings"

	"reverie.jp/reverie/internal/platform/ulid"
)

const (
	usersCollection        = "users"
	callsCollection        = "calls"
	participantsCollection = "participants"
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

func split(name string) ([]string, error) {
	if name == "" {
		return nil, errors.New("empty resource name")
	}
	return strings.Split(name, "/"), nil
}
