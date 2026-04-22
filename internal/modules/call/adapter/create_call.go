package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromCreateCallRequest(ctx context.Context, req *connect.Request[callv1.CreateCallRequest]) (usecase.CreateCallInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.CreateCallInput{}, xerrors.ErrUnauthenticated
	}
	return usecase.CreateCallInput{
		RequesterID: userID,
		Visibility:  fromProtoVisibility(req.Msg.Visibility),
	}, nil
}

func ToCreateCallResponse(output *usecase.CreateCallOutput) *connect.Response[callv1.CreateCallResponse] {
	return connect.NewResponse(&callv1.CreateCallResponse{
		Call: ToCall(output.View),
	})
}
