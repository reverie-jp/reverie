package xerrors

import "connectrpc.com/connect"

var (
	ErrCallNotFound            = New("call_not_found", "call not found", connect.CodeNotFound)
	ErrCallParticipantNotFound = New("call_participant_not_found", "caller is not an active participant of this call", connect.CodeNotFound)
	ErrCallLocked              = New("call_locked", "call is locked and not accepting new participants", connect.CodeFailedPrecondition)
	ErrCallGuestNotAllowed     = New("call_guest_not_allowed", "guests cannot join this call", connect.CodePermissionDenied)
	ErrAlreadyInAnotherCall    = New("already_in_another_call", "user is already participating in another call", connect.CodeFailedPrecondition)
	ErrNotCallHost             = New("not_call_host", "only the host can perform this action", connect.CodePermissionDenied)
)
