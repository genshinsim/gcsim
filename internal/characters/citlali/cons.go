package citlali

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	maxC6Stacks = 40
	c4SkullIcd  = "c4-skull-icd"
)

// Additionally, when Citlali is using her leap, or is Aiming or using her
// Charged Attack in mid-air, her Phlogiston consumption is decreased by 45%.
// NOT IMPLEMENTED
func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if c.Index == atk.Info.ActorIndex {
			return false
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagElementalBurst:
		default:
			return false
		}
		if c.numStellarBlades > 0 {
			em := c.Stat(attributes.EM)
			amt := em * 2
			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Citlali C1 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
					Write("before", atk.Info.FlatDmg).
					Write("addition", amt).
					Write("Stellar Blades left", c.numStellarBlades)
			}
			atk.Info.FlatDmg += amt
			c.numStellarBlades--
		}
		return false
	}, "citlali-c1-on-dmg")
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	buffSelf := make([]float64, attributes.EndStatType)
	buffSelf[attributes.EM] = 125
	buffOther := make([]float64, attributes.EndStatType)
	buffOther[attributes.EM] = 250

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("citlali-c2-em", -1),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.Shields.Get(shield.CitlaliSkill) == nil {
				return nil, false
			}
			return buffSelf, true
		},
	})

	chars := c.Core.Player.Chars()
	for _, char := range chars {
		if c.Index == char.Index {
			continue
		}
		this := char
		this.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("citlali-c2-em", -1),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				// character should be followed by Itzpapa, i.e. the character is active
				if c.Core.Player.Active() != this.Index {
					return nil, false
				}
				if c.Core.Player.Shields.Get(shield.CitlaliSkill) == nil && !c.nightsoulState.HasBlessing() {
					return nil, false
				}
				return buffOther, true
			},
		})
	}
}

func (c *char) c4SkullCB(a combat.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}
	if c.StatusIsActive(c4SkullIcd) {
		return
	}
	c.AddStatus(c4SkullIcd, 8*60, true)

	c.generateNightsoulPoints(16)
	c.AddEnergy("citlali-c4", 8)
	aiSpiritVesselSkull := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Spiritvessel Skull DMG (C4)",
		AttackTag:      attacks.AttackTagNone,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     25,
		FlatDmg:        18 * c.Stat(attributes.EM),
	}
	// TODO: the actual hitmark
	hitmark := spiritVesselSkullHitmark - iceStormHitmark
	c.Core.QueueAttack(
		aiSpiritVesselSkull,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5),
		hitmark,
		hitmark,
	)
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	buffSelf := make([]float64, attributes.EndStatType)
	buffOther := make([]float64, attributes.EndStatType)

	chars := c.Core.Player.Chars()
	for _, char := range chars {
		if c.Index == char.Index {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("citlali-c6", -1),
			Amount: func() ([]float64, bool) {
				buffOther[attributes.PyroP] = 0.015 * c.numC6Stacks
				buffOther[attributes.HydroP] = 0.015 * c.numC6Stacks
				return buffOther, true
			},
		})
	}
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("citlali-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			buffSelf[attributes.DmgP] = 0.025 * c.numC6Stacks
			return buffSelf, true
		},
	})
}
