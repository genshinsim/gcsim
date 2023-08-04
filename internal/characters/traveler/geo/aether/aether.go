package aether

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common/geo"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.AetherGeo, geo.NewChar(0))
}
