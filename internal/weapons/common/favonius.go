package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type Favonius struct {
	Index int
}

func (b *Favonius) SetIndex(idx int) { b.Index = idx }
func (b *Favonius) Init() error      { return nil }

func NewFavonius(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	b := &Favonius{}

	const icdKey = "favonius-cd"

	prob := 0.50 + float64(p.Refine)*0.1
	cd := 810 - p.Refine*90

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		crit := args[3].(bool)
		if dmg == 0 {
			return false
		}
		if !crit {
			return false
		}
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if c.Rand.Float64() > prob {
			return false
		}
		c.Log.NewEvent("favonius proc'd", glog.LogWeaponEvent, char.Index)

		//TODO: used to be 80
		c.QueueParticle("favonius-"+char.Base.Key.String(), 3, attributes.NoElement, char.ParticleDelay)

		//adds a modifier to track icd; this should be fine since it's per char and not global
		char.AddStatus(icdKey, cd, true)

		return false
	}, fmt.Sprintf("favo-%v", char.Base.Key.String()))

	return b, nil
}
