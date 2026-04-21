package handler

import (
	"reverie.jp/reverie/internal/gen/pb/call/v1/callv1connect"
	"reverie.jp/reverie/internal/modules/call/usecase"
)

type Handler struct {
	callv1connect.UnimplementedCallServiceHandler
	createCall               *usecase.CreateCall
	getCall                  *usecase.GetCall
	updateCall               *usecase.UpdateCall
	listPublicCalls          *usecase.ListPublicCalls
	getUserParticipatingCall *usecase.GetUserParticipatingCall
	joinCall                 *usecase.JoinCall
	heartbeatCall            *usecase.HeartbeatCall
	leaveCall                *usecase.LeaveCall
	muteCallParticipant      *usecase.MuteCallParticipant
	unmuteCallParticipant    *usecase.UnmuteCallParticipant
	kickCallParticipant      *usecase.KickCallParticipant
	banCallParticipant       *usecase.BanCallParticipant
	transferCallHost         *usecase.TransferCallHost
	endCall                  *usecase.EndCall
	listCallBans             *usecase.ListCallBans
	unbanCallParticipant     *usecase.UnbanCallParticipant
}

func New(
	createCall *usecase.CreateCall,
	getCall *usecase.GetCall,
	updateCall *usecase.UpdateCall,
	listPublicCalls *usecase.ListPublicCalls,
	getUserParticipatingCall *usecase.GetUserParticipatingCall,
	joinCall *usecase.JoinCall,
	heartbeatCall *usecase.HeartbeatCall,
	leaveCall *usecase.LeaveCall,
	muteCallParticipant *usecase.MuteCallParticipant,
	unmuteCallParticipant *usecase.UnmuteCallParticipant,
	kickCallParticipant *usecase.KickCallParticipant,
	banCallParticipant *usecase.BanCallParticipant,
	transferCallHost *usecase.TransferCallHost,
	endCall *usecase.EndCall,
	listCallBans *usecase.ListCallBans,
	unbanCallParticipant *usecase.UnbanCallParticipant,
) *Handler {
	return &Handler{
		createCall:               createCall,
		getCall:                  getCall,
		updateCall:               updateCall,
		listPublicCalls:          listPublicCalls,
		getUserParticipatingCall: getUserParticipatingCall,
		joinCall:                 joinCall,
		heartbeatCall:            heartbeatCall,
		leaveCall:                leaveCall,
		muteCallParticipant:      muteCallParticipant,
		unmuteCallParticipant:    unmuteCallParticipant,
		kickCallParticipant:      kickCallParticipant,
		banCallParticipant:       banCallParticipant,
		transferCallHost:         transferCallHost,
		endCall:                  endCall,
		listCallBans:             listCallBans,
		unbanCallParticipant:     unbanCallParticipant,
	}
}
