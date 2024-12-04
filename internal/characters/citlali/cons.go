package citlali

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
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
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if c.Index == atk.Info.ActorIndex {
			return false
		}
		if c.numStellarBlades > 0 {
			em := c.NonExtraStat(attributes.EM)
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
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("citlali-c2-em", 20*60),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			buffSelf := make([]float64, attributes.EndStatType)
			buffSelf[attributes.EM] = 125
			return buffSelf, true
		},
	})
	chars := c.Core.Player.Chars()
	for _, char := range chars {
		if c.Index == char.Index {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("citlali-c2-em", 20*60),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				if c.Core.Player.Shields.CharacterIsShielded(char.Index, c.Core.Player.Active()) {
					return nil, false
				}
				buffOther := make([]float64, attributes.EndStatType)
				buffOther[attributes.EM] = 250
				return buffOther, true
			},
		})
	}
}

func (c *char) c4Skull() {
	if c.Base.Cons < 4 {
		return
	}
	if c.StatusIsActive(c4SkullIcd) {
		return
	}
	c.AddStatus(c4SkullIcd, 8*60, false)
	c.nightsoulState.GeneratePoints(16)
	c.tryEnterOpalFireState(c.itzpapaSrc)
	c.AddEnergy("citlali-c4-energy", 8)
	aiSpiritVesselSkull := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Spiritvessel Skull DMG (C4)",
		AttackTag:      attacks.AttackTagNone,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagElementalBurst, // TODO: check this
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     25,
		FlatDmg:        18 * c.NonExtraStat(attributes.EM),
	}
	// TODO: the actual hitmark
	c.Core.QueueAttack(aiSpiritVesselSkull, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()),
		spiritVesselSkullHitmark-iceStormHitmark, spiritVesselSkullHitmark-iceStormHitmark)
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	chars := c.Core.Player.Chars()
	for _, char := range chars {
		if c.Index == char.Index {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("citlali-c6", -1),
			Amount: func() ([]float64, bool) {
				buffC6 := make([]float64, attributes.EndStatType)
				buffC6[attributes.PyroP] = 0.015 * float64(c.numC6Stacks)
				buffC6[attributes.HydroP] = 0.015 * float64(c.numC6Stacks)
				return buffC6, true
			},
		})
	}
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("citlali-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			buffC6 := make([]float64, attributes.EndStatType)
			buffC6[attributes.DmgP] = 0.025 * float64(c.numC6Stacks)
			return buffC6, true
		},
	})
}
