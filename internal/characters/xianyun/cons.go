package xianyun

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var c2buffMod []float64 = nil

const c4Icd = "xianyun-c4-icd"
const c6key = "xianyun-c6"

const c4IcdDur = 5 * 60
const c6Dur = 16 * 60

var c4ratio = []float64{0, 0.5, 0.8, 1.5}
var c6buff = []float64{0, 0.15, 0.35, 0.7}

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	c.SetNumCharges(action.ActionSkill, 2)
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c2buffMod = make([]float64, attributes.EndStatType)
	c2buffMod[attributes.ATKP] = 0.20

	c.a4Max = 18000
	c.a4Ratio = 4
}

func (c *char) c2buff() {
	if c.Base.Cons < 2 {
		return
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("xianyun-c2", 15*60),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return c2buffMod, true
		},
	})
}

// TODO: C4 Xianyun
func (c *char) c4cb() func(a combat.AttackCB) {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if c.StatusIsActive(c4Icd) {
			return
		}

		atk := c.Base.Atk*(1+c.Stat(attributes.ATKP)) + c.Stat(attributes.ATK)
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Mystery Millet Gourmet (C4)",
			Src:     c4ratio[c.skillCounter] * atk,
			Bonus:   c.Stat(attributes.Heal),
		})

		c.AddStatus(c4Icd, c4IcdDur, true)
	}
}

// TODO: C6 Xianyun
func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	c.AddStatus(c6key, c6Dur, true)
	c.SetTag(c6key, 8)
}

func (c *char) c6mod(snap *combat.Snapshot) {
	if c.Base.Cons < 6 {
		return
	}
	if !c.skillWasC6 {
		return
	}
	old := snap.Stats[attributes.CD]
	snap.Stats[attributes.CD] += c6buff[c.skillCounter]
	c.Core.Log.NewEvent("c6 adding crit DMG", glog.LogCharacterEvent, c.Index).
		Write("old", old).
		Write("new", snap.Stats[attributes.CD])
}

func (c *char) c6cb() func(a combat.AttackCB) {
	if c.Base.Cons < 6 {
		return nil
	}
	return func(a combat.AttackCB) {
		if !c.skillWasC6 {
			return
		}
		c.SetTag(c6key, c.Tag(c6key)-1)
		if c.Tag(c6key) <= 0 {
			c.DeleteStatus(c6key)
		}
		c.skillWasC6 = false
	}
}
