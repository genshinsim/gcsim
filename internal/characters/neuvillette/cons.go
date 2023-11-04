package neuvillette

import (
	"strings"

	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	if c.Base.Ascension < 1 {
		return
	}

	c1 := NeuvA1Keys{event.OnCharacterSwap, "neuvillette-a1-c1-onfield"}
	c.a1Statuses = append(c.a1Statuses, c1)
	c.Core.Events.Subscribe(c1.Evt, func(args ...interface{}) bool {
		next := args[1].(int)
		if next == c.Index {
			c.AddStatus(c1.Key, 30*60, false)
		}
		return false
	}, c1.Key)
}

func (c *char) c2() {
	if c.Base.Ascension < 1 {
		return
	}
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("neuvillette-c2", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if strings.Contains(atk.Info.Abil, chargeJudgementName) {
				c2Buff := make([]float64, attributes.EndStatType)
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

		// 4s CD
		if c.Core.F < c.lastc4+4*60 {
			return false
		}
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if c.Index != target {
			return false
		}
		if c.Core.Player.Active() == c.Index && c.Index == target {
			// TODO: find the actual sourcewater droplet spawn shape for Neuv C4
			player := c.Core.Combat.Player()
			center := player.Pos().Add(player.Direction().Normalize().Mul(geometry.Point{X: 3.0, Y: 3.0}))
			pos := geometry.CalcRandomPointFromCenter(center, 0, 2.5, c.Core.Rand)
			common.NewSourcewaterDroplet(c.Core, pos, combat.GadgetTypSourcewaterDropletNeuv)
			c.lastc4 = c.Core.F
		}
		return false
	}, "neuvillette-c4")
}

func (c *char) c6DropletCheck(src int) func() {
	return func() {
		if c.chargeJudgeStartF != src {
			return
		}

		if c.chargeJudgeStartF+c.chargeJudgeDur-c.Core.F <= 60 {
			droplets := c.getSourcewaterDroplets()

			// c6 only absorbs one droplet at a time
			if len(droplets) > 0 {
				c.consumeDroplet(droplets[c.Core.Combat.Rand.Intn(len(droplets))])
				c.chargeJudgeDur += 60
			}
		}

		c.QueueCharTask(c.c6DropletCheck(src), 18)
	}
}

func (c *char) c6(src int) func() {
	return func() {
		if c.chargeJudgeStartF != src {
			return
		}
		if c.Core.F > c.chargeJudgeStartF+c.chargeJudgeDur {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       chargeJudgementName + " (C6)",
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagNeuvilletteC6,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    0.1 * c.MaxHP() * a1Multipliers[c.countA1()],
		}
		// C6 projectile stops on first target hit, with 0.5 rad sphere hitbox.
		// Because we don't simulate the projectile, it's just a circle hit
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{}, 0.5)
		// it looks like the c6 has 29 frames of delay but I didn't count it rigourously
		c.Core.QueueAttack(ai, ap, 29, 29)
		c.Core.QueueAttack(ai, ap, 29, 29)
		c.QueueCharTask(c.c6(src), 120)
	}
}
