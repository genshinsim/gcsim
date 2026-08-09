package citlali

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	iceStormHitmark = 118

	// the logic is cast -> 1s delay -> nyxadd+initial hit
	// -> 0.95s + random [0.5,1.0] delay -> 2nd nyxadd+burst hit
	// From testing, the random delay is uniform
	// const is relative to skull summoning- burst hitmark on q, field hitmark for c4
	spiritVesselSkullHitmarkBase = 87

	iceStormAbil = "Ice Storm DMG"
)

var burstFrames []int

func (c *char) spiritVesselSkullRandHitmark() int {
	return spiritVesselSkullHitmarkBase + c.Core.Rand.Intn(31)
}

func init() {
	burstFrames = frames.InitAbilSlice(113) // Q -> Jump
	burstFrames[action.ActionAttack] = 112
	burstFrames[action.ActionCharge] = 112
	burstFrames[action.ActionSkill] = 111
	burstFrames[action.ActionDash] = 112
	burstFrames[action.ActionWalk] = 112
	burstFrames[action.ActionSwap] = 110
}

func (c *char) Burst(_ map[string]int) (action.Info, error) {
	aiIceStorm := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           iceStormAbil,
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     50,
		Mult:           iceStorm[c.TalentLvlBurst()],
		FlatDmg:        c.a4Dmg(iceStormAbil),
	}
	aiSpiritVesselSkull := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Spiritvessel Skull DMG",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagCitlaliSpiritVessel,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     25,
		Mult:           spiritVessel[c.TalentLvlBurst()],
	}

	// with delay
	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, 15*60)
	c.QueueCharTask(func() {
		c.generateNightsoulPoints(24)
	}, 115)

	// initial hit
	c.Core.QueueAttack(aiIceStorm, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.5), iceStormHitmark, iceStormHitmark)

	// skull hits
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), nil)
	for i, enemy := range enemies {
		if i > 2 {
			break
		}
		c.QueueCharTask(func() {
			c.generateNightsoulPoints(3.0)
			c.Core.QueueAttack(aiSpiritVesselSkull, combat.NewCircleHitOnTarget(enemy.Pos(), nil, 3.5), 0, 0)
		}, iceStormHitmark+c.spiritVesselSkullRandHitmark())
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}
