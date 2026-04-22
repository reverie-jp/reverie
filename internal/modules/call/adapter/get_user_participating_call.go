package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
	"reverie.jp/reverie/internal/platform/resourcename"
)

func FromGetUserParticipatingCallRequest(ctx context.Context, req *connect.Request[callv1.GetUserParticipatingCallRequest]) (usecase.GetUserParticipatingCallInput, error) {
	customID, err := resourcename.ParseUser(req.Msg.Name)
	if err != nil {
		return usecase.GetUserParticipatingCallInput{}, err
	}
	input := usecase.GetUserParticipatingCallInput{
		UserCustomID: customID,
	}
	if userID, ok := interceptor.UserIDFromContext(ctx); ok {
		input.RequesterID = userID
	}
	return input, nil
}

func ToGetUserParticipatingCallResponse(output *usecase.GetUserParticipatingCallOutput) *connect.Response[callv1.GetUserParticipatingCallResponse] {
	return connect.NewResponse(&callv1.GetUserParticipatingCallResponse{
		Call: ToCall(output.View),
	})
}
