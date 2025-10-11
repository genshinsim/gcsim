package ineffa

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key = "ineffa-c1"
	c2Key = "ineffa-c2"

	c4Key    = "ineffa-c4"
	c4IcdKey = "ineffa-c4-icd"

	c6Key    = "ineffa-c6"
	c6IcdKey = "ineffa-c6-icd"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	// TODO: Is this buff hitlag affected on Ineffa only? Or is it hitlag affected per character?
	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(c1Key+"-buff", -1),
			Amount: func(ai info.AttackInfo) (float64, bool) {
				if !c.StatusIsActive(c1Key) {
					return 0, false
				}

				switch ai.AttackTag {
				case attacks.AttackTagReactionLunarCharge:
				case attacks.AttackTagDirectLunarCharged:
				default:
					return 0, false
				}

				bonus := min(c.TotalAtk()/100*0.025, 0.50)
				return bonus, false
			},
		})
	}
}

func (c *char) c1OnShield() {
	if c.Base.Cons < 1 {
		return
	}

	c.AddStatus(c1Key, 20*60, true)
}

func (c *char) c2OnBurst() {
	if c.Base.Cons < 2 {
		return
	}

	c.QueueCharTask(c.addShield, 37)
}

func (c *char) c2MakeCB() func(info.AttackCB) {
	if c.Base.Cons < 2 {
		return nil
	}

	done := false
	return func(ac info.AttackCB) {
		if done {
			return
		}

		if ac.Target.Type() != info.TargettableEnemy {
			return
		}

		e, ok := ac.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		done = true

		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Punishment Edit (C2)",
			AttackTag:        attacks.AttackTagDirectLunarCharged,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDirectLunarCharged,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Electro,
			Mult:             3,
			IgnoreDefPercent: 1,
		}

		c2AtkDone := false
		c2Atk := func() {
			if c2AtkDone {
				return
			}
			c2AtkDone = true
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(e, nil, 4), 0, 0)
		}

		// TODO: What is C2's delay?
		c.QueueCharTask(c2Atk, 60)

		c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
			trg, ok := args[0].(*enemy.Enemy)
			if !ok {
				return false
			}

			if trg.Key() != e.Key() {
				return false
			}

			c2Atk()
			// unsubcribe from the event after
			return true
		}, c2Key)
	}
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnLunarCharged, func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		c.AddEnergy(c4Key, 5)
		c.AddStatus(c4IcdKey, 4*60, true)
		return false
	}, c4Key)
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "A Dawning Morn for You (C6)",
		AttackTag:        attacks.AttackTagDirectLunarCharged,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDirectLunarCharged,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Electro,
		Mult:             1.35,
		IgnoreDefPercent: 1,
	}

	c.Core.Events.Subscribe(event.OnLunarChargedReactionAttack, func(args ...any) bool {
		if !c.StatusIsActive(c1Key) {
			return false
		}
		if c.StatusIsActive(c6IcdKey) {
			return false
		}

		c.AddStatus(c6IcdKey, 3.5*60, true)
		// TODO: damage and snapshot delay?
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4), 10, 10)
		return false
	}, c6Key)
}
