package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// GenerateRefreshToken membuat string acak (32 byte, base64 URL-safe) yang
// dikirim ke klien lewat cookie httpOnly. Hanya hash-nya yang disimpan di DB.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashRefreshToken menghasilkan SHA-256 hex (64 karakter, cocok dengan
// kolom refresh_tokens.token_hash CHAR(64)).
func HashRefreshToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
