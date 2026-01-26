package flameforgedinsight

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.FlameForgedInsight, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Upon causing an Electro-Charged, Lunar-Bloom, Bloom, Crystallize reaction, EM is increased by 60/75/90/105/120% for 15s & restore 12/15/18/21/24 Elemental Energy. This effect can be triggerred at most once every 15s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine
	const icdKey = "flameforgedinsight-icd"
	const buffKey = "mindinbloom"

	//Refinement Scaling
	// EM buff amount is 45 + (Refine Lvl * 15)
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 45 + float64(r)*15 // 60 -> 120
	// Flat energy restore amount is 9 + (Refine Lvl * 3)
	regen := 9 + float64(r)*3 // 12 -> 24

	//nolint:unparam // ignoring for now, event refactor should get rid of bool return of event sub
	// proc logic
	buff := func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		// Only proc for reactions caused by the holder
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		// Only proc for when internal cooldown is not active
		if char.StatusIsActive(icdKey) {
			return false
		}

		// Add internal cooldown on buff for 15s (900 frames at 60hz) if not present
		char.AddStatus(icdKey, 900, true)
		// Add flat energy
		char.AddEnergy("flameforgedinsight", regen)
		// Add Flame Forged insight EM buff mod for 15s (900 frames at 60hz)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 900),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}

	// Proc on opponents only
	buffNoGadget := func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		return buff(args...)
	}

	// Reaction Event hooks, Lunarbloom accepts procs from non-opponents as well
	c.Events.Subscribe(event.OnElectroCharged, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnLunarCharged, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnBloom, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnLunarReactionAttack, buff, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnCrystallizeCryo, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnCrystallizeElectro, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnCrystallizeHydro, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnCrystallizePyro, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnLunarCrystallize, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))

	return w, nil
}
