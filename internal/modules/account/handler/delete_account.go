package handler

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/modules/account/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) DeleteAccount(ctx context.Context, req *connect.Request[accountv1.DeleteAccountRequest]) (*connect.Response[accountv1.DeleteAccountResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	err := h.deleteAccount.Execute(ctx, usecase.DeleteAccountInput{
		UserID:          userID,
		ConfirmCustomID: req.Msg.ConfirmCustomId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.DeleteAccountResponse{}), nil
}
