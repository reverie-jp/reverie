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

func FromListCallBansRequest(ctx context.Context, req *connect.Request[callv1.ListCallBansRequest]) (usecase.ListCallBansInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.ListCallBansInput{}, xerrors.ErrUnauthenticated
	}
	callID, err := resourcename.ParseCall(req.Msg.Parent)
	if err != nil {
		return usecase.ListCallBansInput{}, err
	}
	return usecase.ListCallBansInput{
		RequesterID: userID,
		CallID:      callID,
		PageSize:    req.Msg.PageSize,
		PageToken:   req.Msg.PageToken,
	}, nil
}

func ToListCallBansResponse(output *usecase.ListCallBansOutput) *connect.Response[callv1.ListCallBansResponse] {
	bans := make([]*callv1.CallBan, 0, len(output.Bans))
	for _, b := range output.Bans {
		bans = append(bans, ToCallBan(b))
	}
	return connect.NewResponse(&callv1.ListCallBansResponse{
		Bans:          bans,
		NextPageToken: output.NextPageToken,
	})
}
