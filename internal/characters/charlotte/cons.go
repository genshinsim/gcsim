package charlotte

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1StatusKey     = "charlotte-c1"
	c1HealMsg       = "Verification"
	c6HealMsg       = "charlotte-c6-heal"
	c6CoordinateAtk = "charlotte-c6-coordinate-atk"
	c6IcdKey        = "charlotte-c6-icd"
	c6AttackRadius  = 3.5
	c6HealRadius    = 5
)

func (c *char) c1Heal(char *character.CharWrapper) func() {
	return func() {
		if !char.StatusIsActive(c1StatusKey) {
			return
		}

		stats, _ := c.Stats()
		atk := (c.Base.Atk+c.Weapon.BaseAtk)*(1+stats[attributes.ATKP]) + stats[attributes.ATK]
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  char.Index,
			Message: c1HealMsg,
			Src:     atk * 0.8,
			Bonus:   stats[attributes.Heal],
		})
		char.QueueCharTask(c.c1Heal(char), 2*60)
	}
}

func (c *char) c1() {
	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		src := args[0].(*player.HealInfo)
		if src.Message != healInitialMsg && src.Message != healDotMsg && src.Message != c6HealMsg {
			return false
		}

		target := args[1].(int)
		char := c.Core.Player.ByIndex(target)
		if char.StatusIsActive(c1StatusKey) {
			return false
		}
		char.AddStatus(c1StatusKey, 6*60, true)
		c.c1Heal(char)()

		return false
	}, "charlotte-c1")
}

func (c *char) makeC2CB() combat.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}
	c.c2Hits = 0
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.c2Hits == 3 {
			return
		}
		c.c2Hits++
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.1 * float64(c.c2Hits)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("charlotte-c2", 12*60),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
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
		dmg := 0.1
		ae.Snapshot.Stats[attributes.DmgP] += dmg
		c.Core.Log.NewEvent("charlotte c4 adding dmg%", glog.LogCharacterEvent, c.Index).Write("dmg%", dmg)
		c.AddEnergy("charlotte-c4", 2)
		counter++
		return false
	}, "charlotte-c4")
}

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if c.StatusIsActive(c6IcdKey) {
			return false
		}
		if !t.StatusIsActive(skillHoldMarkKey) {
			return false
		}
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if ae.Info.AttackTag != attacks.AttackTagNormal && ae.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		c.AddStatus(c6IcdKey, 6*60, true)
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               c6CoordinateAtk,
			AttackTag:          attacks.AttackTagElementalBurst,
			ICDTag:             attacks.ICDTagNone,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Cryo,
			Durability:         25,
			Mult:               1.8,
			HitlagFactor:       0.02,
			CanBeDefenseHalted: true,
			IsDeployable:       true,
		}
		pos := t.Pos()
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(pos, nil, c6AttackRadius), 0, 0)
			if c.Core.Combat.Player().IsWithinArea(combat.NewCircleHitOnTarget(pos, nil, c6HealRadius)) {
				stats, _ := c.Stats()
				atk := (c.Base.Atk+c.Weapon.BaseAtk)*(1+stats[attributes.ATKP]) + stats[attributes.ATK]
				c.Core.Player.Heal(player.HealInfo{
					Caller:  c.Index,
					Target:  c.Core.Player.Active(),
					Message: c6HealMsg,
					Src:     atk * 0.42,
					Bonus:   stats[attributes.Heal],
				})
			}
		}, 14)
		return false
	}, "charlotte-c6")
}
