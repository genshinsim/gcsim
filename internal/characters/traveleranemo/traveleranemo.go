package traveleranemo

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl
	qInfuse             core.EleType
	qICDTag             core.ICDTag
	eInfuse             core.EleType
	eICDTag             core.ICDTag
	infuseCheckLocation core.AttackPattern
}

func init() {
	core.RegisterCharFunc(core.TravelerAnemo, NewChar)
}

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
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5
	c.InitCancelFrames()
	c.infuseCheckLocation = core.NewDefCircHit(0.1, false, core.TargettableEnemy, core.TargettablePlayer, core.TargettableObject)

	if c.Base.Cons >= 2 {
		c.c2()
	}

	return &c, nil
}
