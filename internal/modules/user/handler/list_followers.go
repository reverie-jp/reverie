package handler

import (
	"context"
	"encoding/base64"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) ListFollowers(ctx context.Context, req *connect.Request[userv1.ListFollowingRequest]) (*connect.Response[userv1.ListFollowingResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	outputs, err := h.listFollowers.Execute(ctx, req.Msg.UserId, req.Msg.PageToken, req.Msg.PageSize, requestorID)
	if err != nil {
		return nil, err
	}

	users := make([]*userv1.User, len(outputs))
	for i, o := range outputs {
		users[i] = toProtoUser(o)
	}

	var nextPageToken string
	if len(outputs) > 0 {
		last := outputs[len(outputs)-1]
		raw := last.CreateTime.UTC().Format("2006-01-02T15:04:05.999999999Z")
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(raw))
	}

	return connect.NewResponse(&userv1.ListFollowingResponse{
		Users:         users,
		NextPageToken: nextPageToken,
	}), nil
}
