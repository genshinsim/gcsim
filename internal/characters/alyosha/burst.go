package alyosha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstKey = "alyosha-burst"

	fieldFirstTick = 79
	fieldInterval  = 118
)

func init() {
	burstFrames = frames.InitAbilSlice(55) // Q -> D/J
	burstFrames[action.ActionAttack] = 54  // Q -> N1
	burstFrames[action.ActionSkill] = 53   // Q -> E
	burstFrames[action.ActionWalk] = 54    // Q -> W
	burstFrames[action.ActionSwap] = 52    // Q -> Swap
}

// Summons his trusted companion Tugarin to fight alongside him on the field. This also turns an
// area in front of him into a Fulgurite Hunting Field for a short time.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	src := c.Core.F
	c.burstSrc = src
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6)
	c.Core.Tasks.Add(func() { c.burstTicker(src, &ap) }, fieldFirstTick)
	c.AddStatus(burstKey, 14*60+c.c2BurstDur(), true)

	c.SetCDWithDelay(action.ActionBurst, 18*60, 1)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

// Fulgurite Hunting Field
// · Continuously taunts nearby opponents to incite them to attack.
// · Every 2s, deals an instance of AoE Electro DMG to any opponent that enters the field.
//
// Tugarin
// · If there are any opponents nearby, Tugarin will quickly move close to an opponent and maul them every 2s, dealing Electro DMG.
// · If an opponent affected by the Hunter's Mark effect is hit, the Hunter's Mark will also be activated.
// · Where there are multiple opponents, Tugarin will attack those affected by the Hunter's Mark effect first.
func (c *char) burstTicker(src int, ap *info.AttackPattern) {
	if c.burstSrc != src {
		return
	}

	if !c.StatusIsActive(burstKey) {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Fulgurite Hunting Field",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupAlyoshaBurst,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, *ap, 0, 0, c.c2MakeBurstCB())

	aiDog := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Tugarin",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupAlyoshaBurst,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}

	c.Core.Tasks.Add(func() {
		// TODO: should prioritize enemies with hunters mark, but currently just hit the primary target
		apDog := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1)
		c.Core.QueueAttack(aiDog, apDog, 0, 0, c.triggerSkillMarkCB(false))
		c.a1OnTugarin()
		c.c4OnTugarin()
	}, 39)

	c.Core.Tasks.Add(func() { c.burstTicker(src, ap) }, fieldInterval)
}
