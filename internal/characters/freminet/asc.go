package freminet

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a4Key = "freminet-a4-buff"
)

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	a4BuffFunc := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		buff := make([]float64, attributes.EndStatType)
		buff[attributes.DmgP] = 0.4

		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(a4Key, 60*5),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if !strings.Contains(atk.Info.Abil, pressureBaseName) {
					return nil, false
				}
				return buff, true
			},
		})

		c.Core.Log.NewEvent("freminet a4 proc", glog.LogCharacterEvent, c.Index)

		return false
	}

	c.Core.Events.Subscribe(event.OnShatter, a4BuffFunc, "freminet-a4")
}
