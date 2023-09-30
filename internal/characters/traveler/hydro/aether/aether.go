package aether

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common/hydro"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.AetherHydro, hydro.NewChar(0))
}
