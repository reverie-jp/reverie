package usecase

import (
	"context"

	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/xerrors"
)

const (
	defaultBanPageSize = 50
	maxBanPageSize     = 100
)

type ListCallBans struct {
	callRepo    callrepo.Repository
	callGateway callgw.Gateway
}

func NewListCallBans(callRepo callrepo.Repository, callGateway callgw.Gateway) *ListCallBans {
	return &ListCallBans{callRepo: callRepo, callGateway: callGateway}
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

	views, err := uc.callGateway.BuildListCallBanViews(ctx, input.RequesterID, bans)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &ListCallBansOutput{Bans: views, NextPageToken: nextPageToken}, nil
}
