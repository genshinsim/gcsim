package aether

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common/dendro"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.AetherDendro, dendro.NewChar(0))
}
