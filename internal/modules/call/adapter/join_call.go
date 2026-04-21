package adapter

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"reverie.jp/reverie/internal/application/server/interceptor"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
)

func FromJoinCallRequest(ctx context.Context, req *connect.Request[callv1.JoinCallRequest]) (usecase.JoinCallInput, error) {
	input := usecase.JoinCallInput{
		RoomID:           req.Msg.RoomId,
		GuestDisplayName: req.Msg.GuestDisplayName,
	}
	if userID, ok := interceptor.UserIDFromContext(ctx); ok {
		input.RequesterID = userID
	}
	return input, nil
}

func ToJoinCallResponse(output *usecase.JoinCallOutput) *connect.Response[callv1.JoinCallResponse] {
	return connect.NewResponse(&callv1.JoinCallResponse{
		AccessToken: output.AccessToken,
		Url:         output.URL,
		Identity:    output.Identity,
		ExpireTime:  timestamppb.New(output.ExpireTime),
	})
}
