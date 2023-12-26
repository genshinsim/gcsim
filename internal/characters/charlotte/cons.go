package charlotte

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1StatusKey     = "charlotte-c1"
	c1HealMsg       = "Verification"
	c6HealMsg       = "charlotte-c6-heal"
	c6CoordinateAtk = "charlotte-c6-coordinate-atk"
	c6IcdKey        = "charlotte-c6-icd"
	c6Radius        = 2
)

func (c *char) c1Heal() func() {
	return func() {
		if c.Core.Status.Duration(c1StatusKey) > 0 {
			stats, _ := c.Stats()
			atk := (c.Base.Atk+c.Weapon.BaseAtk)*(1+stats[attributes.ATKP]) + stats[attributes.ATK]
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: c1HealMsg,
				Src:     atk * 0.8,
				Bonus:   stats[attributes.Heal],
			})
			c.Core.Tasks.Add(c.c1Heal(), 2*60)
		}
	}
}

func (c *char) c1() {
	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		src := args[0].(*player.HealInfo)
		if src.Message == healInitialMsg || src.Message == healDotMsg || src.Message == c6HealMsg {
			if c.Core.Status.Duration(c1StatusKey) == 0 {
				c.Core.Tasks.Add(c.c1Heal(), 2*60)
			}
			c.Core.Status.Add(c1StatusKey, 60*6)
		}
		return false
	}, "charlotte-c1")
}
func (c *char) c2(ap combat.AttackPattern) {
	enemies := c.Core.Combat.EnemiesWithinArea(ap, nil)
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.1 * float64(min(len(enemies), 3))
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("charlotte-c2", 12*60),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (c *char) c4() {
	counter := 0
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		ae := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if counter == 5 {
			return false
		}
		if !t.StatusIsActive(skillPressMarkKey) && !t.StatusIsActive(skillHoldMarkKey) {
			return false
		}
		if ae.Info.Abil != burstInitialAbil && ae.Info.Abil != burstDotAbil {
			return false
		}
		if counter == 0 {
			c.Core.Tasks.Add(func() {
				counter = 0
			}, 1200)
		}
		ae.Snapshot.Stats[attributes.DmgP] += 0.1
		c.AddEnergy("charlotte-a4", 2)
		counter++
		return false
	}, "charlotte-c4")
}

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		ae := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if c.StatusIsActive(c6IcdKey) {
			return false
		}
		if !t.StatusIsActive(skillHoldMarkKey) {
			return false
		}
		if ae.Info.AttackTag != attacks.AttackTagNormal && ae.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		c.AddStatus(c6IcdKey, 360, true)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       c6CoordinateAtk,
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       1.8,
		}

		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.PrimaryTarget(),
			nil,
			c6Radius,
		)

		stats, _ := c.Stats()
		atk := (c.Base.Atk+c.Weapon.BaseAtk)*(1+stats[attributes.ATKP]) + stats[attributes.ATK]

		c.Core.QueueAttack(ai, ap, 0, 0)

		if c.Core.Combat.Player().IsWithinArea(ap) {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: c6HealMsg,
				Src:     atk * 0.42,
				Bonus:   stats[attributes.Heal],
			})
		}
		return false
	}, "charlotte-c6")
}
