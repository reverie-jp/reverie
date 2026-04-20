package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/modules/account/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromDeleteAccountRequest(ctx context.Context, req *connect.Request[accountv1.DeleteAccountRequest]) (usecase.DeleteAccountInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.DeleteAccountInput{}, xerrors.ErrUnauthenticated
	}
	return usecase.DeleteAccountInput{
		UserID:          userID,
		ConfirmCustomID: req.Msg.ConfirmCustomId,
	}, nil
}

func ToDeleteAccountResponse() *connect.Response[accountv1.DeleteAccountResponse] {
	return connect.NewResponse(&accountv1.DeleteAccountResponse{})
}
