package handler

import (
	"context"

	"connectrpc.com/connect"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/application/server/interceptor"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) ListTimeline(ctx context.Context, req *connect.Request[postv1.ListTimelineRequest]) (*connect.Response[postv1.ListTimelineResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	outputs, err := h.listTimeline.Execute(ctx, usecase.ListTimelineInput{
		Cursor: req.Msg.Cursor,
		Limit:  req.Msg.Limit,
	}, userID)
	if err != nil {
		return nil, err
	}

	posts := make([]*postv1.Post, len(outputs))
	for i, o := range outputs {
		posts[i] = toProtoPost(o)
	}

	return connect.NewResponse(&postv1.ListTimelineResponse{
		Posts:      posts,
		NextCursor: nextCursor(outputs),
	}), nil
}
