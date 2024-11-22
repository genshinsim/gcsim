package astralvulturescrimsonplumage

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.AstralVulturesCrimsonPlumage, NewWeapon)
}

type Weapon struct {
	Index int
	r     float64
	core  *core.Core
	char  *character.CharWrapper
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error {
	counter := 0
	for _, x := range w.core.Player.Chars() {
		if x.Base.Element != w.char.Base.Element {
			counter++
		}
	}
	if counter == 0 {
		return nil
	}

	m := make([]float64, attributes.EndStatType)
	dmg := 0.025*w.r + 0.075
	if counter >= 2 {
		dmg *= 2.4
	}

	w.char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("astralvulturescrimsonplumage-dmg", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagExtra:
				m[attributes.DmgP] = dmg * 2
			case attacks.AttackTagElementalBurst:
				m[attributes.DmgP] = dmg
			default:
				return nil, false
			}
			return m, true
		},
	})

	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		core: c,
		char: char,
		r:    float64(p.Refine),
	}

	atkp := make([]float64, attributes.EndStatType)
	atkp[attributes.ATKP] = 0.06*float64(p.Refine) + 0.18

	for i := event.OnSwirlHydro; i <= event.OnSwirlPyro; i++ {
		c.Events.Subscribe(i, func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			if c.Player.Active() != char.Index {
				return false
			}
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("astralvulturescrimsonplumage-atkp", 12*60),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					return atkp, true
				},
			})
			return false
		}, "astralvulturescrimsonplumage-swirl-"+char.Base.Key.String())
	}

	return w, nil
}
