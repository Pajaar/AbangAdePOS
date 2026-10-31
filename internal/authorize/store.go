// Package authorize memuat izin ke memori supaya tidak perlu query database
// per request:
//   - endpoint: role_kode -> set("METHOD path"), dipakai middleware Authorize
//   - aksi:     role_kode -> set(kode aksi), misalnya lihat_hpp, dipakai executor
package authorize

import (
	"context"
	"fmt"
	"sync"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/module/auth/tck/inservice"
)

type Store struct {
	mu      sync.RWMutex
	perms   map[string]map[string]struct{}
	actions map[string]map[string]struct{}
}

var _ access.IActionChecker = (*Store)(nil)

func NewStore() *Store {
	return &Store{
		perms:   make(map[string]map[string]struct{}),
		actions: make(map[string]map[string]struct{}),
	}
}

// Load membangun ulang kedua map dari database. Dipanggil sekali saat startup
// (fail fast kalau gagal). Headadmin TIDAK mendapat bypass di sini; izinnya
// datang murni dari data seed role_menus.
func (s *Store) Load(ctx context.Context, svc inservice.IPermissionInService) error {
	perms, err := svc.ListAll(ctx)
	if err != nil {
		return err
	}
	actions, err := svc.ListAllActions(ctx)
	if err != nil {
		return err
	}

	nextPerms := make(map[string]map[string]struct{})
	for _, p := range perms {
		add(nextPerms, p.RoleKode, permKey(p.Method, p.Path))
	}
	nextActions := make(map[string]map[string]struct{})
	for _, a := range actions {
		add(nextActions, a.RoleKode, a.AksiKode)
	}

	s.mu.Lock()
	s.perms = nextPerms
	s.actions = nextActions
	s.mu.Unlock()
	return nil
}

// Reload adalah alias Load, dipanggil ulang setelah hak akses role diubah
// lewat PUT /api/v1/roles/:id/akses.
func (s *Store) Reload(ctx context.Context, svc inservice.IPermissionInService) error {
	return s.Load(ctx, svc)
}

func (s *Store) IsAllowed(roleKode, method, path string) bool {
	return s.has(s.perms, roleKode, permKey(method, path))
}

func (s *Store) HasAction(roleKode, aksiKode string) bool {
	return s.has(s.actions, roleKode, aksiKode)
}

func (s *Store) has(m map[string]map[string]struct{}, roleKode, key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := m[roleKode][key]
	return ok
}

func add(m map[string]map[string]struct{}, role, key string) {
	if m[role] == nil {
		m[role] = make(map[string]struct{})
	}
	m[role][key] = struct{}{}
}

func permKey(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
