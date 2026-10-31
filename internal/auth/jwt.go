package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken dikembalikan untuk segala kegagalan parse/verifikasi access
// token (signature salah, expired, algoritma tidak dikenal, dll) — detail
// spesifik tidak dibocorkan ke pemanggil, cukup "token tidak valid".
var ErrInvalidToken = errors.New("token tidak valid")

// Claims: sub = user_id, role = kode role. exp/iat dari RegisteredClaims.
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func GenerateAccessToken(userID uint64, roleKode string, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
		Role: roleKode,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseAccessToken memvalidasi signature dan membatasi algoritma ke HS256
// secara eksplisit (mencegah serangan "alg: none" / downgrade algoritma).
func ParseAccessToken(tokenStr string, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (c *Claims) UserID() (uint64, error) {
	return strconv.ParseUint(c.Subject, 10, 64)
}
