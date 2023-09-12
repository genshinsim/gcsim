package haran

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.HaranGeppakuFutsu, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	icdKey             = "haran-icd"
	maxWavespikeStacks = 2
)

// Obtain 12% All Elemental DMG Bonus. When other nearby party members use
// Elemental Skills, the character equipping this weapon will gain 1 Wavespike
// stack. Max 2 stacks. This effect can be triggered once every 0.3s. When the
// character equipping this weapon uses an Elemental Skill, all stacks of
// Wavespike will be consumed to gain Rippling Upheaval: each stack of Wavespike
// consumed will increase Normal Attack DMG by 20% for 8s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// perm buff
	m := make([]float64, attributes.EndStatType)
	base := 0.09 + float64(r)*0.03
	m[attributes.PyroP] = base
	m[attributes.HydroP] = base
	m[attributes.CryoP] = base
	m[attributes.ElectroP] = base
	m[attributes.AnemoP] = base
	m[attributes.GeoP] = base
	m[attributes.DendroP] = base
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("haran-ele-bonus", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	wavespikeStacks := 0

	nonActiveFn := func() bool {
		// once every 0.3s
		if char.StatusIsActive(icdKey) {
			return false
		}
		// add stacks
		wavespikeStacks++
		if wavespikeStacks > maxWavespikeStacks {
			wavespikeStacks = maxWavespikeStacks
		}
		c.Log.NewEvent("Haran gained a wavespike stack", glog.LogWeaponEvent, char.Index).Write("stack", wavespikeStacks)
		char.AddStatus(icdKey, 18, true)
		return false
	}

	val := make([]float64, attributes.EndStatType)
	activeFn := func() bool {
		val[attributes.DmgP] = (0.15 + float64(r)*0.05) * float64(wavespikeStacks)
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("ripping-upheaval", 480),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagNormal {
					return nil, false
				}
				return val, true
			},
		})
		wavespikeStacks = 0
		return false
	}

	//TODO: this used to be on post. make sure nothing broke here
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() == char.Index {
			return activeFn()
		}
		return nonActiveFn()
	}, fmt.Sprintf("wavespike-%v", char.Base.Key.String()))

	return w, nil
}
