package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type Tmpl struct {
	*core.Core
	*reactable.Reactable
	TargetType  combat.TargettableType
	TargetIndex int
	Hitbox      combat.Circle
	Tags        map[string]int
}
