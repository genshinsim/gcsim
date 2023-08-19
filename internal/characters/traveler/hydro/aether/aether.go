package aether

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common/anemo"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.AetherHydro, anemo.NewChar(0))
}
