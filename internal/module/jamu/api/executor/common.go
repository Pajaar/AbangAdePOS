// Package executor berisi satu executor per use case modul menu jamu.
package executor

import "pos-jamu-backend/internal/common/access"

func canSeeHPP(a access.IActionChecker, roleKode string) bool {
	return a.HasAction(roleKode, access.ActionLihatHPP)
}
