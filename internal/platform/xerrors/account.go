package xerrors

import "connectrpc.com/connect"

var (
	ErrAccountNotFound     = New("account_not_found", "account not found", connect.CodeNotFound)
	ErrCustomIDMismatch    = New("custom_id_mismatch", "custom id does not match", connect.CodeInvalidArgument)
	ErrInvalidRefreshToken = New("invalid_refresh_token", "invalid refresh token", connect.CodeUnauthenticated)
	ErrSocialLoginFailed   = New("social_login_failed", "failed to authenticate with provider", connect.CodeUnauthenticated)
)
