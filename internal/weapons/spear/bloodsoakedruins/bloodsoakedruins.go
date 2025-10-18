package bloodsoakedruins

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.BloodsoakedRuins, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	LunarBuffKey      = "bloodstained-ruins"
	LunarBuffDuration = int(3.5 * 60)
	RequiemKey        = "reqium"
	RequiemDuration   = 6 * 60
	energySrc         = "ruins"
	energyIcdKey      = "ruins-energy-icd"
	energyIcd         = 14 * 60
)

// For 3.5s after using an Elemental Burst,
// the equipping character's Lunar-Charged DMG dealt to opponents is increased by 36%/48%/60%/72%/84%.
// Additionally, after triggering a Lunar-Charged reaction,
// the equipping character will gain Requiem of Ruin: CRIT DMG is increased by 28%/35%/42%/49%/56% for 6s.
// They will also regain 12/13/14/15/16 Elemental Energy.
// Elemental Energy can be restored this way once every 14s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine
	energyRestore := 11 + float64(r)

	// Burst-triggered Bonus ===
	c.Events.Subscribe(event.OnBurst, func(args ...any) bool {
		if c.Player.Active() != char.Index() {
			return false
		}

		// Apply a DMG Bonus modifier for Lunar-Charged attacks
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(LunarBuffKey, LunarBuffDuration),
			Amount: func(atk *info.AttackEvent, t info.Target) ([]float64, bool) {
				// Checking if the attack is Lunar-Charged
				if atk.Info.AttackTag != attacks.AttackTagLunarCharged {
					return nil, false
				}

				mDmg := make([]float64, attributes.EndStatType)
				mDmg[attributes.DmgP] = 0.24 + float64(r)*0.12
				return mDmg, true
			},
		})
		return false
	}, fmt.Sprintf("bloodsoakedruins-burst-%v", char.Base.Key.String()))

	// Lunar-Charged Reaction trigger Bonuses (Requiem of Ruin + Energy) ===
	c.Events.Subscribe(event.OnLunarCharged, func(args ...any) bool {
		if c.Player.Active() != char.Index() {
			return false
		}
		// Requiem of Ruin: +28% CRIT DMG for 6s
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(RequiemKey, RequiemDuration),
			AffectedStat: attributes.CD,
			Amount: func() ([]float64, bool) {
				m := make([]float64, attributes.EndStatType)
				m[attributes.CD] = 0.21 + float64(r)*0.07
				return m, true
			},
		})

		// Restore 12 Energy, but only once every 14s
		if char.StatusIsActive(energyIcdKey) {
			return false
		}
		char.AddStatus(energyIcdKey, energyIcd, true)
		char.AddEnergy(energySrc, energyRestore)

		return false
	}, fmt.Sprintf("bloodsoakedruins-lunarcharged-%v", char.Base.Key.String()))

	return w, nil
}
