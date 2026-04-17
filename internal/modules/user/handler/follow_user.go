package handler

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"reverie.jp/reverie/internal/application/server/interceptor"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/user/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) FollowUser(ctx context.Context, req *connect.Request[userv1.FollowUserRequest]) (*connect.Response[userv1.FollowUserResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	output, err := h.followUser.Execute(ctx, req.Msg.UserId, requestorID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&userv1.FollowUserResponse{
		User: toProtoUser(output),
	}), nil
}

func (h *Handler) UnfollowUser(ctx context.Context, req *connect.Request[userv1.UnfollowUserRequest]) (*connect.Response[userv1.UnfollowUserResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	output, err := h.unfollowUser.Execute(ctx, req.Msg.UserId, requestorID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&userv1.UnfollowUserResponse{
		User: toProtoUser(output),
	}), nil
}

func toProtoUser(out *usecase.GetUserOutput) *userv1.User {
	return &userv1.User{
		Id:             out.ID.String(),
		CustomId:       out.CustomID,
		DisplayName:    out.DisplayName,
		Biography:      &out.Biography,
		IsPrivate:      out.IsPrivate,
		IsMe:           out.IsMe,
		IsFollowing:    out.IsFollowing,
		IsFollowedBy:   out.IsFollowedBy,
		FollowerCount:  int32(out.FollowerCount),
		FollowingCount: int32(out.FollowingCount),
		CreateTime:     timestamppb.New(out.CreateTime),
	}
}
