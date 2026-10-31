// Package authorize memuat permission role->endpoint ke memori supaya
// middleware Authorize tidak perlu query database per request.
package authorize

import (
	"context"
	"fmt"
	"sync"

	"pos-jamu-backend/internal/module/auth/tck/inservice"
)

type Store struct {
	mu    sync.RWMutex
	perms map[string]map[string]struct{}
}

func NewStore() *Store {
	return &Store{perms: make(map[string]map[string]struct{})}
}

// Load membangun ulang map role_kode -> set("METHOD path") dari database.
// Dipanggil sekali saat startup (fail fast kalau gagal). Headadmin TIDAK
// mendapat bypass di sini — izinnya datang murni dari data seed role_menus.
func (s *Store) Load(ctx context.Context, svc inservice.IPermissionInService) error {
	rows, err := svc.ListAll(ctx)
	if err != nil {
		return err
	}

	next := make(map[string]map[string]struct{})
	for _, row := range rows {
		if next[row.RoleKode] == nil {
			next[row.RoleKode] = make(map[string]struct{})
		}
		next[row.RoleKode][permKey(row.Method, row.Path)] = struct{}{}
	}

	s.mu.Lock()
	s.perms = next
	s.mu.Unlock()

	return nil
}

// Reload adalah alias Load, dipanggil ulang setelah hak akses role diubah
// lewat PUT /api/v1/roles/:id/akses (endpoint itu di luar lingkup Tahap 1).
func (s *Store) Reload(ctx context.Context, svc inservice.IPermissionInService) error {
	return s.Load(ctx, svc)
}

func (s *Store) IsAllowed(roleKode, method, path string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	set, ok := s.perms[roleKode]
	if !ok {
		return false
	}
	_, ok = set[permKey(method, path)]
	return ok
}

func permKey(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
