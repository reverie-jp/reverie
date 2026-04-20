package handler

import (
	"context"
	"encoding/base64"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	chatv1 "reverie.jp/reverie/internal/gen/pb/chat/v1"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) ListRooms(ctx context.Context, req *connect.Request[chatv1.ListRoomsRequest]) (*connect.Response[chatv1.ListRoomsResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	outputs, err := h.listRooms.Execute(ctx, requestorID)
	if err != nil {
		return nil, err
	}

	rooms := make([]*chatv1.Room, len(outputs))
	for i, o := range outputs {
		rooms[i] = toProtoRoom(o)
	}

	var nextPageToken string
	if len(outputs) > 0 {
		last := outputs[len(outputs)-1]
		if last.LastMessageAt != nil {
			raw := last.LastMessageAt.UTC().Format("2006-01-02T15:04:05.999999999Z")
			nextPageToken = base64.StdEncoding.EncodeToString([]byte(raw))
		}
	}

	return connect.NewResponse(&chatv1.ListRoomsResponse{
		Rooms:         rooms,
		NextPageToken: nextPageToken,
	}), nil
}

func (h *Handler) CreateDirectRoom(ctx context.Context, req *connect.Request[chatv1.CreateDirectRoomRequest]) (*connect.Response[chatv1.CreateDirectRoomResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	out, err := h.createDirectRoom.Execute(ctx, requestorID, req.Msg.UserId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.CreateDirectRoomResponse{
		Room: toProtoRoom(out),
	}), nil
}

func (h *Handler) MarkRoomAsRead(ctx context.Context, req *connect.Request[chatv1.MarkRoomAsReadRequest]) (*connect.Response[chatv1.MarkRoomAsReadResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	if err := h.markRoomAsRead.Execute(ctx, requestorID, req.Msg.RoomId); err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.MarkRoomAsReadResponse{}), nil
}

func (h *Handler) PinRoom(ctx context.Context, req *connect.Request[chatv1.PinRoomRequest]) (*connect.Response[chatv1.PinRoomResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	out, err := h.pinRoom.Execute(ctx, requestorID, req.Msg.RoomId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.PinRoomResponse{Room: toProtoRoom(out)}), nil
}

func (h *Handler) UnpinRoom(ctx context.Context, req *connect.Request[chatv1.UnpinRoomRequest]) (*connect.Response[chatv1.UnpinRoomResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	out, err := h.unpinRoom.Execute(ctx, requestorID, req.Msg.RoomId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.UnpinRoomResponse{Room: toProtoRoom(out)}), nil
}

func (h *Handler) MuteRoom(ctx context.Context, req *connect.Request[chatv1.MuteRoomRequest]) (*connect.Response[chatv1.MuteRoomResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	out, err := h.muteRoom.Execute(ctx, requestorID, req.Msg.RoomId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.MuteRoomResponse{Room: toProtoRoom(out)}), nil
}

func (h *Handler) UnmuteRoom(ctx context.Context, req *connect.Request[chatv1.UnmuteRoomRequest]) (*connect.Response[chatv1.UnmuteRoomResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	out, err := h.unmuteRoom.Execute(ctx, requestorID, req.Msg.RoomId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.UnmuteRoomResponse{Room: toProtoRoom(out)}), nil
}

func (h *Handler) LeaveRoom(ctx context.Context, req *connect.Request[chatv1.LeaveRoomRequest]) (*connect.Response[chatv1.LeaveRoomResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	if err := h.leaveRoom.Execute(ctx, requestorID, req.Msg.RoomId); err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.LeaveRoomResponse{}), nil
}
