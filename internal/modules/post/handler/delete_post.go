package handler

import (
	"context"

	"connectrpc.com/connect"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/application/server/interceptor"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) DeletePost(ctx context.Context, req *connect.Request[postv1.DeletePostRequest]) (*connect.Response[postv1.DeletePostResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	err := h.deletePost.Execute(ctx, usecase.DeletePostInput{
		PostID: req.Msg.PostId,
	}, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&postv1.DeletePostResponse{}), nil
}
