package amber

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Amber, NewChar)
}

type char struct {
	*character.Tmpl
	bunnies []bunny
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Pyro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.bunnies = make([]bunny, 0, 2)

	if c.Base.Cons >= 4 {
		c.SetNumCharges(core.ActionSkill, 2)
	}

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.a1()

	if c.Base.Cons >= 2 {
		c.overloadExplode()
	}
}

func (c *char) a1() {
	m := make([]float64, core.EndStatType)
	m[core.CR] = .1

	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "amber-a1",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			return m, atk.Info.AttackTag == core.AttackTagElementalBurst
		},
	})
}

func (c *char) a4(a core.AttackCB) {
	if !a.AttackEvent.Info.HitWeakPoint {
		return
	}

	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.15

	c.AddMod(core.CharStatMod{
		Key:    "amber-a4",
		Expiry: c.Core.F + 600,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}
