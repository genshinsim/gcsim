package collei

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("collei-c1", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.Active() != c.Index {
				return m, true
			}
			return nil, false
		},
	})
}

func (c *char) c2() {
	for _, event := range dendroEvents {
		c.Core.Events.Subscribe(event, func(args ...interface{}) bool {
			if c.sproutShouldExtend {
				return false
			}
			if !(c.StatusIsActive(sproutKey) || c.StatusIsActive(skillKey)) {
				return false
			}
			c.sproutShouldExtend = true
			if c.StatusIsActive(sproutKey) {
				c.ExtendStatus(sproutKey, 180)
			}
			c.Core.Log.NewEvent("collei c2 proc", glog.LogCharacterEvent, c.Index)
			return false
		}, "collei-c2")
	}
}

func (c *char) c4() {
	for i, char := range c.Core.Player.Chars() {
		// does not affect collei
		if c.Index == i {
			continue
		}
		amts := make([]float64, attributes.EndStatType)
		amts[attributes.EM] = 60
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("collei-c4", 720),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return amts, true
			},
		})
	}
}

func (c *char) c6(t combat.Target) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Forest of Falling Arrows (C6)",
		AttackTag:  attacks.AttackTagNone, // in game has this as AttackTagColleiC6
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       2,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(t, nil, 4),
		0,
		22,
	)
}
