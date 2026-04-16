package handler

import (
	"context"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) GetPost(ctx context.Context, req *connect.Request[postv1.GetPostRequest]) (*connect.Response[postv1.GetPostResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	output, err := h.getPost.Execute(ctx, usecase.GetPostInput{
		PostID: req.Msg.PostId,
	}, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&postv1.GetPostResponse{
		Post: toProtoPost(output),
	}), nil
}
