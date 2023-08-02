package lumine

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common/geo"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.LumineGeo, geo.NewChar(1))
}
