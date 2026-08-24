package clipboard

import "github.com/atotto/clipboard"

// WriteFunc performs the actual clipboard write; swappable in tests so
// they never touch the real OS clipboard.
var WriteFunc = clipboard.WriteAll

// Copy writes s to the system clipboard.
func Copy(s string) error {
	return WriteFunc(s)
}
