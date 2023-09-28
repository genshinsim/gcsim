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
				c2Buff[attributes.DmgP] = 0.14 * float64(c.countA1())
				return c2Buff, true
			}
			return nil, false
		},
	})
}

func (c *char) c4() {
	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		target := args[1].(int)
		if c.Core.Player.ActiveChar().Index == c.Index && c.Index == target {
			// TODO: find the actual sourcewater droplet spawn shape for Neuv C4
			player := c.Core.Combat.Player()
			center := player.Pos().Add(player.Direction().Normalize().Mul(geometry.Point{X: 3.0, Y: 3.0}))
			pos := geometry.CalcRandomPointFromCenter(center, 0, 2.5, c.Core.Rand)
			common.NewSourcewaterDroplet(c.Core, pos)
		}
		return false
	}, "neuvillette-c4")
}

func (c *char) c6DropletCheck() {
	// I think the c6 droplet check should happen continuously,
	// but currently we cannot modify the duration of an action while it is happening.
	// so right now the c6 check only happens once at the start

	// if c.Core.F-c.chargeSrc >= c.chargeJudgementDur {
	// 	return
	// }
	droplets := make([]*common.SourcewaterDroplet, 0)
	for _, g := range c.Core.Combat.Gadgets() {
		droplet, ok := g.(*common.SourcewaterDroplet)
		if !ok {
			continue
		}
		if droplet.Pos().Distance(c.Core.Combat.Player().Pos()) <= 8 {
			droplets = append(droplets, droplet)
		}
	}

	for _, g := range droplets {
		g.Kill()
		c.healWithDroplets()
		c.chargeJudgementDur += 60
	}
	// c.QueueCharTask(c.c6DropletCheck, 30)
}

func (c *char) c6cb(atk combat.AttackCB) {
	if c.Core.F-c.lastc6 >= 120 {
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
		ap := combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{}, 3, 15)
		c.Core.QueueAttack(ai, ap, 0, 0)
		c.Core.QueueAttack(ai, ap, 0, 0)
	}
}
