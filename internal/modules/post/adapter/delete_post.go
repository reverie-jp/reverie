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

func FromDeletePostRequest(ctx context.Context, req *connect.Request[postv1.DeletePostRequest]) (usecase.DeletePostInput, ulid.ULID, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.DeletePostInput{}, ulid.ULID{}, xerrors.ErrUnauthenticated
	}
	return usecase.DeletePostInput{
		AuthorCustomID: req.Msg.AuthorId,
		ShortID:        req.Msg.ShortId,
	}, userID, nil
}
