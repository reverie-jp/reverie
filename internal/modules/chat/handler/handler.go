package handler

import (
	"reverie.jp/reverie/internal/gen/pb/chat/v1/chatv1connect"
	"reverie.jp/reverie/internal/modules/chat/usecase"
)

type Handler struct {
	chatv1connect.UnimplementedChatServiceHandler
	listRooms             *usecase.ListRooms
	createDirectRoom      *usecase.CreateDirectRoom
	listMessages          *usecase.ListMessages
	sendMessage           *usecase.SendMessage
	markRoomAsRead        *usecase.MarkRoomAsRead
	pinRoom               *usecase.PinRoom
	unpinRoom             *usecase.UnpinRoom
	muteRoom              *usecase.MuteRoom
	unmuteRoom            *usecase.UnmuteRoom
	leaveRoom             *usecase.LeaveRoom
	addMessageReaction    *usecase.AddMessageReaction
	removeMessageReaction *usecase.RemoveMessageReaction
	createGroupRoom       *usecase.CreateGroupRoom
	updateRoom            *usecase.UpdateRoom
	listRoomMembers       *usecase.ListRoomMembers
	addRoomMember         *usecase.AddRoomMember
	removeRoomMember      *usecase.RemoveRoomMember
}

func New(
	listRooms *usecase.ListRooms,
	createDirectRoom *usecase.CreateDirectRoom,
	listMessages *usecase.ListMessages,
	sendMessage *usecase.SendMessage,
	markRoomAsRead *usecase.MarkRoomAsRead,
	pinRoom *usecase.PinRoom,
	unpinRoom *usecase.UnpinRoom,
	muteRoom *usecase.MuteRoom,
	unmuteRoom *usecase.UnmuteRoom,
	leaveRoom *usecase.LeaveRoom,
	addMessageReaction *usecase.AddMessageReaction,
	removeMessageReaction *usecase.RemoveMessageReaction,
	createGroupRoom *usecase.CreateGroupRoom,
	updateRoom *usecase.UpdateRoom,
	listRoomMembers *usecase.ListRoomMembers,
	addRoomMember *usecase.AddRoomMember,
	removeRoomMember *usecase.RemoveRoomMember,
) *Handler {
	return &Handler{
		listRooms:             listRooms,
		createDirectRoom:      createDirectRoom,
		listMessages:          listMessages,
		sendMessage:           sendMessage,
		markRoomAsRead:        markRoomAsRead,
		pinRoom:               pinRoom,
		unpinRoom:             unpinRoom,
		muteRoom:              muteRoom,
		unmuteRoom:            unmuteRoom,
		leaveRoom:             leaveRoom,
		addMessageReaction:    addMessageReaction,
		removeMessageReaction: removeMessageReaction,
		createGroupRoom:       createGroupRoom,
		updateRoom:            updateRoom,
		listRoomMembers:       listRoomMembers,
		addRoomMember:         addRoomMember,
		removeRoomMember:      removeRoomMember,
	}
}
