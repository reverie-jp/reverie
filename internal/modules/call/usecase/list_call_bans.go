package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

const (
	defaultBanPageSize = 50
	maxBanPageSize     = 100
)

type ListCallBans struct {
	callRepo    callrepo.Repository
	userGateway usergw.Gateway
}

func NewListCallBans(callRepo callrepo.Repository, userGateway usergw.Gateway) *ListCallBans {
	return &ListCallBans{callRepo: callRepo, userGateway: userGateway}
}

func (uc *ListCallBans) Execute(ctx context.Context, input ListCallBansInput) (*ListCallBansOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	call, err := uc.callRepo.GetCall(ctx, input.CallID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if call == nil {
		return nil, xerrors.ErrCallNotFound
	}
	if call.HostUserID != input.RequesterID {
		return nil, xerrors.ErrNotCallHost
	}

	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultBanPageSize
	}
	if pageSize > maxBanPageSize {
		pageSize = maxBanPageSize
	}

	bans, err := uc.callRepo.ListCallBans(ctx, call.ID, input.PageToken, pageSize)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	nextPageToken := ""
	if int32(len(bans)) == pageSize {
		nextPageToken = bans[len(bans)-1].UserID.String()
	}

	userIDs := make([]ulid.ULID, 0, len(bans))
	for _, b := range bans {
		userIDs = append(userIDs, b.UserID)
	}
	views, err := uc.userGateway.BuildListViews(ctx, input.RequesterID, userIDs)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	viewByID := make(map[ulid.ULID]*usergw.UserView, len(views))
	for _, v := range views {
		if v != nil && v.User != nil {
			viewByID[v.User.ID] = v
		}
	}

	result := make([]*CallBanView, len(bans))
	for i, b := range bans {
		result[i] = &CallBanView{Ban: b, User: viewByID[b.UserID]}
	}

	return &ListCallBansOutput{Bans: result, NextPageToken: nextPageToken}, nil
}
