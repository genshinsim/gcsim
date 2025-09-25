package pyro

import (
	"slices"

	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1AttackModKey = "travelerpyro-c1"
	c2StatusKey    = "travelerpyro-c2"
	c6AttackModKey = "travelerpyro-c6"
)

func (c *Traveler) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	mDmg := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		this := char
		this.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c1AttackModKey, -1),
			Amount: func(ae *info.AttackEvent, _ info.Target) ([]float64, bool) {
				// char must be active
				if c.Core.Player.Active() != this.Index() {
					return nil, false
				}
				if !c.nightsoulState.HasBlessing() {
					return nil, false
				}
				mDmg[attributes.DmgP] = 0.06
				if this.StatusIsActive(nightsoul.NightsoulBlessingStatus) {
					mDmg[attributes.DmgP] += 0.09
				}
				return mDmg, true
			},
		})
	}
}

func (c *Traveler) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	fReactionHook := func(args ...any) bool {
		if !c.StatusIsActive(c2StatusKey) {
			return false
		}

		if c.c2ActivationsPerSkill < 2 {
			c.c2ActivationsPerSkill++
			c.nightsoulState.GeneratePoints(14)
		}

		return false
	}

	c.Core.Events.Subscribe(event.OnBurning, fReactionHook, "travelerpyro-c2-onburning")
	c.Core.Events.Subscribe(event.OnVaporize, fReactionHook, "travelerpyro-c2-onvaporize")
	c.Core.Events.Subscribe(event.OnMelt, fReactionHook, "travelerpyro-c2-onmelt")
	c.Core.Events.Subscribe(event.OnOverload, fReactionHook, "travelerpyro-c2-onoverload")
	c.Core.Events.Subscribe(event.OnBurgeon, fReactionHook, "travelerpyro-c2-onburgeon")
	c.Core.Events.Subscribe(event.OnSwirlPyro, fReactionHook, "travelerpyro-c2-onswirlpyro")
	c.Core.Events.Subscribe(event.OnCrystallizePyro, fReactionHook, "travelerpyro-c2-oncrystallizepyro")
}

func (c *Traveler) c2OnSkill() {
	if c.Base.Cons < 2 {
		return
	}
	c.AddStatus(c2StatusKey, 12*60, true)
	c.c2ActivationsPerSkill = 0
}

func (c *Traveler) c4AddMod() {
	if c.Base.Cons < 4 {
		return
	}
	mPyroDmg := make([]float64, attributes.EndStatType)
	mPyroDmg[attributes.PyroP] = 0.2
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag("travelerpyro-c4", 9*60),
		Amount: func() ([]float64, bool) {
			return mPyroDmg, true
		},
	})
}

func (c *Traveler) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	mCD := make([]float64, attributes.EndStatType)
	mCD[attributes.CD] = 0.4
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c6AttackModKey, -1),
		Amount: func(ae *info.AttackEvent, _ info.Target) ([]float64, bool) {
			switch ae.Info.AttackTag {
			case attacks.AttackTagNormal:
			case attacks.AttackTagExtra:
			case attacks.AttackTagPlunge:
			default:
				return nil, false
			}
			if !slices.Contains(ae.Info.AdditionalTags, attacks.AdditionalTagNightsoul) {
				return nil, false
			}
			return mCD, true
		},
	})
}
