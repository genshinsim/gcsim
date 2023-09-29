package lumine

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common/hydro"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.LumineHydro, hydro.NewChar(1))
}
