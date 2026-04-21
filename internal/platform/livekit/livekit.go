package livekit

import (
	"time"

	"github.com/livekit/protocol/auth"
)

type Client struct {
	url       string
	apiKey    string
	apiSecret string
}

func NewClient(url, apiKey, apiSecret string) *Client {
	return &Client{
		url:       url,
		apiKey:    apiKey,
		apiSecret: apiSecret,
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
