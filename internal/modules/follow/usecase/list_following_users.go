package usecase

import (
	"context"

	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListFollowingUsers struct {
	followGateway followgw.Gateway
	userGateway   usergw.Gateway
}

func NewListFollowingUsers(followGateway followgw.Gateway, userGateway usergw.Gateway) *ListFollowingUsers {
	return &ListFollowingUsers{followGateway: followGateway, userGateway: userGateway}
}

func (uc *ListFollowingUsers) Execute(ctx context.Context, input ListFollowingUsersInput) (*ListFollowingUsersOutput, error) {
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

	target, err := uc.userGateway.GetUserByCustomID(ctx, input.TargetCustomID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if target == nil {
		return nil, xerrors.ErrUserNotFound
	}

	ids, err := uc.followGateway.ListFollowingIDs(ctx, target.ID, input.PageToken, pageSize)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	nextPageToken := ""
	if int32(len(ids)) == pageSize {
		nextPageToken = ids[len(ids)-1].String()
	}

	views, err := uc.userGateway.BuildListViews(ctx, input.RequesterID, ids)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	ordered := orderViewsByIDs(views, ids)

	return &ListFollowingUsersOutput{Views: ordered, NextPageToken: nextPageToken}, nil
}

func orderViewsByIDs(views []*usergw.UserView, ids []ulid.ULID) []*usergw.UserView {
	byID := make(map[ulid.ULID]*usergw.UserView, len(views))
	for _, v := range views {
		if v != nil && v.User != nil {
			byID[v.User.ID] = v
		}
	}
	out := make([]*usergw.UserView, 0, len(ids))
	for _, id := range ids {
		if v, ok := byID[id]; ok {
			out = append(out, v)
		}
	}
	return out
}
