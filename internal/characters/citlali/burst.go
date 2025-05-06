package citlali

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const (
	iceStormHitmark          = 118
	spiritVesselSkullHitmark = 210

	iceStormAbil = "Ice Storm DMG"
)

var (
	burstFrames []int
)

func init() {
	burstFrames = frames.InitAbilSlice(112) // Q -> N1
	burstFrames[action.ActionSkill] = 111
	burstFrames[action.JumpState] = 113
	burstFrames[action.ActionSwap] = 110
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	aiIceStorm := combat.AttackInfo{
		ActorIndex:     c.Index,
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
	aiSpiritVesselSkull := combat.AttackInfo{
		ActorIndex:     c.Index,
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
	c.ConsumeEnergy(5)
	c.QueueCharTask(func() {
		c.SetCD(action.ActionBurst, 15*60)
	}, 1)
	c.QueueCharTask(func() {
		c.nightsoulState.GeneratePoints(24)
		c.tryEnterOpalFireState(c.itzpapaSrc)
	}, 115)

	// initial hit
	c.Core.QueueAttack(aiIceStorm, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.5), iceStormHitmark, iceStormHitmark)

	// skull hits
	c.QueueCharTask(func() {
		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), nil)
		for i, enemy := range enemies {
			if i > 2 {
				break
			}
			c.Core.QueueAttack(aiSpiritVesselSkull, combat.NewCircleHitOnTarget(enemy.Pos(), nil, 3.5), 0, 0)
		}
		c.nightsoulState.GeneratePoints(float64(3 * min(3, len(enemies))))
		c.tryEnterOpalFireState(c.itzpapaSrc)
	}, spiritVesselSkullHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}
