package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromListFollowingCallsRequest(ctx context.Context, req *connect.Request[callv1.ListFollowingCallsRequest]) (usecase.ListFollowingCallsInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.ListFollowingCallsInput{}, xerrors.ErrUnauthenticated
	}
	return usecase.ListFollowingCallsInput{
		RequesterID: userID,
		PageSize:    req.Msg.PageSize,
		PageToken:   req.Msg.PageToken,
	}, nil
}

func ToListFollowingCallsResponse(output *usecase.ListFollowingCallsOutput) *connect.Response[callv1.ListFollowingCallsResponse] {
	calls := make([]*callv1.Call, 0, len(output.Views))
	for _, v := range output.Views {
		calls = append(calls, ToCall(v))
	}
	return connect.NewResponse(&callv1.ListFollowingCallsResponse{
		Calls:         calls,
		NextPageToken: output.NextPageToken,
	})
}
