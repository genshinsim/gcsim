package fangofthemountainking

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	canopyFavorKey = "canopy-favor"
	skillIcdKey    = "fotmk-skill-icd"
	reactIcdKey    = "fotmk-react-icd"

	stackDuration = 6 * 60
)

func init() {
	core.RegisterWeaponFunc(keys.FangOfTheMountainKing, NewWeapon)
}

type Weapon struct {
	Index        int
	char         *character.CharWrapper
	stackTracker *stackTracker
	buffStack    float64
	mBuff        []float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Gain 1 stack of Canopy's Favor after hitting an opponent with an Elemental Skill. This can be triggered once every 0.5s.
// After a nearby party member triggers a Burning or Burgeon reaction, the equipping character will gain 3 stacks.
// This effect can be triggered once every 2s and can be triggered even when the triggering party member is off-field.
// Canopy's Favor: Elemental Skill and Burst DMG is increased by 10% for 6s. Max 6 stacks. Each stack is counted independently.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		char:         char,
		stackTracker: newStackTracker(6, char.QueueCharTask, &c.F),
		buffStack:    0.10 + float64(p.Refine-1)*0.025,
		mBuff:        make([]float64, attributes.EndStatType),
	}

	//nolint:unparam // ignoring for now, event refactor should get rid of bool return of event sub
	onReact := func(...interface{}) bool {
		if char.StatusIsActive(reactIcdKey) {
			return false
		}
		char.AddStatus(reactIcdKey, 2*60, true)

		w.addStacks(3)
		return false
	}
	c.Events.Subscribe(event.OnBurning, func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		return onReact()
	}, fmt.Sprintf("fangofthemountainking-burning-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnBurgeon, onReact, fmt.Sprintf("fangofthemountainking-burgeon-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != w.char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}

		if char.StatusIsActive(skillIcdKey) {
			return false
		}
		char.AddStatus(skillIcdKey, .5*60, false)

		w.addStacks(1)
		return false
	}, fmt.Sprintf("fangofthemountainking-ondmg-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) addStacks(num int) {
	for i := 0; i < num; i++ {
		w.stackTracker.Add(stackDuration)
	}

	w.char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(canopyFavorKey, stackDuration),
		Amount: func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch a.Info.AttackTag {
			case attacks.AttackTagElementalArt:
			case attacks.AttackTagElementalArtHold:
			case attacks.AttackTagElementalBurst:
			default:
				return nil, false
			}
			w.mBuff[attributes.DmgP] = w.buffStack * float64(w.stackTracker.Count())
			return w.mBuff, true
		},
	})
}
