package usecase

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/xerrors"
)

const (
	defaultListPageSize int32 = 50
	maxListPageSize     int32 = 100
)

type ListPublicCalls struct {
	callRepo    callrepo.Repository
	callGateway callgw.Gateway
}

func NewListPublicCalls(callRepo callrepo.Repository, callGateway callgw.Gateway) *ListPublicCalls {
	return &ListPublicCalls{
		callRepo:    callRepo,
		callGateway: callGateway,
	}
}

func (uc *ListPublicCalls) Execute(ctx context.Context, input ListPublicCallsInput) (*ListPublicCallsOutput, error) {
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultListPageSize
	}
	if pageSize > maxListPageSize {
		pageSize = maxListPageSize
	}

	// Guests see OPEN only; authenticated callers additionally see USERS_ONLY.
	includeUsersOnly := !input.RequesterID.IsZero()

	callIDs, err := uc.callRepo.ListActivePublicCallIDs(ctx, includeUsersOnly, entity.ParticipantStaleSeconds, input.PageToken, pageSize)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	nextPageToken := ""
	if int32(len(callIDs)) == pageSize {
		nextPageToken = callIDs[len(callIDs)-1].String()
	}

	views, err := uc.callGateway.BuildListCallViews(ctx, input.RequesterID, callIDs)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &ListPublicCallsOutput{
		Views:         views,
		NextPageToken: nextPageToken,
	}, nil
}
