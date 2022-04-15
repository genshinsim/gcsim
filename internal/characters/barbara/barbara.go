package barbara

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Barbara, NewChar)
}

type char struct {
	*character.Tmpl
	c6icd      int
	skillInitF int
	// burstBuffExpiry   int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Hydro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassCatalyst
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 4

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.a1()

	if c.Base.Cons >= 1 {
		c.c1(1)
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}
}

func (c *char) a1() {
	c.Core.AddStamMod(func(a core.ActionType) (float64, bool) { // @srl does this activate for the active char?
		if c.Core.Status.Duration("barbskill") >= 0 {
			return -0.12, false
		}
		return 0, false
	}, "barb-a1-stam")
}

func (c *char) c2() {
	val := make([]float64, core.EndStatType)
	val[core.HydroP] = 0.15
	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue
		}
		char.AddMod(core.CharStatMod{
			Key:    "barbara-c2",
			Expiry: -1,
			Amount: func() ([]float64, bool) {
				if c.Core.Status.Duration("barbskill") >= 0 {
					return val, true
				} else {
					return nil, false
				}
			},
		})
	}
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
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}
func (c *char) c1(delay int) {
	c.AddTask(func() {
		c.AddEnergy("barbara-c1", 1)
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
	if c.HP() <= -1 {
		c.HPCurrent = c.MaxHP()
	}

	c.c6icd = c.Core.F + 60*60*15
}
