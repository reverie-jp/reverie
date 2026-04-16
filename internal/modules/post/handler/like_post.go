package handler

import (
	"context"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) LikePost(ctx context.Context, req *connect.Request[postv1.LikePostRequest]) (*connect.Response[postv1.LikePostResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	output, err := h.likePost.Execute(ctx, usecase.LikePostInput{
		PostID: req.Msg.PostId,
	}, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&postv1.LikePostResponse{
		Post: toProtoPost(output),
	}), nil
}

func (h *Handler) UnlikePost(ctx context.Context, req *connect.Request[postv1.UnlikePostRequest]) (*connect.Response[postv1.UnlikePostResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	output, err := h.unlikePost.Execute(ctx, usecase.UnlikePostInput{
		PostID: req.Msg.PostId,
	}, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&postv1.UnlikePostResponse{
		Post: toProtoPost(output),
	}), nil
}
