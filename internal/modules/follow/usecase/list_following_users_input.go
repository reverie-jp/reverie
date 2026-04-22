package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

const (
	defaultListPageSize int32 = 50
	maxListPageSize     int32 = 100
)

type ListFollowingUsersInput struct {
	RequesterID    ulid.ULID
	TargetCustomID string
	PageSize       int32
	PageToken      string
}

func (i ListFollowingUsersInput) Validate() error {
	if i.TargetCustomID == "" {
		return errors.New("target_custom_id is required")
	}
	return nil
}
