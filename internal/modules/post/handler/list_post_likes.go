package handler

import (
	"context"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) ListPostLikes(ctx context.Context, req *connect.Request[postv1.ListPostLikesRequest]) (*connect.Response[postv1.ListPostLikesResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	views, err := h.listPostLikes.Execute(ctx, req.Msg.PostId, requestorID, req.Msg.PageSize)
	if err != nil {
		return nil, err
	}

	users := make([]*userv1.User, len(views))
	for i, v := range views {
		users[i] = toProtoUserFromView(v)
	}

	return connect.NewResponse(&postv1.ListPostLikesResponse{
		Users: users,
	}), nil
}
