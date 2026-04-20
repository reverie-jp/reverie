package adapter

import (
	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func toProviderString(p accountv1.AuthProvider) (string, error) {
	switch p {
	case accountv1.AuthProvider_AUTH_PROVIDER_GOOGLE:
		return "google", nil
	default:
		return "", xerrors.ErrInvalidArgument.WithMessage("unsupported auth provider")
	}
}
