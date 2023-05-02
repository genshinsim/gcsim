package kaveh

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c2Key    = "kaveh-c2"
	c6ICDKey = "kaveh-c6-icd"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.Heal] = 0.25
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("kaveh-c1", 180),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = 0.15
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c2Key, burstDuration),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (c *char) c4() {
	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBaseWithHitlag("kaveh-c4", -1),
		Amount: func(ai combat.AttackInfo) (float64, bool) {
			if ai.AttackTag == attacks.AttackTagBloom {
				return 0.6, false
			}
			return 0, false
		},
	})
}

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal &&
			atk.Info.AttackTag != attacks.AttackTagExtra &&
			atk.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		if !c.StatusIsActive(burstKey) {
			return false
		}
		if c.StatusIsActive(c6ICDKey) {
			return false
		}

		c.AddStatus(c6ICDKey, 180, false)

		ai := combat.AttackInfo{
			Abil:       "Pairidaeza's Dreams (C6)",
			ActorIndex: c.Index,
			AttackTag:  attacks.AttackTagNone,
			ICDTag:     attacks.ICDTagNormalAttack,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Dendro,
			Durability: 25,
			Mult:       0.618,
		}
		ap := combat.NewCircleHitOnTarget(t, nil, 4)
		c.Core.QueueAttack(ai, ap, 0, skillHitmark)
		c.Core.Tasks.Add(func() { c.ruptureDendroCores(ap) }, skillHitmark)

		return false
	}, "kaveh-c6")
}
