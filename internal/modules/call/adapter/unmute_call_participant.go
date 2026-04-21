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

func FromUnmuteCallParticipantRequest(ctx context.Context, req *connect.Request[callv1.UnmuteCallParticipantRequest]) (usecase.UnmuteCallParticipantInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.UnmuteCallParticipantInput{}, xerrors.ErrUnauthenticated
	}
	callID, identity, err := resourcename.ParseCallParticipant(req.Msg.Name)
	if err != nil {
		return usecase.UnmuteCallParticipantInput{}, err
	}
	return usecase.UnmuteCallParticipantInput{
		RequesterID: userID,
		CallID:      callID,
		Identity:    identity,
	}, nil
}

func ToUnmuteCallParticipantResponse() *connect.Response[callv1.UnmuteCallParticipantResponse] {
	return connect.NewResponse(&callv1.UnmuteCallParticipantResponse{})
}
