package handler

import (
	"context"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) ListUserPosts(ctx context.Context, req *connect.Request[postv1.ListUserPostsRequest]) (*connect.Response[postv1.ListUserPostsResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	outputs, err := h.listUserPosts.Execute(ctx, usecase.ListUserPostsInput{
		UserID:    req.Msg.UserId,
		PageToken: req.Msg.PageToken,
		PageSize:  req.Msg.PageSize,
	}, userID)
	if err != nil {
		return nil, err
	}

	posts := make([]*postv1.Post, len(outputs))
	for i, o := range outputs {
		posts[i] = toProtoPost(o)
	}

	return connect.NewResponse(&postv1.ListUserPostsResponse{
		Posts:         posts,
		NextPageToken: nextPageToken(outputs),
	}), nil
}
