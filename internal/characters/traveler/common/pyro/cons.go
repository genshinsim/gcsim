package pyro

import (
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1AttackModKey = "travelerpyro-c1"
	c2StatusKey    = "travelerpyro-c2-key"
	c6AttackModKey = "travelerpyro-c6"
)

func (c *Traveler) c1AddMod() {
	if c.Base.Cons < 1 {
		return
	}
	mDmg := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		this := char
		this.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c1AttackModKey, -1),
			Amount: func(ae *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				// char must be active
				if c.Core.Player.Active() != this.Index {
					return nil, false
				}
				mDmg[attributes.DmgP] = 0.06
				if this.StatusIsActive(nightsoul.NightsoulBlessingStatus) {
					mDmg[attributes.DmgP] += 0.06
				}
				return mDmg, true
			},
		})
	}
}

func (c *Traveler) c2() {
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

func (c *Traveler) c6AddMod() {
	if c.Base.Cons < 6 {
		return
	}
	mCD := make([]float64, attributes.EndStatType)
	mCD[attributes.CD] = 0.4
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c6AttackModKey, -1),
		Amount: func(ae *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			switch ae.Info.AttackTag {
			case attacks.AttackTagNormal:
			case attacks.AttackTagExtra:
			case attacks.AttackTagPlunge:
			default:
				return nil, false
			}
			return mCD, true
		},
	})
}
