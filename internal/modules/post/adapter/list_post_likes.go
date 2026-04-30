package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type ListPostLikesParams struct {
	AuthorCustomID string
	ShortID        string
	UserID         ulid.ULID
	PageSize       int32
}

func FromListPostLikesRequest(ctx context.Context, req *connect.Request[postv1.ListPostLikesRequest]) (ListPostLikesParams, error) {
	userID, _ := interceptor.UserIDFromContext(ctx)
	return ListPostLikesParams{
		AuthorCustomID: req.Msg.AuthorId,
		ShortID:        req.Msg.ShortId,
		UserID:         userID,
		PageSize:       req.Msg.PageSize,
	}, nil
}

func ToListPostLikesResponse(views []*usergw.UserView) *connect.Response[postv1.ListPostLikesResponse] {
	users := make([]*userv1.User, len(views))
	for i, v := range views {
		users[i] = useradapter.ToUser(v)
	}
	return connect.NewResponse(&postv1.ListPostLikesResponse{Users: users})
}
