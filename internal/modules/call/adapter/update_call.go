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

func FromUpdateCallRequest(ctx context.Context, req *connect.Request[callv1.UpdateCallRequest]) (usecase.UpdateCallInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.UpdateCallInput{}, xerrors.ErrUnauthenticated
	}
	if req.Msg.Call == nil {
		return usecase.UpdateCallInput{}, xerrors.ErrInvalidArgument.WithMessage("call is required")
	}
	callID, err := resourcename.ParseCall(req.Msg.Call.Name)
	if err != nil {
		return usecase.UpdateCallInput{}, err
	}
	var mask []string
	if req.Msg.UpdateMask != nil {
		mask = req.Msg.UpdateMask.Paths
	}
	return usecase.UpdateCallInput{
		RequesterID: userID,
		CallID:      callID,
		Visibility:  fromProtoVisibility(req.Msg.Call.Visibility),
		UpdateMask:  mask,
	}, nil
}

func ToUpdateCallResponse(output *usecase.UpdateCallOutput) *connect.Response[callv1.UpdateCallResponse] {
	return connect.NewResponse(&callv1.UpdateCallResponse{
		Call: ToCall(output.Call, output.Host),
	})
}
