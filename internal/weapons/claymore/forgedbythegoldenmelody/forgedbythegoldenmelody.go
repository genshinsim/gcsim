package forgedbythegoldenmelody

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	atkBuffKey           = "forged-by-the-golden-melody-atk"
	emBuffKey            = "forged-by-the-golden-melody-em"
	stellarBuffKey       = "forged-by-the-golden-melody-stellar"
	atkBuffContraKey     = "forged-by-the-golden-melody-contrapuntal-atk"
	emBuffContraKey      = "forged-by-the-golden-melody-contrapuntal-em"
	stellarBuffContraKey = "forged-by-the-golden-melody-contrapuntal-stellar"
	contraICDKey         = "forged-by-the-golden-melody-contrapuntal-icd"
)

type Weapon struct {
	Index       int
	curBuff     int
	char        *character.CharWrapper
	atkBuff     []float64
	emBuff      []float64
	stellarBuff float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error {
	w.curBuff = 2
	w.switchBuff()
	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		char: char,
	}
	r := p.Refine

	w.atkBuff = make([]float64, attributes.EndStatType)
	w.atkBuff[attributes.ATKP] = 0.135 + float64(r)*0.045

	w.emBuff = make([]float64, attributes.EndStatType)
	w.emBuff[attributes.EM] = 90 + float64(r)*30

	w.stellarBuff = 0.21 + float64(r)*0.07

	onStellar := func(args ...any) {
		if char.StatusIsActive(contraICDKey) {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		switch w.curBuff {
		case 0:
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(atkBuffContraKey, 12*60),
				AffectedStat: attributes.ATKP,
				Amount: func() []float64 {
					return w.atkBuff
				},
			})
		case 1:
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(emBuffContraKey, 12*60),
				AffectedStat: attributes.EM,
				Amount: func() []float64 {
					return w.emBuff
				},
			})
		case 2:
			w.char.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBaseWithHitlag(stellarBuffContraKey, 12*60),
				Amount: func(ai info.AttackInfo) float64 {
					if ai.AttackTag.IsStellar() {
						return w.stellarBuff
					}
					return 0
				},
			})
		}

		char.AddStatus(contraICDKey, 12*60, true)
	}

	c.Events.Subscribe(event.OnStellarConduct, onStellar, fmt.Sprintf("forged-by-the-golden-melody-on-stellar-conduct-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnStellarSwirl, onStellar, fmt.Sprintf("forged-by-the-golden-melody-on-stellar-swirl-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) switchBuff() {
	switch w.curBuff {
	case 0:
		w.char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(emBuffKey, 10*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return w.emBuff
			},
		})
	case 1:
		w.char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(stellarBuffKey, 10*60),
			Amount: func(ai info.AttackInfo) float64 {
				if ai.AttackTag.IsStellar() {
					return w.stellarBuff
				}
				return 0
			},
		})
	default:
		w.char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(atkBuffKey, 10*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return w.atkBuff
			},
		})
	}
	w.curBuff = (w.curBuff + 1) % 3
	w.char.QueueCharTask(w.switchBuff, 10*60)
}
