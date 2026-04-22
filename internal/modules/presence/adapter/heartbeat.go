package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	presencev1 "reverie.jp/reverie/internal/gen/pb/presence/v1"
	"reverie.jp/reverie/internal/modules/presence/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromHeartbeatRequest(ctx context.Context, _ *connect.Request[presencev1.HeartbeatRequest]) (usecase.HeartbeatInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.HeartbeatInput{}, xerrors.ErrUnauthenticated
	}
	return usecase.HeartbeatInput{RequesterID: requesterID}, nil
}

func ToHeartbeatResponse(_ *usecase.HeartbeatOutput) *connect.Response[presencev1.HeartbeatResponse] {
	return connect.NewResponse(&presencev1.HeartbeatResponse{})
}
