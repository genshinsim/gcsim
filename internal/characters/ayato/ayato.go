package ayato

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Ayato, NewChar)
}

type char struct {
	*tmpl.Character
	stacks            int
	stacksMax         int
	shunsuikenCounter int
	c6Ready           bool
}

const (
	particleICDKey = "ayato-particle-icd"
)

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum

	c.shunsuikenCounter = 3

	c.stacksMax = 4
	if c.Base.Cons >= 2 {
		c.stacksMax = 5
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.onExitField()
	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}

func (c *char) AdvanceNormalIndex() {
	c.NormalCounter++

	if c.StatusIsActive(skillBuffKey) {
		if c.NormalCounter == c.shunsuikenCounter {
			c.NormalCounter = 0
		}
	} else {
		if c.NormalCounter == c.NormalHitNum {
			c.NormalCounter = 0
		}
	}
}

// TODO: maybe move infusion out of snapshot?
func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.StatusIsActive(skillBuffKey) {
		switch ai.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		default:
			return ds
		}
		// namisen doesn't affect c6
		if ai.Abil == c6Abil {
			return ds
		}
		ai.Element = attributes.Hydro
		// add namisen stack
		flatdmg := c.MaxHP() * skillpp[c.TalentLvlSkill()] * float64(c.stacks)
		ai.FlatDmg += flatdmg
		c.Core.Log.NewEvent("namisen add damage", glog.LogCharacterEvent, c.Index).
			Write("damage_added", flatdmg).
			Write("stacks", c.stacks).
			Write("expiry", c.StatusExpiry(skillBuffKey))
	}
	return ds
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		if c.StatusIsActive(skillBuffKey) {
			return 17
		}
		return 15
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
