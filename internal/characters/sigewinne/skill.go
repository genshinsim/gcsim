package sigewinne

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/internal/template/sourcewaterdroplet"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames [][]int

const (
	skillPressCDStart     = 16
	skillPressHitmark     = 35
	skillShortHoldCDStart = 40
	skillShortHoldHitmark = 66
	skillHoldCDStart      = 66
	skillHoldHitmark      = 90
	skillCD               = 18

	bubbleHitInterval = 107
	bubbleRadius      = 1
	bubbleTierBuff    = 0.05
	skillKey          = "sigewinne-skill"

	skillDropletOffset          = 0.5
	skillDropletSpawnTimeOffset = 4

	skillAlignedICD     = 10 * 60
	skillAlignedHitmark = 40
	skillAlignedICDKey  = "sigewinne-aligned-icd"

	particleCount  = 4
	particleICDKey = "sigewinne-particle-icd"

	hpDebtEnergyRatio = 2000.
)

func init() {
	skillFrames = make([][]int, 3)
	// skill (press) -> x
	skillFrames[0] = frames.InitAbilSlice(41) // burst
	skillFrames[0][action.ActionAttack] = 39
	skillFrames[0][action.ActionCharge] = 40
	skillFrames[0][action.ActionWalk] = 40

	// skill (short hold) -> x
	skillFrames[1] = frames.InitAbilSlice(56)

	// skill (hold) -> x
	skillFrames[2] = frames.InitAbilSlice(89) // na
	skillFrames[2][action.ActionWalk] = 86
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.burstEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Super Saturated Syringing with Elemental Skill", c.Base.Key)
	}

	// TODO: rework it like kirara/sayu skill?
	skillHitmark := skillPressHitmark
	skillCDStart := skillPressCDStart
	c.currentBubbleTier = 0
	hold, ok := p["hold"]
	if !ok {
		hold = 0
	}
	if hold == 1 {
		skillHitmark = skillShortHoldHitmark
		skillCDStart = skillShortHoldCDStart
		c.currentBubbleTier = 1
	} else if hold == 2 {
		skillHitmark = skillHoldHitmark
		skillCDStart = skillHoldCDStart
		c.currentBubbleTier = 2
	}

	c.generateSkillSnapshot()
	c.AddStatus(skillKey, -1, false)
	c.particleGenerated = false
	c.lastSummonSrc = c.Core.F

	c.SetCDWithDelay(action.ActionSkill, skillCD*60, skillCDStart)
	c.Core.Tasks.Add(c.spawnDroplets, skillCDStart+skillDropletSpawnTimeOffset)
	c.Core.Tasks.Add(c.bolsteringBubblebalm(c.lastSummonSrc, 0), skillHitmark)

	if c.Base.Ascension >= 1 {
		c.a1Self()
	}
	if c.Base.Cons >= 2 {
		c.Core.Tasks.Add(c.addC2Shield(skillFrames[hold][action.ActionWalk]), 1)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[hold]),
		AnimationLength: skillFrames[hold][action.InvalidAction],
		CanQueueAfter:   skillFrames[hold][action.ActionAttack],
		State:           action.SkillState,
	}, nil
}

func (c *char) bolsteringBubblebalm(src, tick int) func() {
	return func() {
		if src != c.lastSummonSrc {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}

		// Damage
		target := c.Core.Combat.PrimaryTarget()
		c.Core.QueueAttackWithSnap(
			c.skillAttackInfo,
			c.skillSnapshot,
			combat.NewCircleHitOnTarget(target, nil, bubbleRadius),
			0,
			c.particleCB,
		)
		c.surgingBladeTask(target)
		c.bubbleTierLoseTask(tick)

		// Healing
		c.bubbleHealing()

		if tick == c.bubbleHitLimit-1 {
			c.bubbleFinalHealing()
			c.DeleteStatus(skillKey)
			return
		}

		// TODO: hitlag affected?
		c.Core.Tasks.Add(c.bolsteringBubblebalm(src, tick+1), bubbleHitInterval)
	}
}

func (c *char) spawnDroplets() {
	player := c.Core.Combat.Player()
	for j := 0; j < 2; j++ {
		pos := geometry.CalcRandomPointFromCenter(
			geometry.CalcOffsetPoint(
				player.Pos(),
				geometry.Point{Y: 1.5},
				player.Direction(),
			),
			0.3,
			1,
			c.Core.Rand,
		)
		sourcewaterdroplet.New(c.Core, pos, combat.GadgetTypSourcewaterDropletSigewinne)
	}
}

func (c *char) surgingBladeTask(target combat.Target) {
	if c.StatusIsActive(skillAlignedICDKey) {
		return
	}
	c.AddStatus(skillAlignedICDKey, skillAlignedICD, true)

	aiThorn := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Spiritbreath Thorn (" + c.Base.Key.Pretty() + ")",
		AttackTag:    attacks.AttackTagElementalArt,
		ICDTag:       attacks.ICDTagNone,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypePierce,
		Element:      attributes.Hydro,
		Durability:   0,
		FlatDmg:      surgingBladeDMG[c.TalentLvlSkill()] * c.MaxHP(),
		HitlagFactor: 0.01,
	}
	c.Core.QueueAttack(
		aiThorn,
		combat.NewCircleHitOnTarget(target, nil, 3),
		skillAlignedHitmark,
		skillAlignedHitmark,
	)
}

func (c *char) bubbleHealing() {
	if !c.StatusIsActive(skillKey) {
		return
	}

	// heal everyone except Sigewinne
	for _, other := range c.Core.Player.Chars() {
		if other.Index == c.Index {
			continue
		}
		skillBonus := float64(c.currentBubbleTier) * bubbleTierBuff
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index,
			Target:  other.Index,
			Message: "Bolstering Bubblebalm Healing",
			Src:     bolsteringBubblebalmHealingPct[c.TalentLvlSkill()]*c.MaxHP() + bolsteringBubblebalmHealingFlat[c.TalentLvlSkill()],
			Bonus:   c.Stat(attributes.Heal) + skillBonus,
		})
	}
}

func (c *char) bubbleFinalHealing() {
	if !c.StatusIsActive(skillKey) {
		return
	}
	// heal only Sigewinne
	skillBonus := float64(c.currentBubbleTier) * bubbleTierBuff
	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: "Bolstering Bubblebalm Healing",
		Src:     finalBounceHealing[c.TalentLvlSkill()] * c.MaxHP(),
		Bonus:   c.Stat(attributes.Heal) + skillBonus,
	})
}

func (c *char) bubbleTierDamageMod() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("sigewinne-bubble-tier", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt:
			case attacks.AttackTagElementalArtHold:
			default:
				return nil, false
			}
			if c.currentBubbleTier == 0 {
				return nil, false
			}
			if atk.Info.Abil != c.skillAttackInfo.Abil {
				return nil, false
			}
			m[attributes.DmgP] = float64(c.currentBubbleTier) * bubbleTierBuff
			return m, true
		},
	})
}

func (c *char) energyBondClearMod() {
	// TODO: override healing functions?
	c.Core.Events.Subscribe(event.OnHPDebt, func(args ...interface{}) bool {
		index := args[0].(int)
		if index != c.Index {
			return false
		}
		debtChange := args[1].(float64)
		if debtChange < 0 {
			c.collectedHpDebt += -float32(debtChange)
		}
		if c.CurrentHPDebt() > 0 {
			return false
		}

		energyAmt := min(5., c.collectedHpDebt/hpDebtEnergyRatio)
		c.collectedHpDebt = 0
		c.AddEnergy("sigewinne-skill", float64(energyAmt))
		return false
	}, "sigewinne-hpdebt-hook")
}

func (c *char) bubbleTierLoseTask(tick int) {
	if c.Base.Cons < 1 || tick > 2 {
		c.currentBubbleTier--
		c.currentBubbleTier = max(c.currentBubbleTier, 0)
	}
}

func (c *char) particleCB(ac combat.AttackCB) {
	if ac.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.particleGenerated {
		return
	}

	// once per skill
	c.particleGenerated = true
	c.Core.QueueParticle(c.Base.Key.String(), particleCount, attributes.Hydro, c.ParticleDelay)
}

func (c *char) generateSkillSnapshot() {
	c.skillAttackInfo = combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Rebound Hydrotherapy",
		AttackTag:    attacks.AttackTagElementalArt,
		ICDTag:       attacks.ICDTagElementalArt,
		ICDGroup:     attacks.ICDGroupSigewinne,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Hydro,
		Durability:   25,
		FlatDmg:      bolsteringBubblebalmDMG[c.TalentLvlSkill()] * c.MaxHP(),
		HitlagFactor: 0.02,
	}
	c.skillSnapshot = c.Snapshot(&c.skillAttackInfo)
}
