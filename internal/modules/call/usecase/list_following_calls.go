package usecase

import (
	"context"

	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListFollowingCalls struct {
	callRepo    callrepo.Repository
	callGateway callgw.Gateway
}

func NewListFollowingCalls(callRepo callrepo.Repository, callGateway callgw.Gateway) *ListFollowingCalls {
	return &ListFollowingCalls{
		callRepo:    callRepo,
		callGateway: callGateway,
	}
}

func (uc *ListFollowingCalls) Execute(ctx context.Context, input ListFollowingCallsInput) (*ListFollowingCallsOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultListPageSize
	}
	if pageSize > maxListPageSize {
		pageSize = maxListPageSize
	}

	callIDs, err := uc.callRepo.ListActiveCallIDsForFollower(ctx, input.RequesterID, participantStaleSeconds, input.PageToken, pageSize)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	nextPageToken := ""
	if int32(len(callIDs)) == pageSize {
		nextPageToken = callIDs[len(callIDs)-1].String()
	}

	views, err := uc.callGateway.BuildListViews(ctx, input.RequesterID, callIDs)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &ListFollowingCallsOutput{
		Views:         views,
		NextPageToken: nextPageToken,
	}, nil
}
