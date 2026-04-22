package call

import (
	"time"

	"reverie.jp/reverie/internal/gen/pb/call/v1/callv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	"reverie.jp/reverie/internal/modules/call/handler"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/modules/call/usecase"
	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/livekit"
)

func InitModule(
	q *sqlc.Queries,
	userGateway usergw.Gateway,
	followGateway followgw.Gateway,
	notificationGateway notificationgw.Gateway,
	lk *livekit.Client,
	tokenTTL time.Duration,
) callv1connect.CallServiceHandler {
	callRepo := callrepo.New(q)
	callGateway := callgw.New(callRepo, userGateway)

	createCall := usecase.NewCreateCall(callRepo, callGateway, followGateway, notificationGateway)
	getCall := usecase.NewGetCall(callRepo, callGateway)
	updateCall := usecase.NewUpdateCall(callRepo, callGateway)
	listPublicCalls := usecase.NewListPublicCalls(callRepo, callGateway)
	listFollowingCalls := usecase.NewListFollowingCalls(callRepo, callGateway)
	getUserParticipatingCall := usecase.NewGetUserParticipatingCall(callRepo, userGateway, callGateway)
	joinCall := usecase.NewJoinCall(callRepo, userGateway, lk, tokenTTL)
	heartbeatCall := usecase.NewHeartbeatCall(callRepo)
	leaveCall := usecase.NewLeaveCall(callRepo, lk)
	muteCallParticipant := usecase.NewMuteCallParticipant(callRepo, lk)
	unmuteCallParticipant := usecase.NewUnmuteCallParticipant(callRepo, lk)
	kickCallParticipant := usecase.NewKickCallParticipant(callRepo, lk)
	banCallParticipant := usecase.NewBanCallParticipant(callRepo, lk)
	transferCallHost := usecase.NewTransferCallHost(callRepo, userGateway, callGateway)
	endCall := usecase.NewEndCall(callRepo, lk)
	listCallBans := usecase.NewListCallBans(callRepo, callGateway)
	unbanCallParticipant := usecase.NewUnbanCallParticipant(callRepo)

	return handler.New(
		createCall,
		getCall,
		updateCall,
		listPublicCalls,
		listFollowingCalls,
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
