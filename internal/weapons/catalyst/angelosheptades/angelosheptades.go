package angelosheptades

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	angelosHeptadesEnergyICDKey = "angelos-heptades-energy-icd"
	angelosHeptadesEnergyKey    = "angelos-heptades-energy"
)

func init() {
	core.RegisterWeaponFunc(keys.AngelosHeptades, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(core *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATK] = 0.09 + float64(r)*0.03

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("angelos-heptades-atk", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() []float64 {
			return m
		},
	})

	core.Events.Subscribe(event.OnShielded, func(args ...any) {
		shd := args[0].(shield.Shield)
		if shd.ShieldOwner() != char.Index() {
			return
		}
		// TODO: Not sure if the character needs to be on the field
		if core.Player.Active() != char.Index() {
			return
		}

		buff := 0.07 + float64(r)*0.03
		buffCap := 0.18 + float64(r)*0.08

		for _, c := range core.Player.Chars() {
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBase("pathfinders-light", 20*60),
				AffectedStat: attributes.DmgP,
				Amount: func() []float64 {
					n := make([]float64, attributes.EndStatType)
					n[attributes.DmgP] = min(char.TotalAtk()/1000.0*buff, buffCap)
					if c.Index() != core.Player.Active() {
						return nil
					}

					return n
				},
			})
		}

		if char.StatusIsActive(angelosHeptadesEnergyICDKey) {
			return
		}

		energyAmount := 13.0 + float64(r)

		char.AddStatus(angelosHeptadesEnergyICDKey, 14*60, false)
		char.AddEnergy(angelosHeptadesEnergyKey, energyAmount)
	}, fmt.Sprintf("angelos-heptades-shielded-%v", char.Base.Key.String()))
	return w, nil
}
