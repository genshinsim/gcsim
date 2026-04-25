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

const energyIcdKey = "bloodsoakedruins-energy-icd"

func init() {
	core.RegisterWeaponFunc(keys.BloodsoakedRuins, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

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
	lcBonus := 0.24 + float64(r)*0.12
	mCrit := make([]float64, attributes.EndStatType)
	mCrit[attributes.CD] = 0.21 + float64(r)*0.07

	c.Events.Subscribe(event.OnBurst, func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag("bloodsoakedruins-lc", 3.5*60),
			Amount: func(ai info.AttackInfo) float64 {
				if ai.AttackTag != attacks.AttackTagReactionLunarCharge && ai.AttackTag != attacks.AttackTagDirectLunarCharged {
					return 0
				}
				return lcBonus
			},
		})
	}, fmt.Sprintf("bloodsoakedruins-burst-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnLunarCharged, func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}

		ae := args[1].(*info.AttackEvent)
		if ae.Info.ActorIndex != char.Index() {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("bloodsoakedruins-cd", 6*60),
			AffectedStat: attributes.CD,
			Amount: func() []float64 {
				return mCrit
			},
		})

		if char.StatusIsActive(energyIcdKey) {
			return
		}
		char.AddStatus(energyIcdKey, 14*60, true)
		char.AddEnergy("bloodsoakedruins", energyRestore)
	}, fmt.Sprintf("bloodsoakedruins-lunarcharged-%v", char.Base.Key.String()))

	return w, nil
}
