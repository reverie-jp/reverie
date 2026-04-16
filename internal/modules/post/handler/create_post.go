package handler

import (
	"context"

	"connectrpc.com/connect"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/application/server/interceptor"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) CreatePost(ctx context.Context, req *connect.Request[postv1.CreatePostRequest]) (*connect.Response[postv1.CreatePostResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	output, err := h.createPost.Execute(ctx, usecase.CreatePostInput{
		Text:      req.Msg.Text,
		ReplyToID: req.Msg.ReplyToId,
		RepostID:  req.Msg.RepostId,
	}, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&postv1.CreatePostResponse{
		Post: toProtoPost(output),
	}), nil
}
