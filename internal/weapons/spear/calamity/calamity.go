package calamity

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.CalamityQueller, NewWeapon)
}

type Weapon struct {
	Index        int
	stacks       int
	char         *character.CharWrapper
	c            *core.Core
	icd          int
	lastBuffGain int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func (w *Weapon) incStacks() func() {
	return func() {
		if w.stacks < 6 {
			w.stacks++
			if w.stacks != 6 {
				w.char.QueueCharTask(w.incStacks(), w.icd) // check again in 1s if stacks are not max
			}
		}
		w.c.Log.NewEvent("calamity gained stack", glog.LogWeaponEvent, w.char.Index).
			Write("stacks", w.stacks)
	}
}
func (w *Weapon) checkBuffExpiry(src int) func() {
	return func() {
		if w.lastBuffGain != src {
			w.c.Log.NewEvent("calamity buff expiry check ignored, src diff", glog.LogWeaponEvent, w.char.Index).
				Write("src", src).
				Write("new src", w.lastBuffGain)
			return
		}
		w.stacks = 0
		w.c.Log.NewEvent("calamity buff expired", glog.LogWeaponEvent, w.char.Index).
			Write("src", src).
			Write("lastBuffGain", w.lastBuffGain)
	}
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	// Gain 12% All Elemental DMG Bonus. Obtain Consummation for 20s after using
	// an Elemental Skill, causing ATK to increase by 3.2% per second. This ATK
	// increase has a maximum of 6 stacks. When the character equipped with this
	// weapon is not on the field, Consummation's ATK increase is doubled.
	w := &Weapon{
		char: char,
		c:    c,
	}
	r := p.Refine

	// fixed elemental dmg bonus
	dmg := .09 + float64(r)*.03
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = dmg
	m[attributes.HydroP] = dmg
	m[attributes.CryoP] = dmg
	m[attributes.ElectroP] = dmg
	m[attributes.AnemoP] = dmg
	m[attributes.GeoP] = dmg
	m[attributes.DendroP] = dmg
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("calamity-dmg", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	const buffKey = "calamity-consummation"
	buffDuration := 1200 // 20s * 60
	w.icd = 60           // 1s * 60

	// atk increase per stack after using skill
	// double bonus if not on field
	atkbonus := .024 + float64(r)*.008
	skillPressBonus := make([]float64, attributes.EndStatType)
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}

		// asummes that stacks are not reset on refreshing calamity buff
		w.lastBuffGain = c.F
		char.QueueCharTask(w.checkBuffExpiry(c.F), buffDuration)
		char.QueueCharTask(w.incStacks(), w.icd)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, buffDuration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				atk := atkbonus * float64(w.stacks)
				if c.Player.Active() != char.Index {
					atk *= 2
				}
				skillPressBonus[attributes.ATKP] = atk

				return skillPressBonus, true
			},
		})

		return false
	}, fmt.Sprintf("calamity-queller-%v", char.Base.Key.String()))

	return w, nil
}
