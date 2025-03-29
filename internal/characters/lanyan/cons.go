package lanyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2Icd = "lanyan-c2-icd"

var c1Hitmarks = []int{38, 64, 90}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		if !c.hasShield() {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}

		if c.StatusIsActive(c2Icd) {
			return false
		}
		c.AddStatus(c2Icd, 2*60, true)

		c.restoreShield(0.4)
		return false
	}, "lanyan-c2")
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 60
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("lanyan-c4", 12*60),
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}
