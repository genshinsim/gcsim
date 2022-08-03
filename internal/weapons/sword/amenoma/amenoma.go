package amenoma

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.AmenomaKageuchi, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	seeds := make([]int, 3) //keep track the seeds
	refund := 4.5 + 1.5*float64(r)
	const icdKey = "amenoma-icd"

	//TODO: this used to be on postskill. make sure nothing broke here
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		// add 1 seed
		if char.StatusIsActive(icdKey) {
			return false
		}
		// find oldest seed to overwrite
		index := 0
		old := seeds[0]

		for i, v := range seeds {
			if v < old {
				old = v
				index = i
			}
		}
		seeds[index] = c.F + 30*60

		c.Log.NewEvent("amenoma proc'd", glog.LogWeaponEvent, char.Index).
			Write("index", index).
			Write("seeds", seeds)

		char.AddStatus(icdKey, 300, true) //5 sec icd

		return false
	}, fmt.Sprintf("amenoma-skill-%v", char.Base.Key.String()))

	//TODO: this used to be on postburst. make sure nothing broke here
	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		count := 0
		for i, v := range seeds {
			if v > c.F {
				count++
			}
			seeds[i] = 0
		}
		if count == 0 {
			return false
		}
		//regen energy after 2 seconds
		char.QueueCharTask(func() {
			char.AddEnergy("amenoma", refund*float64(count))
		}, 120+60) //added 1 extra sec for burst animation but who knows if this is true

		return false
	}, fmt.Sprintf("amenoma-burst-%v", char.Base.Key.String()))
	return w, nil
}
