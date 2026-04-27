package handler

import (
	"context"
	"encoding/base64"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/user/adapter"
	"reverie.jp/reverie/internal/modules/user/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) SearchUsers(ctx context.Context, req *connect.Request[userv1.SearchUsersRequest]) (*connect.Response[userv1.SearchUsersResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	views, err := h.searchUsers.Execute(ctx, usecase.SearchUsersInput{
		Query:     req.Msg.Query,
		PageToken: req.Msg.PageToken,
		PageSize:  req.Msg.PageSize,
	}, requestorID)
	if err != nil {
		return nil, err
	}

	users := make([]*userv1.User, len(views))
	for i, v := range views {
		users[i] = adapter.ToUser(v)
	}

	var nextPageToken string
	if len(views) > 0 {
		last := views[len(views)-1]
		raw := last.User.CreateTime.UTC().Format("2006-01-02T15:04:05.999999999Z")
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(raw))
	}

	return connect.NewResponse(&userv1.SearchUsersResponse{
		Users:         users,
		NextPageToken: nextPageToken,
	}), nil
}
