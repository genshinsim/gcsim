package ayaka

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl

	icdC1          int
	c6CDTimerAvail bool // Flag that controls whether the 0.5 C6 CD timer is available to be started
}

func init() {
	core.RegisterCharFunc(core.Ayaka, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSword
	c.CharZone = core.ZoneInazuma
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	c.icdC1 = -1
	c.c6CDTimerAvail = false

	// Start with C6 ability active
	if c.Base.Cons == 6 {
		c.c6CDTimerAvail = true
	}

	c.InitCancelFrames()

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	// Start with C6 ability active
	if c.Base.Cons == 6 {
		c.c6AddBuff()
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		f, ok := p["f"]
		if !ok {
			return 10 //tap = 36 frames, so under 1 second
		}
		//for every 1 second passed, consume extra 15
		extra := f / 60
		return float64(10 + 15*extra)
	case core.ActionCharge:
		return 20
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}
