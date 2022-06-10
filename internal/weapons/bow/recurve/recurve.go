package recurve

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	//TODO: Defeating an opponent restores 8% HP.
	core.RegisterWeaponFunc(keys.RecurveBow, common.NewNoEffect)
}
