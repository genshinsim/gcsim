package kinich

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c4IcdKey = "kinich-c4-icd-key"
	c6Abil   = "Scalespiker Cannon (C6)"
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	// "After Kinich lands from Canopy Hunter: Riding High's mid-air swing,
	// his Movement SPD will increase by 30% for 6s." is not implemented
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("kinich-c1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold:
			default:
				return nil, false
			}
			if atk.Info.Abil != scalespikerAbil {
				return nil, false
			}

			m[attributes.CD] = 1
			return m, true
		},
	})
}

func (c *char) c2ResShredCB(a combat.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	if a.Damage == 0 {
		return
	}

	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag("kinich-c2", 6*60),
		Ele:   attributes.Dendro,
		Value: -0.3,
	})
}

func (c *char) c2Bonus(ai *combat.AttackInfo) (combat.Snapshot, float64) {
	s := c.Snapshot(ai)
	if c.Base.Cons < 2 {
		return s, 3.0
	}
	if c.c2AoeIncreased {
		return s, 3.0
	}
	c.c2AoeIncreased = true
	s.Stats[attributes.DmgP] += 1.0
	c.Core.Log.NewEvent("Kinich C2 Damage Bonus", glog.LogCharacterEvent, c.Index).
		Write("final", s.Stats[attributes.DmgP])
	return s, 5.0
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	if c.StatusIsActive(c4IcdKey) {
		return
	}
	c.AddStatus(c4IcdKey, 2.8*60, true)
	c.AddEnergy("kinich-c4", 5)

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("kinich-c4-dmgp", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			m[attributes.DmgP] = 0.7
			return m, true
		},
	})
}

func (c *char) c6(ai combat.AttackInfo, s *combat.Snapshot, radius float64, target combat.Target, travel int) {
	if c.Base.Cons < 6 {
		return
	}
	ai.Abil = c6Abil
	var next combat.Target = c.Core.Combat.RandomEnemyWithinArea(combat.NewCircleHitOnTarget(target, nil, radius), func(t combat.Enemy) bool {
		return target.Key() != t.Key()
	})
	if next == nil {
		next = target
	}
	ap := combat.NewCircleHitOnTarget(next, nil, radius)
	c.Core.QueueAttackWithSnap(ai, *s, ap, scalespikerHitmark+travel, c.particleCB, c.a1CB, c.c2ResShredCB)
}
