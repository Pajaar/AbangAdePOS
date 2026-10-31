// Package executor berisi satu executor per use case modul produk. Executor
// hanya bergantung pada kontrak tck/inservice dan common/access.
package executor

import (
	"pos-jamu-backend/internal/common/access"
)

const (
	defaultLimit int32 = 50
)

// canSeeHPP: role boleh menerima field HPP (izin aksi lihat_hpp).
func canSeeHPP(a access.IActionChecker, roleKode string) bool {
	return a.HasAction(roleKode, access.ActionLihatHPP)
}

func limitOrDefault(n int32) int32 {
	if n <= 0 {
		return defaultLimit
	}
	return n
}
