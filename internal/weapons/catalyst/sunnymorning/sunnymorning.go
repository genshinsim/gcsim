package sunnymorning

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

func init() {
	core.RegisterWeaponFunc(keys.SunnyMorningSleepIn, NewWeapon)
}

type Weapon struct {
	Index       int
	c           *core.Core
	self        *character.CharWrapper
	emBuffSwirl []float64
	emBuffSkill []float64
	emBuffBurst []float64
}

func (w *Weapon) SetIndex(idx int) {
	w.Index = idx
}

func (w *Weapon) Init() error {
	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		c:    c,
		self: char,
	}
	r := p.Refine

	multiplier := 0.75 + (0.25 * float64(r))

	w.emBuffSwirl = make([]float64, attributes.EndStatType)
	w.emBuffSwirl[attributes.EM] = 120 * multiplier

	w.emBuffSkill = make([]float64, attributes.EndStatType)
	w.emBuffSkill[attributes.EM] = 96 * multiplier

	w.emBuffBurst = make([]float64, attributes.EndStatType)
	w.emBuffBurst[attributes.EM] = 32 * multiplier

	frameSwirlBuffApplied := -1
	swirlFunc := func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		// Although from the description it can be implied that anyone's swirl can trigger it,
		// only the wielder swirls trigger.
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if frameSwirlBuffApplied == c.F {
			// avoid doing this 2 times on double swirls
			return false
		}
		frameSwirlBuffApplied = c.F

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("sunny-morning-swirl", 6*60),
			Amount: func() ([]float64, bool) {
				return w.emBuffSwirl, true
			},
		})

		return false
	}

	skillFunc := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("sunny-morning-skill", 9*60),
			Amount: func() ([]float64, bool) {
				return w.emBuffSkill, true
			},
		})
		return false
	}

	burstFunc := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
			return false
		}

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBase("sunny-morning-burst", 30*60),
			Amount: func() ([]float64, bool) {
				return w.emBuffBurst, true
			},
		})

		return false
	}

	c.Events.Subscribe(event.OnSwirlElectro, swirlFunc, fmt.Sprintf("sunny-morning-%v-electro-swirl", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSwirlCryo, swirlFunc, fmt.Sprintf("sunny-morning-%v-cryo-swirl", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSwirlHydro, swirlFunc, fmt.Sprintf("sunny-morning-%v-hydro-swirl", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSwirlPyro, swirlFunc, fmt.Sprintf("sunny-morning-%v-pyro-swirl", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyHit, skillFunc, fmt.Sprintf("sunny-morning-%v-skill-hit", char.Base.Key.String()))
	c.Events.Subscribe(event.OnEnemyHit, burstFunc, fmt.Sprintf("sunny-morning-%v-burst-hit", char.Base.Key.String()))

	return w, nil
}
