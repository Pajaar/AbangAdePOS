package auth

import (
	"testing"
	"time"
)

func TestHashAndComparePassword(t *testing.T) {
	hash, err := HashPassword("rahasia123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}

	if err := ComparePassword(hash, "rahasia123"); err != nil {
		t.Errorf("ComparePassword seharusnya cocok: %v", err)
	}

	if err := ComparePassword(hash, "salah"); err == nil {
		t.Error("ComparePassword seharusnya gagal untuk password salah")
	}
}

func TestGenerateAndParseAccessToken(t *testing.T) {
	secret := "test-secret"

	token, err := GenerateAccessToken(42, "headadmin", secret, 15*time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	claims, err := ParseAccessToken(token, secret)
	if err != nil {
		t.Fatalf("ParseAccessToken: %v", err)
	}

	if claims.Role != "headadmin" {
		t.Errorf("role = %q, mau %q", claims.Role, "headadmin")
	}

	userID, err := claims.UserID()
	if err != nil {
		t.Fatalf("UserID: %v", err)
	}
	if userID != 42 {
		t.Errorf("userID = %d, mau 42", userID)
	}
}

func TestParseAccessTokenWrongSecret(t *testing.T) {
	token, err := GenerateAccessToken(1, "cashier", "secret-a", 15*time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	if _, err := ParseAccessToken(token, "secret-b"); err == nil {
		t.Error("ParseAccessToken seharusnya gagal dengan secret yang salah")
	}
}

func TestParseAccessTokenExpired(t *testing.T) {
	token, err := GenerateAccessToken(1, "cashier", "secret", -1*time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	if _, err := ParseAccessToken(token, "secret"); err == nil {
		t.Error("ParseAccessToken seharusnya gagal untuk token expired")
	}
}

func TestGenerateAndHashRefreshToken(t *testing.T) {
	plain, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}
	if plain == "" {
		t.Fatal("GenerateRefreshToken menghasilkan string kosong")
	}

	plain2, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}
	if plain == plain2 {
		t.Error("dua refresh token berturut-turut seharusnya tidak sama")
	}

	hash := HashRefreshToken(plain)
	if len(hash) != 64 {
		t.Errorf("panjang hash = %d, mau 64", len(hash))
	}

	hashAgain := HashRefreshToken(plain)
	if hash != hashAgain {
		t.Error("HashRefreshToken seharusnya deterministik untuk input yang sama")
	}
}
