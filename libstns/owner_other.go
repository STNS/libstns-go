//go:build !unix

package libstns

import "os"

func isOwnedByCurrentUser(fi os.FileInfo) bool {
	return true
}
