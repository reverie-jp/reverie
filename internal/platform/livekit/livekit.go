package livekit

import (
	"context"
	"time"

	"github.com/livekit/protocol/auth"
	livekitproto "github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type Client struct {
	url       string
	apiKey    string
	apiSecret string
	roomSvc   *lksdk.RoomServiceClient
}

func NewClient(url, apiKey, apiSecret string) *Client {
	return &Client{
		url:       url,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		roomSvc:   lksdk.NewRoomServiceClient(toHTTPURL(url), apiKey, apiSecret),
	}
}

func (c *Client) URL() string {
	return c.url
}

type JoinTokenParams struct {
	RoomID      string
	Identity    string
	DisplayName string
	TTL         time.Duration
}

func (c *Client) CreateJoinToken(params JoinTokenParams) (string, error) {
	at := auth.NewAccessToken(c.apiKey, c.apiSecret)
	canPublish := true
	canSubscribe := true
	canPublishData := true
	// LiveKit's "VideoGrant" covers all real-time media permissions. Restricting
	// sources to "microphone" keeps this an audio-only call — clients cannot
	// publish camera or screen share even if they try.
	grant := &auth.VideoGrant{
		RoomJoin:          true,
		Room:              params.RoomID,
		CanPublish:        &canPublish,
		CanSubscribe:      &canSubscribe,
		CanPublishData:    &canPublishData,
		CanPublishSources: []string{"microphone"},
	}
	at.SetVideoGrant(grant).
		SetIdentity(params.Identity).
		SetName(params.DisplayName).
		SetValidFor(params.TTL)
	return at.ToJWT()
}

// MuteParticipantMicrophone server-side mutes (or unmutes) every audio track
// published by the given identity. Missing participants / tracks are ignored so
// callers can issue the request even if the participant just disconnected.
func (c *Client) MuteParticipantMicrophone(ctx context.Context, roomID, identity string, muted bool) error {
	info, err := c.roomSvc.GetParticipant(ctx, &livekitproto.RoomParticipantIdentity{
		Room:     roomID,
		Identity: identity,
	})
	if err != nil {
		return err
	}
	for _, track := range info.Tracks {
		if track.Type != livekitproto.TrackType_AUDIO {
			continue
		}
		if _, err := c.roomSvc.MutePublishedTrack(ctx, &livekitproto.MuteRoomTrackRequest{
			Room:     roomID,
			Identity: identity,
			TrackSid: track.Sid,
			Muted:    muted,
		}); err != nil {
			return err
		}
	}
	return nil
}

// RemoveParticipant disconnects the participant from the LiveKit room. The
// caller is responsible for updating DB state.
func (c *Client) RemoveParticipant(ctx context.Context, roomID, identity string) error {
	_, err := c.roomSvc.RemoveParticipant(ctx, &livekitproto.RoomParticipantIdentity{
		Room:     roomID,
		Identity: identity,
	})
	return err
}

// toHTTPURL converts a LiveKit ws:// / wss:// URL to http:// / https:// for
// the admin REST API.
func toHTTPURL(u string) string {
	switch {
	case len(u) >= 6 && u[:6] == "wss://":
		return "https://" + u[6:]
	case len(u) >= 5 && u[:5] == "ws://":
		return "http://" + u[5:]
	default:
		return u
	}
}
