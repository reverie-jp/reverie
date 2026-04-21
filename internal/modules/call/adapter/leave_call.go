package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
	"reverie.jp/reverie/internal/platform/resourcename"
)

func FromLeaveCallRequest(ctx context.Context, req *connect.Request[callv1.LeaveCallRequest]) (usecase.LeaveCallInput, error) {
	callID, err := resourcename.ParseCall(req.Msg.Name)
	if err != nil {
		return usecase.LeaveCallInput{}, err
	}
	input := usecase.LeaveCallInput{CallID: callID}
	if userID, ok := interceptor.UserIDFromContext(ctx); ok {
		input.RequesterID = userID
		input.Identity = "user:" + userID.String()
	} else {
		input.Identity = req.Msg.GuestIdentity
	}
	return input, nil
}

func ToLeaveCallResponse() *connect.Response[callv1.LeaveCallResponse] {
	return connect.NewResponse(&callv1.LeaveCallResponse{})
}
