package starsealed

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type Weapon struct {
	Index  int
	stacks int
	char   *character.CharWrapper
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{char: char}
	multiplier := 0.75 + (0.25 * float64(p.Refine))
	critdmg := 0.24 * multiplier
	skillem := 48 * multiplier
	stackem := 24 * multiplier
	stackssccd := 0.12 * multiplier
	constantCD := make([]float64, attributes.EndStatType)
	constantCD[attributes.CD] = critdmg
	skillEM := make([]float64, attributes.EndStatType)
	skillEM[attributes.EM] = skillem
	stackEM := make([]float64, attributes.EndStatType)

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("Starsealed-constant-CDMG", -1),
		Amount: func() []float64 {
			return constantCD
		},
	})

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		// don't proc if someone else used a skill
		if c.Player.Active() != char.Index() {
			return
		}
		// clear stacks
		w.stacks = 0
		// refresh duration if already active
		if char.StatModIsActive("starsealed-skill-buff") {
			char.ExtendStatus(
				"starsealed-skill-buff",
				20*60-char.StatusDuration("starsealed-skill-buff"),
			)
			return
		}
		// add EM buff
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("starsealed-skill-buff", 20*60),
			Amount: func() []float64 {
				return skillEM
			},
		})
	}, fmt.Sprintf("Starsealed-on-skill-%v", char.Base.Key.String()))

	onReaction := func(args ...any) {
		// no stacks if the buff is not active
		if !char.StatModIsActive("starsealed-skill-buff") {
			return
		}
		// adds stack on team reaction, max 3 stacks
		if w.stacks < 3 {
			w.stacks++
		}
	}
	// activates the stacking function when a reaction happens
	for e := event.ReactionEventStartDelim + 1; e < event.ReactionEventEndDelim; e++ {
		c.Events.Subscribe(
			e,
			onReaction,
			fmt.Sprintf("starsealed-reaction-%s-%d", char.Base.Key.String(), e),
		)
	}
	// adds the stacking EM
	char.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag("Starsealed-stack-buff", -1),
		Amount: func() []float64 {
			if !char.StatModIsActive("starsealed-skill-buff") {
				return nil
			}
			stackEM[attributes.EM] = stackem * float64(w.stacks)
			return stackEM
		},
	})

	// adds the stacking SSC crit dmg
	sscCrit := make([]float64, attributes.EndStatType)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("starsealed-ssc-crit", -1),
		Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
			if !char.StatModIsActive("starsealed-skill-buff") {
				return nil
			}
			if !w.isSSC() || atk.Info.AttackTag != attacks.AttackTagDirectStellarConduct {
				return nil
			}
			sscCrit[attributes.CD] = stackssccd * float64(w.stacks)
			return sscCrit
		},
	})

	return w, nil
}

func (w *Weapon) isSSC() bool {
	return w.char.StatusIsActive(reactable.PolestarFieldKey)
}
