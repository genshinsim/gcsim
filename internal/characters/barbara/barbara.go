package barbara

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("barbara", NewChar)
}

type char struct {
	*character.Tmpl
	stacks int
	c6icd  int
	// burstBuffExpiry   int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassCatalyst
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 4

	c.a1()

	if c.Base.Cons >= 1 {
		c.c1(1)
	}
	c.onSkillStackCount() //doesnt' do anything yet
	if c.Base.Cons >= 6 {
		c.c6()
	}
	return &c, nil
}

func (c *char) a1() {
	c.Core.AddStamMod(func(a core.ActionType) (float64, bool) { // @srl does this activate for the active char?
		if c.Core.Status.Duration("barbara-field") >= 0 {
			return -0.12, false
		}
		return 0, false
	})
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	/*
		Returns character stamina consumption for specified action.
	*/
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 50
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}
func (c *char) c1(delay int) {
	c.AddTask(func() {
		c.AddEnergy(1)
		c.c1(0)
	}, "barbara-c1", delay+10*60)
}

// inspired from hutao c6
func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.Index { //trigger only when not barbara
			c.checkc6()
		}
		return false
	}, "barbara-c6")
}

func (c *char) checkc6() {
	if c.Base.Cons < 6 {
		return
	}
	if c.Core.F < c.c6icd && c.c6icd != 0 {
		return
	}
	//if dead, revive back to 1 hp
	if c.HPCurrent < 0 {
		c.HPCurrent = c.HPMax
	}

	c.c6icd = c.Core.F + 60*60*15
}
