package xerrors

import "connectrpc.com/connect"

var (
	ErrCannotFollowSelf = New("cannot_follow_self", "cannot follow yourself", connect.CodeInvalidArgument)
)
