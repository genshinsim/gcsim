package albedo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstHitmark = 96

func init() {
	burstFrames = frames.InitAbilSlice(96)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	hits, ok := p["bloom"]
	if !ok {
		hits = 2 //default 2 hits
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rite of Progeniture: Tectonic Tide",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       burst[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	//check stacks
	if c.Base.Cons >= 2 && c.Core.Status.Duration("albedoc2") > 0 {
		ai.FlatDmg += (snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]) * float64(c.Tags["c2"])
		c.Tags["c2"] = 0
	}

	//TODO: damage frame
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(3, false, combat.TargettableEnemy), burstHitmark)

	// Blooms are generated on a slight delay from initial hit
	// TODO: No precise frame data, guessing correct delay
	ai.Abil = "Rite of Progeniture: Tectonic Tide (Blossom)"
	ai.Mult = burstPerBloom[c.TalentLvlSkill()]
	for i := 0; i < hits; i++ {
		c.Core.QueueAttackWithSnap(ai, c.bloomSnapshot, combat.NewDefCircHit(3, false, combat.TargettableEnemy), burstHitmark+30+i*5)
	}

	//Party wide EM buff
	// a4: burst increase party em by 125 for 10s
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 125
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{Base: modifier.NewBase("albedo-a4", 600), AffectedStat: attributes.EM, Amount: func() ([]float64, bool) {
			return m, true
		}})
	}

	c.SetCDWithDelay(action.ActionBurst, 720, 80)
	c.ConsumeEnergy(80)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}
