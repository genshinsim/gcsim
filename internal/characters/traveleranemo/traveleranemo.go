package traveleranemo

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.AetherAnemo, NewChar(0))
	core.RegisterCharFunc(keys.LumineAnemo, NewChar(1))
}

type char struct {
	*tmpl.Character
	qInfuse             attributes.Element
	qICDTag             combat.ICDTag
	eAbsorb             attributes.Element
	eICDTag             combat.ICDTag
	absorbCheckLocation combat.AttackPattern
	gender              int
}

func NewChar(gender int) core.NewCharacterFunc {
	return func(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
		c := char{
			gender: gender,
		}
		c.Character = tmpl.NewWithWrapper(s, w)

		c.Base.Element = attributes.Anemo
		c.EnergyMax = 60
		c.BurstCon = 3
		c.SkillCon = 5
		c.NormalHitNum = normalHitNum
		c.absorbCheckLocation = combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableGadget)

		w.Character = &c

		return nil
	}
}

func (c *char) Init() error {
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}
