package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword menghasilkan bcrypt hash untuk disimpan di users.password_hash.
// Jangan pernah simpan password asli.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword mengembalikan nil kalau password cocok dengan hash.
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
