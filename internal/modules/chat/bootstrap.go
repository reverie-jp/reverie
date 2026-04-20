package chat

import (
	"reverie.jp/reverie/internal/gen/pb/chat/v1/chatv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/modules/chat/handler"
	"reverie.jp/reverie/internal/modules/chat/usecase"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

func InitModule(q sqlc.Querier, userGateway usergw.Gateway) chatv1connect.ChatServiceHandler {
	gateway := chatgw.New(q, userGateway)
	return handler.New(
		usecase.NewListRooms(gateway),
		usecase.NewCreateDirectRoom(gateway),
		usecase.NewListMessages(gateway),
		usecase.NewSendMessage(gateway),
		usecase.NewMarkRoomAsRead(gateway),
		usecase.NewPinRoom(gateway),
		usecase.NewUnpinRoom(gateway),
		usecase.NewMuteRoom(gateway),
		usecase.NewUnmuteRoom(gateway),
		usecase.NewLeaveRoom(gateway),
		usecase.NewAddMessageReaction(gateway),
		usecase.NewRemoveMessageReaction(gateway),
		usecase.NewCreateGroupRoom(gateway),
		usecase.NewUpdateRoom(gateway),
		usecase.NewListRoomMembers(gateway),
		usecase.NewAddRoomMember(gateway),
		usecase.NewRemoveRoomMember(gateway),
	)
}
