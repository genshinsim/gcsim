package heizou

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Heizou, NewChar)
}

type char struct {
	*character.Tmpl
	decStack            int
	infuseCheckLocation core.AttackPattern
	a1icd               int
	c1icd               int
	c1buff              []float64
	a4buff              []float64
	burstTaggedCount    int
}

const (
	eCD = 600
)

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Anemo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 5
	c.SkillCon = 3
	c.BurstCon = 5
	c.InitCancelFrames()
	c.a1icd = 0
	c.c1icd = 0

	c.infuseCheckLocation = core.NewDefCircHit(0.1, false, core.TargettableEnemy, core.TargettablePlayer, core.TargettableObject)

	if c.Base.Cons >= 1 {
		c.c1()
	}

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.a1()

	c.a4buff = make([]float64, core.EndStatType)
	c.a4buff[core.EM] = 80

	if c.Base.Cons >= 1 {
		c.c1buff = make([]float64, core.EndStatType)
		c.c1buff[core.AtkSpd] = .15
	}

}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		return 0
	}
}
