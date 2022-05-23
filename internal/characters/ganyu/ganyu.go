package ganyu

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Ganyu, NewChar)
}

type char struct {
	*character.Tmpl
	a1Expiry int
	c4Stacks int
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
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 6
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue
	c.InitCancelFrames()

	c.a1Expiry = -1
	c.c4Stacks = 0

	if c.Base.Cons >= 2 {
		c.SetNumCharges(core.ActionSkill, 2)
	}

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	if c.Base.Cons >= 1 {
		c.c1()
	}
}
