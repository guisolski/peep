package clipboard

import (
	"encoding/json"

	"github.com/atotto/clipboard"
	"github.com/guisolski/peep/model"
)

// CopyValue copies the leaf value (or summary) of n to the system clipboard.
func CopyValue(n *model.Node) error {
	return clipboard.WriteAll(n.Summary())
}

// CopyPath copies the JSON path of n (e.g. ".address.city") to the clipboard.
func CopyPath(n *model.Node) error {
	return clipboard.WriteAll(n.Path())
}

// CopySubtree marshals n back to indented JSON and copies it.
func CopySubtree(n *model.Node) error {
	b, err := json.MarshalIndent(n.ToInterface(), "", "  ")
	if err != nil {
		return err
	}
	return clipboard.WriteAll(string(b))
}
