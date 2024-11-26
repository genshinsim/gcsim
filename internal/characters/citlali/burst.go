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
	spiritVesselSkullHitmark = 223

	iceStormAbil = "Ice Storm DMG"
)

var (
	burstFrames []int
)

func init() {
	burstFrames = frames.InitAbilSlice(133) // Q -> Swap
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
		Mult:           9.677,
		FlatDmg:        c.a4Dmg(iceStormAbil),
	}
	aiSpiritVesselSkull := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Spiritvessel Skull DMG",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagElementalBurst, // TODO: check this
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     25,
		Mult:           2.419,
	}
	c.ConsumeEnergy(5)
	c.SetCD(action.ActionBurst, 15*60)
	c.nightsoulState.GeneratePoints(24)
	c.ActivateItzpapa(c.itzpapaSrc)
	c.Core.QueueAttack(aiIceStorm, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.5), iceStormHitmark, iceStormHitmark)
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), nil)
	c.QueueCharTask(func() {
		c.nightsoulState.GeneratePoints(float64(3 * len(enemies)))
		c.ActivateItzpapa(c.itzpapaSrc)
	}, spiritVesselSkullHitmark)
	for _, enemy := range enemies {
		c.Core.QueueAttack(aiSpiritVesselSkull, combat.NewSingleTargetHit(enemy.Key()), spiritVesselSkullHitmark, spiritVesselSkullHitmark)
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionBurst],
		State:           action.BurstState,
	}, nil
}
