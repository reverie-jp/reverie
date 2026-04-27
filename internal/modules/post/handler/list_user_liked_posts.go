package handler

import (
	"context"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) ListUserLikedPosts(ctx context.Context, req *connect.Request[postv1.ListUserLikedPostsRequest]) (*connect.Response[postv1.ListUserLikedPostsResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	outputs, err := h.listUserLikedPosts.Execute(ctx, req.Msg.UserId, req.Msg.PageToken, req.Msg.PageSize, requestorID)
	if err != nil {
		return nil, err
	}

	posts := make([]*postv1.Post, len(outputs))
	for i, o := range outputs {
		posts[i] = toProtoPost(o)
	}

	return connect.NewResponse(&postv1.ListUserLikedPostsResponse{
		Posts:         posts,
		NextPageToken: nextPageToken(outputs),
	}), nil
}
