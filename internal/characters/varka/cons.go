package varka

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c4Key          = "varka-c4"
	c6FreeSkillKey = "varka-c6-free-skill"
	c6FreeCAKey    = "varka-c6-free-charge"
	c6Dur          = 3 * 60
)

func (c *char) c1OnSkill() {
	if c.Base.Cons < 1 {
		return
	}
	c.c1Extra = 1
	c.fourWindsChargesAva = 1
}

func (c *char) c1OnSpecialSkill() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}

	if c.c1Extra > 0 {
		c.c1Extra = 0
		return 2.0
	}
	return 1.0
}

func (c *char) c2OnSpecialSkill() {
	if c.Base.Cons < 2 {
		return
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Varka C2",
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       20,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           8,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagVarkaSpecial},
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7)
	// TODO: Get c2 hitmark
	c.Core.QueueAttack(ai, ap, 10, 10)
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnSwirlPyro, c.makeC4CB(attributes.Pyro), c4Key)
	c.Core.Events.Subscribe(event.OnSwirlHydro, c.makeC4CB(attributes.Hydro), c4Key)
	c.Core.Events.Subscribe(event.OnSwirlElectro, c.makeC4CB(attributes.Electro), c4Key)
	c.Core.Events.Subscribe(event.OnSwirlCryo, c.makeC4CB(attributes.Cryo), c4Key)
}

func (c *char) makeC4CB(ele attributes.Element) func(...any) {
	mAnemo := make([]float64, attributes.EndStatType)
	mAnemo[attributes.AnemoP] = 0.2

	dmgP := attributes.EleToDmgP(ele)
	mEle := make([]float64, attributes.EndStatType)
	mEle[dmgP] = 0.2
	return func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c4Key+"-"+attributes.Anemo.String(), 10*60),
				AffectedStat: attributes.AnemoP,
				Amount: func() []float64 {
					return mAnemo
				},
			})
		}

		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c4Key+"-"+ele.String(), 10*60),
				AffectedStat: dmgP,
				Amount: func() []float64 {
					return mEle
				},
			})
		}
	}
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	m := make([]float64, attributes.EndStatType)

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("varka-c6-cdmg", -1),
		AffectedStat: attributes.AnemoP,
		Amount: func() []float64 {
			if c.a4Stacks == 0 {
				return nil
			}
			m[attributes.CD] = float64(c.a4Stacks) * 0.2
			return m
		},
	})
}

func (c *char) c6OnSkill() {
	if c.Base.Cons < 6 {
		return
	}

	c.AddStatus(c6FreeCAKey, c6Dur, true)
}

func (c *char) c6OnSkillCA() {
	if c.Base.Cons < 6 {
		return
	}

	c.AddStatus(c6FreeSkillKey, c6Dur, true)
}

func (c *char) c6FreeCA() bool {
	if c.Base.Cons < 6 {
		return false
	}

	// does this also need to be in E state?
	return c.StatusIsActive(c6FreeCAKey)
}

func (c *char) c6FreeSkill() bool {
	if c.Base.Cons < 6 {
		return false
	}

	// does this also need to be in E state?
	return c.StatusIsActive(c6FreeSkillKey)
}
