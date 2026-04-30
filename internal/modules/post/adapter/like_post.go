package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromLikePostRequest(ctx context.Context, req *connect.Request[postv1.LikePostRequest]) (usecase.LikePostInput, ulid.ULID, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.LikePostInput{}, ulid.ULID{}, xerrors.ErrUnauthenticated
	}
	return usecase.LikePostInput{
		AuthorCustomID: req.Msg.AuthorId,
		ShortID:        req.Msg.ShortId,
	}, userID, nil
}

func ToLikePostResponse(output *usecase.PostOutput) *connect.Response[postv1.LikePostResponse] {
	return connect.NewResponse(&postv1.LikePostResponse{Post: ToPost(output)})
}

func FromUnlikePostRequest(ctx context.Context, req *connect.Request[postv1.UnlikePostRequest]) (usecase.UnlikePostInput, ulid.ULID, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.UnlikePostInput{}, ulid.ULID{}, xerrors.ErrUnauthenticated
	}
	return usecase.UnlikePostInput{
		AuthorCustomID: req.Msg.AuthorId,
		ShortID:        req.Msg.ShortId,
	}, userID, nil
}

func ToUnlikePostResponse(output *usecase.PostOutput) *connect.Response[postv1.UnlikePostResponse] {
	return connect.NewResponse(&postv1.UnlikePostResponse{Post: ToPost(output)})
}
