package usecase

import (
	"context"

	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListUserFollowers struct {
	followGateway followgw.Gateway
	userGateway   usergw.Gateway
}

func NewListUserFollowers(followGateway followgw.Gateway, userGateway usergw.Gateway) *ListUserFollowers {
	return &ListUserFollowers{followGateway: followGateway, userGateway: userGateway}
}

func (uc *ListUserFollowers) Execute(ctx context.Context, input ListUserFollowersInput) (*ListUserFollowersOutput, error) {
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

	ids, err := uc.followGateway.ListFollowerIDs(ctx, target.ID, input.PageToken, pageSize)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	nextPageToken := ""
	if int32(len(ids)) == pageSize {
		nextPageToken = ids[len(ids)-1].String()
	}

	views, err := uc.userGateway.BuildListUserViews(ctx, input.RequesterID, ids)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	ordered := orderViewsByIDs(views, ids)

	return &ListUserFollowersOutput{Views: ordered, NextPageToken: nextPageToken}, nil
}
