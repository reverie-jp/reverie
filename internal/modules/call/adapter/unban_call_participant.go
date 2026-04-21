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

func FromUnbanCallParticipantRequest(ctx context.Context, req *connect.Request[callv1.UnbanCallParticipantRequest]) (usecase.UnbanCallParticipantInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.UnbanCallParticipantInput{}, xerrors.ErrUnauthenticated
	}
	callID, targetUserID, err := resourcename.ParseCallBan(req.Msg.Name)
	if err != nil {
		return usecase.UnbanCallParticipantInput{}, err
	}
	return usecase.UnbanCallParticipantInput{
		RequesterID: userID,
		CallID:      callID,
		UserID:      targetUserID,
	}, nil
}

func ToUnbanCallParticipantResponse() *connect.Response[callv1.UnbanCallParticipantResponse] {
	return connect.NewResponse(&callv1.UnbanCallParticipantResponse{})
}
