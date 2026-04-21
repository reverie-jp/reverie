package call

import (
	"time"

	"reverie.jp/reverie/internal/gen/pb/call/v1/callv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/call/handler"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/modules/call/usecase"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/livekit"
)

func InitModule(q *sqlc.Queries, userGateway usergw.Gateway, lk *livekit.Client, tokenTTL time.Duration) callv1connect.CallServiceHandler {
	callRepo := callrepo.New(q)

	createCall := usecase.NewCreateCall(callRepo, userGateway)
	getCall := usecase.NewGetCall(callRepo, userGateway)
	updateCall := usecase.NewUpdateCall(callRepo, userGateway)
	listPublicCalls := usecase.NewListPublicCalls(callRepo, userGateway)
	getUserParticipatingCall := usecase.NewGetUserParticipatingCall(callRepo, userGateway)
	joinCall := usecase.NewJoinCall(callRepo, userGateway, lk, tokenTTL)
	heartbeatCall := usecase.NewHeartbeatCall(callRepo)
	leaveCall := usecase.NewLeaveCall(callRepo)
	muteCallParticipant := usecase.NewMuteCallParticipant(callRepo, lk)
	unmuteCallParticipant := usecase.NewUnmuteCallParticipant(callRepo, lk)
	kickCallParticipant := usecase.NewKickCallParticipant(callRepo, lk)
	banCallParticipant := usecase.NewBanCallParticipant(callRepo, lk)
	transferCallHost := usecase.NewTransferCallHost(callRepo, userGateway)
	endCall := usecase.NewEndCall(callRepo, lk)
	listCallBans := usecase.NewListCallBans(callRepo, userGateway)
	unbanCallParticipant := usecase.NewUnbanCallParticipant(callRepo)

	return handler.New(
		createCall,
		getCall,
		updateCall,
		listPublicCalls,
		getUserParticipatingCall,
		joinCall,
		heartbeatCall,
		leaveCall,
		muteCallParticipant,
		unmuteCallParticipant,
		kickCallParticipant,
		banCallParticipant,
		transferCallHost,
		endCall,
		listCallBans,
		unbanCallParticipant,
	)
}
