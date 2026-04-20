package handler

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reverie.jp/reverie/internal/application/server/interceptor"
	chatv1 "reverie.jp/reverie/internal/gen/pb/chat/v1"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) CreateGroupRoom(ctx context.Context, req *connect.Request[chatv1.CreateGroupRoomRequest]) (*connect.Response[chatv1.CreateGroupRoomResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}
	out, err := h.createGroupRoom.Execute(ctx, requestorID, req.Msg.Name, req.Msg.MemberIds)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&chatv1.CreateGroupRoomResponse{Room: toProtoRoom(out)}), nil
}

func (h *Handler) UpdateRoom(ctx context.Context, req *connect.Request[chatv1.UpdateRoomRequest]) (*connect.Response[chatv1.UpdateRoomResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}
	out, err := h.updateRoom.Execute(ctx, requestorID, req.Msg.RoomId, req.Msg.Name)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&chatv1.UpdateRoomResponse{Room: toProtoRoom(out)}), nil
}

func (h *Handler) ListRoomMembers(ctx context.Context, req *connect.Request[chatv1.ListRoomMembersRequest]) (*connect.Response[chatv1.ListRoomMembersResponse], error) {
	_, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}
	members, err := h.listRoomMembers.Execute(ctx, req.Msg.RoomId)
	if err != nil {
		return nil, err
	}
	protoMembers := make([]*userv1.User, len(members))
	for i, m := range members {
		bio := m.Biography
		protoMembers[i] = &userv1.User{
			Id:          m.ID.String(),
			CustomId:    m.CustomID,
			DisplayName: m.DisplayName,
			Biography:   &bio,
			IsPrivate:   m.IsPrivate,
			CreateTime:  timestamppb.New(m.CreateTime),
		}
	}
	return connect.NewResponse(&chatv1.ListRoomMembersResponse{Members: protoMembers}), nil
}

func (h *Handler) AddRoomMember(ctx context.Context, req *connect.Request[chatv1.AddRoomMemberRequest]) (*connect.Response[chatv1.AddRoomMemberResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}
	out, err := h.addRoomMember.Execute(ctx, requestorID, req.Msg.RoomId, req.Msg.UserId)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&chatv1.AddRoomMemberResponse{Room: toProtoRoom(out)}), nil
}

func (h *Handler) RemoveRoomMember(ctx context.Context, req *connect.Request[chatv1.RemoveRoomMemberRequest]) (*connect.Response[chatv1.RemoveRoomMemberResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}
	_, err := h.removeRoomMember.Execute(ctx, requestorID, req.Msg.RoomId, req.Msg.UserId)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&chatv1.RemoveRoomMemberResponse{}), nil
}
