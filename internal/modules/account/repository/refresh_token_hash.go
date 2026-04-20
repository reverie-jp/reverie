package repository

import (
	"crypto/sha256"
	"encoding/hex"
)

// hashRefreshToken returns the SHA-256 hex digest of a raw refresh token
// so that only the digest is persisted. Given a raw token, the digest is
// deterministic and suitable for equality lookup.
func hashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
