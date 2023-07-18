package electro

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

type char struct {
	*tmpl.Character
	abundanceAmulets      int
	burstC6Hits           int
	burstC6WillGiveEnergy bool
	burstSnap             combat.Snapshot
	burstAtk              *combat.AttackEvent
	burstSrc              int
	gender                int
}

func NewChar(gender int) core.NewCharacterFunc {
	return func(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
		c := char{
			gender: gender,
		}
		c.Character = tmpl.NewWithWrapper(s, w)

		c.Base.Element = attributes.Electro
		c.EnergyMax = 80
		c.BurstCon = 3
		c.SkillCon = 5
		c.NormalHitNum = normalHitNum

		w.Character = &c

		return nil
	}
}

func (c *char) Init() error {
	c.burstProc()
	return nil
}
