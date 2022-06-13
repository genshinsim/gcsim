package yoimiya

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Yoimiya, NewChar)
}

type char struct {
	*tmpl.Character
	a1stack  int
	lastPart int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Pyro
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassBow
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = character.ZoneInazuma

	w.Character = &c

	return nil
}

func (c *char) Init() error {
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
	if c.Core.Status.Duration("yoimiyaskill") > 0 && ai.AttackTag == combat.AttackTagNormal {
		ai.Element = attributes.Pyro
		ai.Mult = skill[c.TalentLvlSkill()] * ai.Mult
		c.Core.Log.NewEvent("skill mult applied", glog.LogCharacterEvent, c.Index, "prev", ai.Mult, "next", skill[c.TalentLvlSkill()]*ai.Mult, "char", c.Index)
	}

	return ds
}
