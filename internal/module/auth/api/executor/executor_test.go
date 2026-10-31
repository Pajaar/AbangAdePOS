package executor

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"pos-jamu-backend/internal/auth"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/auth/tck/domain"
	"pos-jamu-backend/internal/util/apperror"
)

// Executor hanya bergantung pada interface tck/inservice, jadi bisa diuji
// dengan implementasi palsu tanpa database.

type fakeUsers struct {
	user domain.User
	err  error
}

func (f fakeUsers) Authenticate(context.Context, string, string) (domain.User, error) {
	return f.user, f.err
}
func (f fakeUsers) GetByID(context.Context, uint64) (domain.User, error) { return f.user, f.err }
func (f fakeUsers) Create(context.Context, domain.NewUser) (uint64, error) {
	return 0, errors.New("not used")
}

type fakeSessions struct {
	startCalled bool
	gotClient   domain.ClientInfo
	session     domain.Session
	err         error
}

func (f *fakeSessions) Start(_ context.Context, u domain.User, c domain.ClientInfo) (domain.Session, error) {
	f.startCalled = true
	f.gotClient = c
	return f.session, f.err
}
func (f *fakeSessions) Rotate(_ context.Context, _ string, c domain.ClientInfo) (domain.Session, error) {
	f.gotClient = c
	return f.session, f.err
}
func (f *fakeSessions) End(context.Context, string) error { return f.err }

var testAT = AccessTokenConfig{Secret: "test-secret", TTL: 15 * time.Minute}

func TestLoginExecutor_Success(t *testing.T) {
	user := domain.User{ID: 7, Nama: "Ade", Username: "ade", RoleKode: "cashier", Aktif: true}
	sessions := &fakeSessions{session: domain.Session{UserID: 7, RoleKode: "cashier", RefreshToken: "rt-plain", RefreshExpiresAt: time.Now().Add(time.Hour)}}
	ex := NewLoginExecutor(fakeUsers{user: user}, sessions, testAT)

	res, err := ex.Process(context.Background(), dto.LoginRequest{Username: "ade", Password: "x", UserAgent: "ua", IPAddress: "1.2.3.4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.RefreshToken != "rt-plain" {
		t.Errorf("refresh token = %q", res.RefreshToken)
	}
	if res.Response.User.Role != "cashier" || res.Response.User.ID != 7 {
		t.Errorf("user info = %+v", res.Response.User)
	}
	if sessions.gotClient.UserAgent != "ua" || sessions.gotClient.IPAddress != "1.2.3.4" {
		t.Errorf("client info tidak diteruskan: %+v", sessions.gotClient)
	}
	claims, err := auth.ParseAccessToken(res.Response.AccessToken, testAT.Secret)
	if err != nil || claims.Role != "cashier" {
		t.Errorf("access token tidak valid: %v %+v", err, claims)
	}
}

func TestLoginExecutor_FailedDoesNotStartSession(t *testing.T) {
	sessions := &fakeSessions{}
	ex := NewLoginExecutor(fakeUsers{err: apperror.Unauthorized("Username atau password salah")}, sessions, testAT)

	_, err := ex.Process(context.Background(), dto.LoginRequest{Username: "x", Password: "y"})
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Status != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %v", err)
	}
	if sessions.startCalled {
		t.Error("session tidak boleh dibuat saat login gagal")
	}
}

func TestRefreshExecutor_PassesErrorThrough(t *testing.T) {
	ex := NewRefreshExecutor(&fakeSessions{err: apperror.Unauthorized("Sesi tidak valid")}, testAT)
	_, err := ex.Process(context.Background(), dto.RefreshRequest{RefreshToken: "old"})
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Status != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %v", err)
	}
}

func TestMeExecutor_NotFoundBecomesUnauthorized(t *testing.T) {
	ex := NewMeExecutor(fakeUsers{err: apperror.NotFound("User tidak ditemukan")})
	_, err := ex.Process(context.Background(), dto.MeRequest{UserID: 1})
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Status != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %v", err)
	}
}
