package usecase

import (
	"context"

	chatgw "reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type MarkRoomAsRead struct{ gateway chatgw.Gateway }
type PinRoom struct{ gateway chatgw.Gateway }
type UnpinRoom struct{ gateway chatgw.Gateway }
type MuteRoom struct{ gateway chatgw.Gateway }
type UnmuteRoom struct{ gateway chatgw.Gateway }
type LeaveRoom struct{ gateway chatgw.Gateway }

func NewMarkRoomAsRead(g chatgw.Gateway) *MarkRoomAsRead { return &MarkRoomAsRead{gateway: g} }
func NewPinRoom(g chatgw.Gateway) *PinRoom               { return &PinRoom{gateway: g} }
func NewUnpinRoom(g chatgw.Gateway) *UnpinRoom           { return &UnpinRoom{gateway: g} }
func NewMuteRoom(g chatgw.Gateway) *MuteRoom             { return &MuteRoom{gateway: g} }
func NewUnmuteRoom(g chatgw.Gateway) *UnmuteRoom         { return &UnmuteRoom{gateway: g} }
func NewLeaveRoom(g chatgw.Gateway) *LeaveRoom           { return &LeaveRoom{gateway: g} }

func (uc *MarkRoomAsRead) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr string) error {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return xerrors.ErrInvalidArgument
	}
	return uc.gateway.MarkRoomAsRead(ctx, roomID, requestorID)
}

func (uc *PinRoom) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr string) (*RoomOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.PinRoom(ctx, requestorID, roomID)
	if err != nil {
		return nil, err
	}
	return toRoomOutput(view, requestorID), nil
}

func (uc *UnpinRoom) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr string) (*RoomOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.UnpinRoom(ctx, requestorID, roomID)
	if err != nil {
		return nil, err
	}
	return toRoomOutput(view, requestorID), nil
}

func (uc *MuteRoom) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr string) (*RoomOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.MuteRoom(ctx, requestorID, roomID)
	if err != nil {
		return nil, err
	}
	return toRoomOutput(view, requestorID), nil
}

func (uc *UnmuteRoom) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr string) (*RoomOutput, error) {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument
	}
	view, err := uc.gateway.UnmuteRoom(ctx, requestorID, roomID)
	if err != nil {
		return nil, err
	}
	return toRoomOutput(view, requestorID), nil
}

func (uc *LeaveRoom) Execute(ctx context.Context, requestorID ulid.ULID, roomIDStr string) error {
	roomID, err := ulid.Parse(roomIDStr)
	if err != nil {
		return xerrors.ErrInvalidArgument
	}
	return uc.gateway.LeaveRoom(ctx, roomID, requestorID)
}
