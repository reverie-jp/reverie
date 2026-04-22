package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type CallView struct {
	Call *entity.Call
	Host *usergw.UserView
	// ActiveParticipants is populated on list endpoints for avatar stacks
	// on call cards. Includes both authenticated participants and guests.
	// Ordered by first_join_time. Empty on single-call reads.
	ActiveParticipants []*CallParticipantView
}

type CallParticipantView struct {
	Participant          *entity.CallParticipant
	User                 *usergw.UserView
	IsCurrentlyConnected bool
}

type CallBanView struct {
	Ban  *entity.CallBan
	User *usergw.UserView
}

type Gateway interface {
	BuildCallView(ctx context.Context, requesterID, callID ulid.ULID) (*CallView, error)
	BuildListCallViews(ctx context.Context, requesterID ulid.ULID, callIDs []ulid.ULID) ([]*CallView, error)
	BuildListParticipantViews(ctx context.Context, requesterID ulid.ULID, participants []*entity.CallParticipant) ([]*CallParticipantView, error)
	BuildListCallBanViews(ctx context.Context, requesterID ulid.ULID, bans []*entity.CallBan) ([]*CallBanView, error)
}

type gatewayImpl struct {
	repo        repository.Repository
	userGateway usergw.Gateway
}

func New(repo repository.Repository, userGateway usergw.Gateway) Gateway {
	return &gatewayImpl{
		repo:        repo,
		userGateway: userGateway,
	}
}
