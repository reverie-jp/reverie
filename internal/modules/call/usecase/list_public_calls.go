package usecase

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

const (
	defaultListPageSize int32 = 50
	maxListPageSize     int32 = 100
)

type ListPublicCalls struct {
	callRepo    callrepo.Repository
	userGateway usergw.Gateway
}

func NewListPublicCalls(callRepo callrepo.Repository, userGateway usergw.Gateway) *ListPublicCalls {
	return &ListPublicCalls{
		callRepo:    callRepo,
		userGateway: userGateway,
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

	calls, err := uc.callRepo.ListActivePublicCalls(ctx, participantStaleSeconds, input.PageToken, pageSize)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	// Cursor is taken from the DB-ordered last row, not the post-filter last
	// row, so that the next page resumes after the DB scan boundary.
	nextPageToken := ""
	if int32(len(calls)) == pageSize {
		nextPageToken = calls[len(calls)-1].ID.String()
	}

	// Guests see OPEN only; authenticated callers additionally see USERS_ONLY.
	if input.RequesterID.IsZero() {
		filtered := make([]*entity.Call, 0, len(calls))
		for _, c := range calls {
			if c.Visibility == entity.CallVisibilityOpen {
				filtered = append(filtered, c)
			}
		}
		calls = filtered
	}

	hostIDs := make([]ulid.ULID, 0, len(calls))
	seen := make(map[ulid.ULID]bool, len(calls))
	for _, c := range calls {
		if !seen[c.HostUserID] {
			hostIDs = append(hostIDs, c.HostUserID)
			seen[c.HostUserID] = true
		}
	}
	views, err := uc.userGateway.BuildListViews(ctx, input.RequesterID, hostIDs)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	hostsByID := make(map[ulid.ULID]*usergw.UserView, len(views))
	for _, v := range views {
		if v != nil && v.User != nil {
			hostsByID[v.User.ID] = v
		}
	}

	return &ListPublicCallsOutput{
		Calls:         calls,
		HostsByID:     hostsByID,
		NextPageToken: nextPageToken,
	}, nil
}
