package amenoma

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterWeaponFunc(keys.AmenomaKageuchi, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	seeds := []string{"amenoma-seed-0", "amenoma-seed-1", "amenoma-seed-2"}
	refund := 4.5 + 1.5*float64(r)
	const icdKey = "amenoma-icd"

	// TODO: this used to be on postskill. make sure nothing broke here
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
		for i, s := range seeds {
			if char.StatusExpiry(s) < char.StatusExpiry(seeds[index]) {
				index = i
			}
		}
		char.AddStatus(seeds[index], 30*60, true)

		c.Log.NewEvent("amenoma proc'd", glog.LogWeaponEvent, char.Index).
			Write("index", index)

		char.AddStatus(icdKey, 300, true) // 5 sec icd

		return false
	}, fmt.Sprintf("amenoma-skill-%v", char.Base.Key.String()))

	// TODO: this used to be on postburst. make sure nothing broke here
	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		count := 0
		for _, s := range seeds {
			if char.StatusIsActive(s) {
				count++
			}
			char.DeleteStatus(s)
		}
		if count == 0 {
			return false
		}
		// regen energy after 2 seconds
		char.QueueCharTask(func() {
			char.AddEnergy("amenoma", refund*float64(count))
		}, 120+60) // added 1 extra sec for burst animation but who knows if this is true

		return false
	}, fmt.Sprintf("amenoma-burst-%v", char.Base.Key.String()))
	return w, nil
}
