package handler

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/modules/account/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) GetAccount(ctx context.Context, req *connect.Request[accountv1.GetAccountRequest]) (*connect.Response[accountv1.GetAccountResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	output, err := h.getAccount.Execute(ctx, usecase.GetAccountInput{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.GetAccountResponse{
		Account: &accountv1.Account{
			Id:          output.ID.String(),
			CustomId:    output.CustomID,
			DisplayName: output.DisplayName,
			AvatarUrl:   output.AvatarURL,
		},
	}), nil
}
