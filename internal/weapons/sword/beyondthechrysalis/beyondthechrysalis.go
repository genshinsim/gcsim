package beyondthechrysalis

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	cdKey          = "beyond-the-chrysalis-crit-dmg"
	stellarBuffKey = "beyond-the-chrysalis-stellar-buff"
	energyICDKey   = "beyond-the-chrysalis-energy-icd"
)

type Weapon struct {
	Index     int
	buffCount int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	w.buffCount = 0

	critBuff := make([]float64, attributes.EndStatType)
	critBuff[attributes.CD] = 0.4 + 0.16*float64(r)

	stellarBuff := 0.27 + 0.09*float64(r)

	energy := 4.5 + 0.5*float64(r)

	onSkillOrBurst := func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}

		switch w.buffCount {
		case 0:
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag(cdKey, 10*60),
				Amount: func() []float64 {
					return critBuff
				},
			})
		case 1:
			char.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBaseWithHitlag(stellarBuffKey, 10*60),
				Amount: func(ai info.AttackInfo) float64 {
					switch ai.AttackTag {
					case
						attacks.AttackTagDirectStellarSwirl,
						attacks.AttackTagReactionStellarSwirl:
					default:
						return 0
					}
					return stellarBuff
				},
			})
		case 2:
			if !char.StatusIsActive(energyICDKey) {
				char.AddEnergy("beyond-the-chrysalis-energy", energy)
				char.AddStatus(energyICDKey, 4*60, true)
			}
		}
		w.buffCount = (w.buffCount + 1) % 3
	}

	c.Events.Subscribe(event.OnSkill, onSkillOrBurst, "beyond-the-chrysalis-on-skill-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnBurst, onSkillOrBurst, "beyond-the-chrysalis-on-burst-"+char.Base.Key.String())

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)

		if prev != char.Index() {
			return
		}

		w.buffCount = 0
		char.DeleteStatMod(cdKey)
		char.DeleteReactBonusMod(stellarBuffKey)
		char.DeleteStatus(energyICDKey)
	}, "beyond-the-chrysalis-swap-"+char.Base.Key.String())

	return w, nil
}
