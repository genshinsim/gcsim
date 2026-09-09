package exaiphanesblade

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const onHitICDKey = "exaiphanes-blade-on-hit-icd"

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	switch char.Base.Key {
	case keys.AetherAnemo:
	case keys.AetherGeo:
	case keys.AetherElectro:
	case keys.AetherDendro:
	case keys.AetherHydro:
	case keys.AetherPyro:
	case keys.LumineAnemo:
	case keys.LumineGeo:
	case keys.LumineElectro:
	case keys.LumineDendro:
	case keys.LumineHydro:
	case keys.LuminePyro:
	default:
		return w, nil
	}

	if r >= 2 {
		ele, ok := p.Params["elements"]
		if !ok {
			ele = 7
		}

		critBuff := make([]float64, attributes.EndStatType)
		critBuff[attributes.CD] = float64(ele) * 0.06

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("exaiphanes-blade-cd", -1),
			AffectedStat: attributes.CD,
			Amount: func() []float64 {
				return critBuff
			},
		})
	}

	atkBuff := make([]float64, attributes.EndStatType)
	atkBuff[attributes.ATKP] = 0.12 + float64(r)*0.04

	energy := 3.0

	if r >= 3 {
		energy = 5.0
		atkBuff[attributes.ATKP] = float64(r) * 0.08
	}

	onHit := func(args ...any) {
		if char.StatusIsActive(onHitICDKey) {
			return
		}
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		char.AddStatus(onHitICDKey, 5*60, true)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("exaiphanes-blade-atk", 8*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return atkBuff
			},
		})
		char.AddEnergy("exaiphanes-blade-energy", energy)
	}

	c.Events.Subscribe(event.OnEnemyHit, onHit, "exaiphanes-blade-on-hit")

	return w, nil
}
