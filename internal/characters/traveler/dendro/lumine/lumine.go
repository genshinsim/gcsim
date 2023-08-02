package lumine

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common/dendro"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.LumineDendro, dendro.NewChar(1))
}
