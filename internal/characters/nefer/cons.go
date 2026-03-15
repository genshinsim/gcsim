package nefer

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2PhantasmBonusKey = "nefer-c2-phantasm"

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c2PhantasmBonusKey, -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if !strings.HasPrefix(atk.Info.Abil, "Phantasm Performance") {
				return nil
			}
			bonus := c.c2PhantasmBonus()
			if bonus <= 0 {
				return nil
			}
			stats := make([]float64, attributes.EndStatType)
			stats[attributes.DmgP] = bonus
			return stats
		},
	})
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	if !c.ascendantGleam {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagDirectLunarBloom {
			return
		}
		atk.Info.Elevation += 0.15
	}, "nefer-c6-elevation")
}
