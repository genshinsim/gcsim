package starsealed

import (
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

const (
	skillEMKey = "starsealed-em-buff"
)

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{char: char}
	multiplier := 0.75 + (0.25 * float64(p.Refine))
	critDmg := 0.24 * multiplier
	skillEM := 48 * multiplier
	stackEM := 24 * multiplier
	stackSscCD := 0.12 * multiplier

	mEM := make([]float64, attributes.EndStatType)
	mSscCD := make([]float64, attributes.EndStatType)
	mCD := make([]float64, attributes.EndStatType)
	mCD[attributes.CD] = critDmg

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("starsealed-cd", -1),
		Amount: func() []float64 {
			return mCD
		},
	})

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		// don't proc if someone else used a skill
		if c.Player.Active() != char.Index() {
			return
		}
		// clear stacks
		w.stacks = 0
		// readds the buff if already active - refreshes it
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(skillEMKey, 20*60),
			Amount: func() []float64 {
				mEM[attributes.EM] = skillEM + stackEM*float64(w.stacks)
				return mEM
			},
		})
	}, "starsealed-on-skill-"+char.Base.Key.String())

	onReaction := func(args ...any) {
		// no stacks if the buff is not active
		if !char.StatModIsActive(skillEMKey) {
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
			"starsealed-reaction-"+char.Base.Key.String(),
		)
	}

	// adds the stacking SSC crit dmg
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("starsealed-ssc-cd", -1),
		Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
			if !char.StatModIsActive(skillEMKey) {
				return nil
			}
			if !w.isStellarRadiance() {
				return nil
			}
			if atk.Info.AttackTag != attacks.AttackTagDirectStellarConduct {
				return nil
			}
			if w.stacks == 0 {
				return nil
			}
			mSscCD[attributes.CD] = stackSscCD * float64(w.stacks)
			return mSscCD
		},
	})

	return w, nil
}

func (w *Weapon) isStellarRadiance() bool {
	return w.char.StatusIsActive(reactable.PolestarFieldKey)
}
