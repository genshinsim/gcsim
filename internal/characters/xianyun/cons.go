package xianyun

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2Key = "xianyun-c2"
const c4Icd = "xianyun-c4-icd"
const c6Key = "xianyun-c6"

const c2Dur = 15 * 60
const c4IcdDur = 5 * 60
const c6Dur = 16 * 60

var c4Ratio = []float64{0, 0.5, 0.8, 1.5}
var c6Buff = []float64{0, 0.15, 0.35, 0.7}
var c2BuffMod []float64

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

	c2BuffMod = make([]float64, attributes.EndStatType)
	c2BuffMod[attributes.ATKP] = 0.20

	c.a4Max = 18000
	c.a4Ratio = 4
}

func (c *char) c2buff() {
	if c.Base.Cons < 2 {
		return
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c2Key, c2Dur),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return c2BuffMod, true
		},
	})
}

func (c *char) c4cb() func(a combat.AttackCB) {
	if c.Base.Cons < 4 {
		return nil
	}

	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}

		if c.StatusIsActive(c4Icd) {
			return
		}

		atk := c.Base.Atk*(1+c.Stat(attributes.ATKP)) + c.Stat(attributes.ATK)
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Mystery Millet Gourmet (C4)",
			Src:     c4Ratio[c.skillCounter] * atk,
			Bonus:   c.Stat(attributes.Heal),
		})

		c.AddStatus(c4Icd, c4IcdDur, true)
	}
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	c.AddStatus(c6Key, c6Dur, true)
	c.SetTag(c6Key, 8)
}

func (c *char) c6mod(snap *combat.Snapshot) {
	if c.Base.Cons < 6 {
		return
	}
	old := snap.Stats[attributes.CD]
	snap.Stats[attributes.CD] += c6Buff[c.skillCounter]
	c.Core.Log.NewEvent("c6 adding crit DMG", glog.LogCharacterEvent, c.Index).
		Write("old", old).
		Write("new", snap.Stats[attributes.CD])
}
