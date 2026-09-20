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
	weaponKey     = "golden-melody"
	atkBuffSuffix = "-atk"
	emBuffSuffix  = "-em"
	stellarSuffix = "-stellar"

	contraICDKey = "golden-melody-contrapuntal-icd"
)

type Weapon struct {
	Index       int
	curBuff     BuffType
	char        *character.CharWrapper
	atkBuff     []float64
	emBuff      []float64
	stellarBuff float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		char: char,
	}
	r := p.Refine

	w.curBuff = BuffType(p.Params["buff_type"] % 3)

	w.atkBuff = make([]float64, attributes.EndStatType)
	w.atkBuff[attributes.ATKP] = 0.135 + float64(r)*0.045

	w.emBuff = make([]float64, attributes.EndStatType)
	w.emBuff[attributes.EM] = 90 + float64(r)*30

	w.stellarBuff = 0.21 + float64(r)*0.07

	w.applyBuff(false, 10*60)
	w.char.QueueCharTask(w.switchBuff, 10*60)

	onStellar := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		if char.StatusIsActive(contraICDKey) {
			return
		}

		char.AddStatus(contraICDKey, 12*60, true)

		w.applyBuff(true, 12*60)
	}

	c.Events.Subscribe(event.OnStellarConduct, onStellar, "forged-by-the-golden-melody-on-stellar-conduct-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnStellarSwirl, onStellar, "forged-by-the-golden-melody-on-stellar-swirl-"+char.Base.Key.String())

	return w, nil
}

func (w *Weapon) switchBuff() {
	w.curBuff = w.curBuff.Next()
	w.applyBuff(false, 10*60)
	w.char.QueueCharTask(w.switchBuff, 10*60)
}

func (w *Weapon) applyBuff(contrapuntal bool, dur int) {
	middle := ""
	if contrapuntal {
		middle = "-contrapuntal"
	}

	switch w.curBuff {
	case EmBuff:
		w.char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(weaponKey+middle+emBuffSuffix, dur),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return w.emBuff
			},
		})
	case StellarBuff:
		w.char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(weaponKey+middle+stellarSuffix, dur),
			Amount: func(ai info.AttackInfo) float64 {
				if ai.AttackTag.IsStellar() {
					return w.stellarBuff
				}
				return 0
			},
		})
	case AtkBuff:
		w.char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(weaponKey+middle+atkBuffSuffix, dur),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return w.atkBuff
			},
		})
	default:
		panic(fmt.Sprintf("golden melody: unknown buff type: %v", w.curBuff))
	}
}

type BuffType int

const (
	AtkBuff BuffType = iota
	EmBuff
	StellarBuff
)

func (t BuffType) Next() BuffType {
	return (t + 1) % 3
}
