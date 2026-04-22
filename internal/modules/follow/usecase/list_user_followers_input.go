package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type ListUserFollowersInput struct {
	RequesterID    ulid.ULID
	TargetCustomID string
	PageSize       int32
	PageToken      string
}

func (i ListUserFollowersInput) Validate() error {
	if i.TargetCustomID == "" {
		return errors.New("target_custom_id is required")
	}
	return nil
}
