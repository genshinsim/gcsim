package yanfei

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Yanfei, NewChar)
}

type char struct {
	*tmpl.Character
	maxTags           int
	sealStamReduction float64
	sealExpiry        int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Pyro
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassCatalyst
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneLiyue

	c.maxTags = 3
	if c.Base.Cons >= 6 {
		c.maxTags = 4
	}

	c.sealStamReduction = 0.15
	if c.Base.Cons >= 1 {
		c.sealStamReduction = 0.25
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.onExitField()
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		if c.Core.F > c.sealExpiry {
			c.Tags["seal"] = 0
		}
		stacks := c.Tags["seal"]
		return 50 * (1 - c.sealStamReduction*float64(stacks))
	}
	return c.Character.ActionStam(a, p)
}

// Hook that clears yanfei burst status and seals when she leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		c.Tags["seal"] = 0
		c.sealExpiry = c.Core.F - 1
		c.Core.Status.Delete("yanfeiburst")
		return false
	}, "yanfei-exit")
}
