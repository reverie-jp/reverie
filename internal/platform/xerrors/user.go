package xerrors

import "connectrpc.com/connect"

var (
	ErrUserNotFound = New("user_not_found", "user not found", connect.CodeNotFound)
)
