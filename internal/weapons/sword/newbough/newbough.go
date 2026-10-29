package newbough

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type Weapon struct {
	stacks int
	Index  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	skillKey   = "newbough-skill"
	buffKey    = "newbough-verdant"
	icdKey     = "newbough-icd"
	stellarKey = "newbough-verdant-stellar"
)

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	normalATK := 0.03 + 0.01*float64(p.Refine)
	normalEM := 15 + 5*float64(p.Refine)
	radianceATK := 0.045 + 0.015*float64(p.Refine)
	radianceDmg := 0.06 + 0.02*float64(p.Refine)
	val := make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		if c.Player.Active() == char.Index() {
			char.AddStatus(skillKey, 12*60, true)
		}
	}, fmt.Sprintf("newbough-skill-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if !char.StatusIsActive(skillKey) {
			return
		}
		if char.StatusIsActive(icdKey) {
			return
		}

		if !char.StatModIsActive(buffKey) {
			w.stacks = 0
		}
		if w.stacks < 3 {
			w.stacks++
		}

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(buffKey, 6*60),
			Amount: func() []float64 {
				clear(val)
				if isRadiance(char) {
					val[attributes.ATKP] = radianceATK * float64(w.stacks)
				} else {
					val[attributes.ATKP] = normalATK * float64(w.stacks)
					val[attributes.EM] = normalEM * float64(w.stacks)
				}
				return val
			},
		})
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(stellarKey, 6*60),
			Amount: func(ai info.AttackInfo) float64 {
				if isRadiance(char) && ai.AttackTag.IsStellar() {
					return radianceDmg * float64(w.stacks)
				}
				return 0
			},
		})

		char.AddStatus(icdKey, 60, true)
	}, fmt.Sprintf("newbough-verdant-%v", char.Base.Key.String()))

	return w, nil
}

func isRadiance(char *character.CharWrapper) bool {
	return char.StatusIsActive(reactable.PolestarFieldKey) || char.StatusIsActive(reactable.RadianceSwirlKey)
}
