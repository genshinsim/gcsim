package skyrider

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.SkyriderGreatsword, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
	buff   []float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//On hit, Normal or Charged Attacks increase ATK by 6% for 6s. Max 4 stacks. Can occur once every 0.5s.
	w := &Weapon{}
	r := p.Refine

	atkbuff := 0.05 + float64(r)*0.01
	w.buff = make([]float64, attributes.EndStatType)
	const icdKey = "skyrider-icd"

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		//if hit lands after all stack should have fallen off, reset to 0
		if !char.StatModIsActive("skyrider") {
			w.stacks = 0
		}

		if w.stacks < 4 {
			w.stacks++
			//update buff
			w.buff[attributes.ATKP] = float64(w.stacks) * atkbuff
		}

		//extend buff timer
		char.AddStatus(icdKey, 30, true)

		//every whack adds a stack while under 4 and refreshes buff
		//lasts 6 seconds
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("skyrider", 360),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return w.buff, true
			},
		})

		return false
	}, fmt.Sprintf("skyrider-greatsword-%v", char.Base.Key.String()))

	return w, nil
}
