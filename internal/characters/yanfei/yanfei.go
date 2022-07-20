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
	sealCount         int
	burstBuff         []float64
	a1buff            []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
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
	c.a1buff = make([]float64, attributes.EndStatType)
	c.burstBuff = make([]float64, attributes.EndStatType)
	c.burstBuff[attributes.DmgP] = burstBonus[c.TalentLvlBurst()]
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
		if !c.StatusIsActive(sealBuffKey) {
			c.sealCount = 0
		}
		return 50 * (1 - c.sealStamReduction*float64(c.sealCount))
	}
	return c.Character.ActionStam(a, p)
}

// Hook that clears yanfei burst status and seals when she leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		c.sealCount = 0
		c.DeleteStatus(sealBuffKey)
		c.Core.Status.Delete("yanfeiburst")
		return false
	}, "yanfei-exit")
}
