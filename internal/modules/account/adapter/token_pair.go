package adapter

import (
	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
)

func toTokenPair(accessToken, refreshToken string) *accountv1.TokenPair {
	return &accountv1.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
