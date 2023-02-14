package yoimiya

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Yoimiya, NewChar)
}

type char struct {
	*tmpl.Character
	a1Stack   int
	a1Bonus   []float64
	a4Bonus   []float64
	abApplied bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1Bonus = make([]float64, attributes.EndStatType)
	c.a4Bonus = make([]float64, attributes.EndStatType)
	c.a1()
	c.onExit()
	c.burstHook()
	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}

func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	//infusion to normal attack only
	if c.StatusIsActive(skillKey) && ai.AttackTag == combat.AttackTagNormal {
		ai.Element = attributes.Pyro
		ai.Mult = skill[c.TalentLvlSkill()] * ai.Mult
		c.Core.Log.NewEvent("skill mult applied", glog.LogCharacterEvent, c.Index).
			Write("prev", ai.Mult).
			Write("next", skill[c.TalentLvlSkill()]*ai.Mult).
			Write("char", c.Index)
	}

	return ds
}
