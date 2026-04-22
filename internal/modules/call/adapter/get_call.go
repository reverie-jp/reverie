package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
	"reverie.jp/reverie/internal/platform/resourcename"
)

func FromGetCallRequest(ctx context.Context, req *connect.Request[callv1.GetCallRequest]) (usecase.GetCallInput, error) {
	callID, err := resourcename.ParseCall(req.Msg.Name)
	if err != nil {
		return usecase.GetCallInput{}, err
	}
	input := usecase.GetCallInput{CallID: callID}
	if userID, ok := interceptor.UserIDFromContext(ctx); ok {
		input.RequesterID = userID
		input.ViewerIdentity = "user:" + userID.String()
	} else {
		input.ViewerIdentity = req.Msg.GuestIdentity
	}
	return input, nil
}

func ToGetCallResponse(output *usecase.GetCallOutput) *connect.Response[callv1.GetCallResponse] {
	participants := make([]*callv1.CallParticipant, len(output.Participants))
	for i, p := range output.Participants {
		participants[i] = ToCallParticipant(p)
	}
	return connect.NewResponse(&callv1.GetCallResponse{
		Call:         ToCall(output.View),
		Participants: participants,
	})
}
