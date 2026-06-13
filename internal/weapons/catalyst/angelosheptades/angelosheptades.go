package angelosheptades

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	angelosHeptadesEnergyICDKey  = "angelos-heptades-energy-icd"
	angelosHeptadesEnergyKey     = "angelos-heptades-energy"
	angelosHeptadesHolderBuffKey = "angelos-heptades-buff"
)

func init() {
	core.RegisterWeaponFunc(keys.AngelosHeptades, NewWeapon)
}

type Weapon struct {
	Index   int
	Core    *core.Core
	Char    *character.CharWrapper
	BuffSrc int
	BuffAmt float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error {
	n := make([]float64, attributes.EndStatType)
	for _, c := range w.Core.Player.Chars() {
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("pathfinders-light", -1),
			Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
				if !w.Char.StatusIsActive(angelosHeptadesHolderBuffKey) {
					return nil
				}

				n[attributes.DmgP] = w.BuffAmt
				if c.Index() != w.Core.Player.Active() {
					if c.IsHexerei {
						n[attributes.DmgP] *= 0.5
						return n
					}
					return nil
				}

				return n
			},
		})
	}

	return nil
}

func NewWeapon(core *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		Core:    core,
		Char:    char,
		BuffSrc: -1,
		BuffAmt: 0,
	}
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

		w.BuffSrc = core.F
		char.AddStatus(angelosHeptadesHolderBuffKey, 20*60, true)

		w.updateBuff(char, core, r, w.BuffSrc)

		if char.StatusIsActive(angelosHeptadesEnergyICDKey) {
			return
		}

		energyAmount := 13.0 + float64(r)

		char.AddStatus(angelosHeptadesEnergyICDKey, 14*60, true)
		char.AddEnergy(angelosHeptadesEnergyKey, energyAmount)
	}, fmt.Sprintf("angelos-heptades-shielded-%v", char.Base.Key.String()))
	return w, nil
}

func (w *Weapon) updateBuff(char *character.CharWrapper, core *core.Core, r, src int) {
	if w.BuffSrc != src {
		return
	}

	buff := 0.07 + float64(r)*0.03
	buffCap := 0.18 + float64(r)*0.08
	w.BuffAmt = min(char.NonExtraStat(attributes.ATK)/1000.0*buff, buffCap)

	if core.F-w.BuffSrc >= 20*60 {
		w.BuffAmt = 0
		return
	}

	core.Log.NewEvent("angelos heptades buff updated", glog.LogWeaponEvent, char.Index()).
		Write("src", w.BuffSrc).
		Write("value", w.BuffAmt).
		Write("holder-atk", char.NonExtraStat(attributes.ATK))

	core.Tasks.Add(func() {
		w.updateBuff(char, core, r, src)
	}, 0.5*60)
}
