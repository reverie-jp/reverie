package handler

import (
	"context"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) ListPostReposts(ctx context.Context, req *connect.Request[postv1.ListPostRepostsRequest]) (*connect.Response[postv1.ListPostRepostsResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	outputs, err := h.listPostReposts.Execute(ctx, usecase.ListPostRepostsInput{
		PostID:    req.Msg.PostId,
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

	return connect.NewResponse(&postv1.ListPostRepostsResponse{
		Posts:         posts,
		NextPageToken: nextPageToken(outputs),
	}), nil
}
