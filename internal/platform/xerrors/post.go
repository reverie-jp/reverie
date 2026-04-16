package xerrors

import "connectrpc.com/connect"

var (
	ErrPostNotFound = New("post_not_found", "post not found", connect.CodeNotFound)
)
