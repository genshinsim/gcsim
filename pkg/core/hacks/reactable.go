package hacks

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

// TODO: place holder function for creating a reactable until we move off
// old reactable system entirely
func NewReactable(t info.Target, c *core.Core) info.Reactable {
	r := &reactable.Reactable{}
	r.Init(t, c)
	return r
}
