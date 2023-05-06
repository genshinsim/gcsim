package kaveh

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1ICDKey = "kaveh-a1-icd"
	a4Key    = "kaveh-a4"
	a4ICDKey = "kaveh-a4-icd"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnPlayerHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagBloom &&
			atk.Info.AttackTag != attacks.AttackTagHyperbloom &&
			atk.Info.AttackTag != attacks.AttackTagBurgeon {
			return false
		}
		if c.StatusIsActive(a1ICDKey) {
			return false
		}
		c.AddStatus(a1ICDKey, 30, false)
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "Creator's Undertaking (A1)",
			Src:     3.0 * c.Stat(attributes.EM),
			Bonus:   c.Stat(attributes.Heal),
		})
		return false
	}, "kaveh-a1")
}

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(a4Key, burstDuration),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			m[attributes.EM] = float64(25 * c.a4Stacks)
			return m, true
		},
	})
}

func (c *char) a4AddStacksHandler() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		if c.a4Stacks >= 4 {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal &&
			atk.Info.AttackTag != attacks.AttackTagExtra &&
			atk.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}
		if !c.StatusIsActive(burstKey) {
			return false
		}
		if c.StatusIsActive(a4ICDKey) {
			return false
		}

		c.AddStatus(a4ICDKey, 6, true)
		c.a4Stacks++
		return false
	}, "kaveh-a4")
}
