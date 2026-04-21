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

func FromEndCallRequest(ctx context.Context, req *connect.Request[callv1.EndCallRequest]) (usecase.EndCallInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.EndCallInput{}, xerrors.ErrUnauthenticated
	}
	callID, err := resourcename.ParseCall(req.Msg.Name)
	if err != nil {
		return usecase.EndCallInput{}, err
	}
	return usecase.EndCallInput{
		RequesterID: userID,
		CallID:      callID,
	}, nil
}

func ToEndCallResponse() *connect.Response[callv1.EndCallResponse] {
	return connect.NewResponse(&callv1.EndCallResponse{})
}
