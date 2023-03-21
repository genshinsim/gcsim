package travelergeo

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
	core.RegisterCharFunc(keys.AetherGeo, NewChar(0))
	core.RegisterCharFunc(keys.LumineGeo, NewChar(1))
}

type char struct {
	*tmpl.Character
	skillCD     int
	burstArea   combat.AttackPattern // needed for c1
	c1TickCount int
	gender      int
}

func NewChar(gender int) core.NewCharacterFunc {
	return func(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
		c := char{
			gender: gender,
		}
		c.Character = tmpl.NewWithWrapper(s, w)

		c.Base.Element = attributes.Geo
		c.EnergyMax = 60
		c.BurstCon = 3
		c.SkillCon = 5
		c.NormalHitNum = normalHitNum

		c.skillCD = 6 * 60

		w.Character = &c

		return nil
	}
}

func (c *char) Init() error {
	c.a1()
	// setup number of C1 ticks
	c.c1TickCount = 15
	if c.Base.Cons >= 6 {
		c.c1TickCount = 20
	}
	return nil
}
