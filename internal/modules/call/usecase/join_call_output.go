package usecase

import "time"

type JoinCallOutput struct {
	AccessToken string
	URL         string
	Identity    string
	ExpireTime  time.Time
}
