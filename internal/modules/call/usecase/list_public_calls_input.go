package usecase

import "reverie.jp/reverie/internal/platform/ulid"

type ListPublicCallsInput struct {
	// Zero means the caller is a guest. Authenticated callers additionally
	// see USERS_ONLY calls.
	RequesterID ulid.ULID
	// <= 0 means default (50). Clamped to max (100).
	PageSize int32
	// Empty means first page. Otherwise the ULID cursor from a prior response.
	PageToken string
}
