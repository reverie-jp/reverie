package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
)

func FromListPublicCallsRequest(ctx context.Context, req *connect.Request[callv1.ListPublicCallsRequest]) usecase.ListPublicCallsInput {
	input := usecase.ListPublicCallsInput{
		PageSize:  req.Msg.PageSize,
		PageToken: req.Msg.PageToken,
	}
	if userID, ok := interceptor.UserIDFromContext(ctx); ok {
		input.RequesterID = userID
	}
	return input
}

func ToListPublicCallsResponse(output *usecase.ListPublicCallsOutput) *connect.Response[callv1.ListPublicCallsResponse] {
	calls := make([]*callv1.Call, 0, len(output.Views))
	for _, v := range output.Views {
		calls = append(calls, ToCall(v.Call, v.Host))
	}
	return connect.NewResponse(&callv1.ListPublicCallsResponse{
		Calls:         calls,
		NextPageToken: output.NextPageToken,
	})
}
