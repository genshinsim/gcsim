package neuvillette

import (
	"strings"

	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4ICDKey = "neuvillette-c4-icd"
const c6ICDKey = "neuvillette-c6-icd"

func (c *char) c1() {
	if c.Base.Ascension < 1 {
		return
	}

	c1 := NeuvA1Keys{event.OnCharacterSwap, "neuvillette-a1-c1-onfield"}
	c.a1Statuses = append(c.a1Statuses, c1)
	c.Core.Events.Subscribe(c1.Evt, func(args ...interface{}) bool {
		next := args[1].(int)
		if next == c.Index {
			c.AddStatus(c1.Key, 30*60, true)
		}
		return false
	}, c1.Key)
}

func (c *char) c2() {
	if c.Base.Ascension < 1 {
		return
	}
	c2Buff := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("neuvillette-c2", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if strings.Contains(atk.Info.Abil, chargeJudgementName) {
				c2Buff[attributes.CD] = 0.14 * float64(c.countA1())
				return c2Buff, true
			}
			return nil, false
		},
	})
}

func (c *char) c4() {
	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		target := args[1].(int)

		if c.Core.Player.Active() != c.Index {
			return false
		}
		if c.Index != target {
			return false
		}
		if c.StatusIsActive(c4ICDKey) {
			return false
		}

		// 4s CD
		c.AddStatus(c4ICDKey, 4*60, true)
		player := c.Core.Combat.Player()
		common.NewSourcewaterDroplet(
			c.Core,
			geometry.CalcRandomPointFromCenter(
				geometry.CalcOffsetPoint(
					player.Pos(),
					geometry.Point{Y: 8},
					player.Direction(),
				),
				1.3,
				3,
				c.Core.Rand,
			),
			combat.GadgetTypSourcewaterDropletNeuv,
		)
		c.Core.Combat.Log.NewEvent("C4: Spawned 1 droplet", glog.LogCharacterEvent, c.Index)

		return false
	}, "neuvillette-c4")
}

func (c *char) c6DropletCheck(src int) func() {
	return func() {
		if c.chargeJudgeStartF != src {
			return
		}

		if c.Core.F > c.chargeJudgeStartF+c.tickAnimLength {
			return
		}

		if c.chargeJudgeStartF+c.chargeJudgeDur-c.Core.F <= 60 {
			droplets := c.getSourcewaterDropletsC6()

			// c6 only absorbs one droplet at a time
			if len(droplets) > 0 {
				c.Core.Combat.Log.NewEvent("C6: Picked up 1 droplet", glog.LogCharacterEvent, c.Index).
					Write("prev-charge-duration", c.chargeJudgeDur).
					Write("curr-charge-duration", c.chargeJudgeDur+60)

				// take first droplet
				c.consumeDroplet(droplets[0])
				c.chargeJudgeDur += 60
			}
		}

		c.QueueCharTask(c.c6DropletCheck(src), 18)
	}
}

func (c *char) c6cb(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(c6ICDKey) {
		return
	}
	c.AddStatus(c6ICDKey, 2*60, true)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       chargeJudgementName + " (C6)",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNeuvilletteC6,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    0.1 * c.MaxHP() * a1Multipliers[c.countA1()],
	}
	// C6 projectile stops on first target hit, with 0.5 rad sphere hitbox.
	// Because we don't simulate the projectile, it's just a circle hit
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5)
	// it looks like the c6 has 29 frames of delay but I didn't count it rigourously
	c.Core.QueueAttack(ai, ap, 29, 29)
	c.Core.QueueAttack(ai, ap, 29, 29)
}
