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
	Index       int
	c           *core.Core
	self        *character.CharWrapper
	critDMG []float64
	SkillEM []float64
	stacks int
	EMStack float64
	SSCCDMGStack float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
	c:    c,
	self: char,
	}
	r := p.Refine

	multiplier := 0.75 + (0.25 * float64(r))

	w.critDMG = make([]float64, attributes.EndStatType)
	w.critDMG[attributes.CD] = 0.24 * multiplier

	w.SkillEM = make([]float64, attributes.EndStatType)
	w.SkillEM[attributes.EM] = 48 * multiplier

	w.EMStack = 24 * multiplier
	w.SSCCDMGStack = 0.12 * multiplier

	char.AddStatMod(character.StatMod{
			Base: modifier.NewBase("Starsealed", -1),	
			Amount: func() []float64 {
				return w.critDMG
			},
		})
		
	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		// don't proc if someone else used a skill
		if c.Player.Active() != char.Index() {
			return
		}

		// activates the passive stacks + clears them
		w.stacks = 0
		char.AddStatus("active-stacks", 20*60, true)

		// add buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("starsealed-skill-buff", 20*60),
			Amount: func() []float64 {
				return w.SkillEM
			},
		})
	}, fmt.Sprintf("Starsealed-on-skill-%v", char.Base.Key.String()))

	onReaction := func(args ...any) {
		// no stacks if the buff is not active
		if !char.StatusIsActive("active-stacks") {
			return
		}
		//adds stack on team reaction, max 3 stacks
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
	//adds the stacking EM
	m := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
    	Base: modifier.NewBaseWithHitlag("Starsealed-stack-buff", -1),
    	Amount: func() []float64 {
			if !char.StatusIsActive("active-stacks") {
				return nil
			}
        	m[attributes.EM] = w.EMStack * float64(w.stacks)
        	return m
    		},
		})
		fmt.Sprintf("Starsealed-stack-buff-%v", char.Base.Key.String())

	//adds the stacking SSC crit dmg	
	sscCrit := make([]float64, attributes.EndStatType)
	char.AddAttackMod(character.AttackMod{
    	Base: modifier.NewBaseWithHitlag("starsealed-ssc-crit", -1),
    	Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
			if !char.StatusIsActive("active-stacks") {
				return nil
			}
        	if !w.isSSC() || atk.Info.AttackTag != attacks.AttackTagDirectStellarConduct {
            	return nil
        	}
        	sscCrit[attributes.CD] = w.SSCCDMGStack * float64(w.stacks)
        	return sscCrit
    	},
	})

	return w, nil
}

func (w *Weapon) isSSC() bool {
	return w.self.StatusIsActive(reactable.PolestarFieldKey)
}
