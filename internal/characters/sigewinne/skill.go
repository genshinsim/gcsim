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
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	skillFrames              [][]int
	skillDropletRandomRanges = [][]float64{{0.5, 1.}, {0.5, 1.}}
)

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
	SkillDropletSpawnTimeOffset = 4

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
	skillFrames[0] = frames.InitAbilSlice(40) // walk
	skillFrames[0][action.ActionAttack] = 40
	skillFrames[0][action.ActionCharge] = 40
	skillFrames[0][action.ActionBurst] = 40
	skillFrames[0][action.ActionDash] = 40
	skillFrames[0][action.ActionJump] = 40
	skillFrames[0][action.ActionSwap] = 40

	// skill (short hold) -> x
	skillFrames[1] = frames.InitAbilSlice(56) // walk
	skillFrames[1][action.ActionAttack] = 56
	skillFrames[1][action.ActionCharge] = 56
	skillFrames[1][action.ActionBurst] = 56
	skillFrames[1][action.ActionDash] = 56
	skillFrames[1][action.ActionJump] = 56
	skillFrames[1][action.ActionSwap] = 56

	// skill (hold) -> x
	skillFrames[2] = frames.InitAbilSlice(86) // walk
	skillFrames[2][action.ActionAttack] = 86
	skillFrames[2][action.ActionCharge] = 86
	skillFrames[2][action.ActionBurst] = 86
	skillFrames[2][action.ActionDash] = 86
	skillFrames[2][action.ActionJump] = 86
	skillFrames[2][action.ActionSwap] = 86
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.burstEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Super Saturated Syringing with Elemental Skill", c.CharWrapper.Base.Key)
	}
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

	if c.Base.Ascension >= 1 {
		c.AddStatus(convalescenceKey, skillCD*60, false)
		c.SetTag(convalescenceKey, 10)
		buff := make([]float64, attributes.EndStatType)
		buff[attributes.HydroP] = a1DmgBuff
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("sigewinne-a1-hydro-percent", skillCD*60),
			Amount: func(a *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				return buff, true
			},
		})
	}

	c.AddStatus(skillKey, -1, false)
	src := c.Core.F
	c.lastSummonSrc = src
	c.Core.Log.NewEvent("Summoned Bostering Bubblebalm", glog.LogCharacterEvent, c.Index)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rebound Hydrotherapy",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArtHydro,
		ICDGroup:   attacks.ICDGroupSigewinne,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    bolsteringBubblebalmDMG[c.TalentLvlSkill()] * c.MaxHP(),
	}

	c.SetCDWithDelay(action.ActionSkill, skillCD*60, skillCDStart)
	c.Core.Tasks.Add(c.spawnDroplets(), skillCDStart+SkillDropletSpawnTimeOffset)

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, bubbleRadius), skillHitmark, skillHitmark, c.particleCB)
	c.Core.Tasks.Add(c.bolsteringBubblebalm(src, 1), skillHitmark+bubbleHitInterval)
	c.Core.Tasks.Add(c.surgingBladeTask(), skillHitmark)
	c.Core.Tasks.Add(c.bubbleTierLoseTask(0), skillHitmark+1)

	// Healing
	c.Core.Tasks.Add(c.bubbleHealing(src), skillHitmark)

	if c.Base.Cons >= 2 {
		c.Core.Tasks.Add(c.addC2Shield, 1)
		c.Core.Tasks.Add(c.removeC2Shield, skillFrames[hold][action.ActionWalk])
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[hold]),
		AnimationLength: skillFrames[hold][action.InvalidAction],
		CanQueueAfter:   skillFrames[hold][action.ActionDash],
		State:           action.SkillState,
	}, nil
}

func (c *char) bolsteringBubblebalm(src, tick int) func() {
	return func() {
		if tick >= c.bubbleHitLimit {
			c.DeleteStatus(skillKey)
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		if src != c.lastSummonSrc {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Rebound Hydrotherapy",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArtHydro,
			ICDGroup:   attacks.ICDGroupSigewinne,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    bolsteringBubblebalmDMG[c.TalentLvlSkill()] * c.MaxHP(),
		}

		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, bubbleRadius), 0, 0, c.particleCB)
		c.Core.Tasks.Add(c.bolsteringBubblebalm(src, tick+1), bubbleHitInterval)
		c.Core.Tasks.Add(c.surgingBladeTask(), bubbleHitInterval)
		c.Core.Tasks.Add(c.bubbleTierLoseTask(tick), bubbleHitInterval+1)

		// Healing
		if c.bubbleHitLimit == tick-1 {
			c.Core.Tasks.Add(c.bubbleFinalHealing(src), bubbleHitInterval)
		}
		c.Core.Tasks.Add(c.bubbleHealing(src), bubbleHitInterval)
	}
}

func (c *char) spawnDroplets() func() {
	return func() {
		player := c.Core.Combat.Player()
		playerX := player.Pos().X
		playerY := player.Pos().Y
		targetX := c.Core.Combat.PrimaryTarget().Pos().X
		targetY := c.Core.Combat.PrimaryTarget().Pos().Y
		vecX := targetX - playerX
		vecY := targetY - playerY
		dropletOffsetX := vecX * skillDropletOffset
		dropletOffsetY := vecY * skillDropletOffset
		for j := 0; j < 2; j++ {
			sourcewaterdroplet.New(
				c.Core,
				geometry.CalcRandomPointFromCenter(
					geometry.CalcOffsetPoint(
						player.Pos(),
						geometry.Point{X: dropletOffsetX, Y: dropletOffsetY},
						player.Direction(),
					),
					skillDropletRandomRanges[j][0],
					skillDropletRandomRanges[j][1],
					c.Core.Rand,
				),
				combat.GadgetTypSourcewaterDropletNeuv,
			)
			c.Core.Combat.Log.NewEvent("Skill: Spawned 3 droplets", glog.LogCharacterEvent, c.Index)
		}
	}
}

func (c *char) surgingBladeTask() func() {
	return func() {
		aiThorn := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Spiritbreath Thorn (" + c.Base.Key.Pretty() + ")",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagNone,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSpear,
			Element:            attributes.Hydro,
			Durability:         0,
			FlatDmg:            surgingBladeDMG[c.TalentLvlSkill()] * c.MaxHP(),
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
		}

		if c.StatusIsActive(skillAlignedICDKey) {
			return
		}
		c.AddStatus(skillAlignedICDKey, skillAlignedICD, true)

		skillPos := c.Core.Combat.PrimaryTarget().Pos()
		c.Core.QueueAttack(
			aiThorn,
			combat.NewCircleHitOnTarget(skillPos, nil, 3),
			skillAlignedHitmark,
			skillAlignedHitmark,
		)
	}
}

func (c *char) bubbleHealing(src int) func() {
	return func() {
		if src != c.lastSummonSrc {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		// heal everyone except Sigewinne
		for _, other := range c.Core.Player.Chars() {
			if other.Index == c.Index {
				continue
			}
			c.Core.Player.Heal(info.HealInfo{
				Caller:  c.Index,
				Target:  other.Index,
				Message: "Bolstering Bubblebalm Healing",
				Src:     bolsteringBubblebalmHealingPct[c.TalentLvlSkill()]*c.MaxHP() + bolsteringBubblebalmHealingFlat[c.TalentLvlSkill()],
				Bonus:   c.Stat(attributes.Heal),
			})
		}
	}
}

func (c *char) bubbleFinalHealing(src int) func() {
	return func() {
		if src != c.lastSummonSrc {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}
		// heal only Sigewinne
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "Bolstering Bubblebalm Healing",
			Src:     finalBounceHealing[c.TalentLvlSkill()] * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
	}
}

func (c *char) bubbleTierDamageMod() {
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("sigewinne-bubble-tier-damage-buff", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if !c.StatusIsActive(skillKey) {
				return nil, false
			}
			if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
				return nil, false
			}
			dmgAmt := make([]float64, attributes.EndStatType)
			dmgAmt[attributes.DmgP] = min(2., float64(c.currentBubbleTier)) * bubbleTierBuff
			return dmgAmt, true
		},
	})

	c.AddHealBonusMod(character.HealBonusMod{
		Base: modifier.NewBase("sigewinne-bubble-tier-heal-buff", -1),
		Amount: func() (float64, bool) {
			if c.StatusIsActive(skillKey) {
				return min(2, float64(c.currentBubbleTier)) * bubbleTierBuff, true
			}
			return 0, false
		},
	})
}

func (c *char) energyBondClearMod() {
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
		c.AddEnergy("sigewinne-hpdebt-clear-energy", float64(energyAmt))
		return false
	}, "sigewinne-hpdebt-hook")
}

func (c *char) bubbleTierLoseTask(tick int) func() {
	return func() {
		if c.Base.Cons < 1 || tick > 2 {
			c.currentBubbleTier--
			c.currentBubbleTier = max(c.currentBubbleTier, 0)
		}
	}
}

func (c *char) particleCB(ac combat.AttackCB) {
	if ac.Target.Type() != targets.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, skillCD*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), particleCount, attributes.Hydro, c.ParticleDelay)
}
