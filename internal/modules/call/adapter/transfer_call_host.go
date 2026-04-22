package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
	"reverie.jp/reverie/internal/platform/resourcename"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromTransferCallHostRequest(ctx context.Context, req *connect.Request[callv1.TransferCallHostRequest]) (usecase.TransferCallHostInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.TransferCallHostInput{}, xerrors.ErrUnauthenticated
	}
	callID, err := resourcename.ParseCall(req.Msg.Name)
	if err != nil {
		return usecase.TransferCallHostInput{}, err
	}
	newHostCustomID, err := resourcename.ParseUser(req.Msg.NewHost)
	if err != nil {
		return usecase.TransferCallHostInput{}, err
	}
	return usecase.TransferCallHostInput{
		RequesterID:     userID,
		CallID:          callID,
		NewHostCustomID: newHostCustomID,
	}, nil
}

func ToTransferCallHostResponse(output *usecase.TransferCallHostOutput) *connect.Response[callv1.TransferCallHostResponse] {
	return connect.NewResponse(&callv1.TransferCallHostResponse{
		Call: ToCall(output.View),
	})
}
