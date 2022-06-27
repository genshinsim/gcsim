package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type Favonius struct {
	Index int
}

func (b *Favonius) SetIndex(idx int) { b.Index = idx }
func (b *Favonius) Init() error      { return nil }

func NewFavonius(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	b := &Favonius{}

	prob := 0.50 + float64(p.Refine)*0.1
	cd := 810 - p.Refine*90
	icd := 0

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if !crit {
			return false
		}
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if icd > c.F {
			return false
		}

		if c.Rand.Float64() > prob {
			return false
		}
		c.Log.NewEvent("favonius proc'd", glog.LogWeaponEvent, char.Index)

		c.QueueParticle("favonius-"+char.Base.Key.String(), 3, attributes.NoElement, 80)

		icd = c.F + cd

		return false
	}, fmt.Sprintf("favo-%v", char.Base.Key.String()))

	return b, nil
}
