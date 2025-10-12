package flins

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
	c1Key    = "flins-c1"
	c1IcdKey = "flins-c1-icd"
	c2Key    = "flins-c2"
	c4Key    = "flins-c4"
	c6Key    = "flins-c6"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnLunarCharged, func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		c.AddEnergy(c1Key, 8)
		c.AddStatus(c1IcdKey, 5.5*60, true)
		return false
	}, c1Key)
}

func (c *char) c1SkillCD() int {
	if c.Base.Cons < 1 {
		return 6 * 60
	}
	return 4 * 60
}

func (c *char) c2OnSkill() {
	if c.Base.Cons < 2 {
		return
	}
	c.AddStatus(c2Key, 6*60, true)
}

func (c *char) c2MakeAtkCB() func(info.AttackCB) {
	if c.Base.Cons < 2 {
		return nil
	}
	if !c.StatusIsActive(c2Key) {
		return nil
	}
	c.DeleteStatus(c2Key)
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
			Abil:             "The Devil's Wall (C2)",
			AttackTag:        attacks.AttackTagDirectLunarCharged,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDirectLunarCharged,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Electro,
			Mult:             0.5,
			IgnoreDefPercent: 1,
		}
		// TODO: What is C2's delay?
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(e, nil, 4), 20, 20)
	}
}

func (c *char) c2GleamInit() {
	if c.Base.Cons < 2 {
		return
	}

	if c.getMoonsignLevel() < 2 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		if c.Core.Player.Active() != c.Index() {
			return false
		}

		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		atk := args[1].(*info.AttackEvent)

		if atk.Info.ActorIndex != c.Index() {
			return false
		}
		if atk.Info.Element != attributes.Electro {
			return false
		}

		t.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag(c2Key, 7*60),
			Ele:   attributes.Electro,
			Value: -0.25,
		})
		return false
	}, c2Key+"-gleam")
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.2

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c4Key, -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (c *char) c4A4() (float64, float64) {
	if c.Base.Cons < 4 {
		return 0.08, 160.0
	}
	return 0.1, 220.0
}

func (c *char) c6Init() {
	flinsMult := 1.35
	otherMult := 1.0

	if c.getMoonsignLevel() >= 2 {
		flinsMult += 0.1
		otherMult += 0.1
	}

	// TODO: How to do elevate?
	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) bool {
		atk := args[0].(*info.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarCharged:
		case attacks.AttackTagReactionLunarCharge:
		default:
			return false
		}

		if atk.Info.ActorIndex == c.Index() {
			atk.Info.Mult *= flinsMult
			atk.Info.FlatDmg *= flinsMult
			return false
		}

		if c.getMoonsignLevel() < 2 {
			return false
		}

		atk.Info.Mult *= otherMult
		atk.Info.FlatDmg *= otherMult

		return false
	}, c6Key)
}
