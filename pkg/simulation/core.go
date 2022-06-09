package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func NewCore(seed int64, debug bool, cfg *ast.ActionList) (*core.Core, error) {
	return core.New(seed, debug)
}
