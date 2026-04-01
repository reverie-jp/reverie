package xerrors

import "connectrpc.com/connect"

// Common errors shared across modules.
var (
	ErrInvalidArgument  = New("invalid_argument", "invalid argument", connect.CodeInvalidArgument)
	ErrUnauthenticated  = New("unauthenticated", "authentication required", connect.CodeUnauthenticated)
	ErrNotFound         = New("not_found", "resource not found", connect.CodeNotFound)
	ErrAlreadyExists    = New("already_exists", "resource already exists", connect.CodeAlreadyExists)
	ErrPermissionDenied = New("permission_denied", "permission denied", connect.CodePermissionDenied)
	ErrInternal         = New("internal", "internal error", connect.CodeInternal)
)
