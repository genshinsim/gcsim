package sara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstStart = 47           // lines up with cd start
const burstInitialHitmark = 51  // Initial Hit
const burstClusterHitmark = 100 // First Cluster Hit

func init() {
	burstFrames = frames.InitAbilSlice(80) // Q -> CA
	burstFrames[action.ActionAttack] = 78  // Q -> N1
	burstFrames[action.ActionSkill] = 57   // Q -> E
	burstFrames[action.ActionDash] = 58    // Q -> D
	burstFrames[action.ActionJump] = 58    // Q -> J
	burstFrames[action.ActionSwap] = 56    // Q -> Swap
}

// Implements burst handling.
// Casts down Tengu Juurai: Titanbreaker, dealing AoE Electro DMG. Afterwards, Tengu Juurai: Titanbreaker spreads out into 4 consecutive bouts of Tengu Juurai: Stormcluster, dealing AoE Electro DMG.
// Tengu Juurai: Titanbreaker and Tengu Juurai: Stormcluster can provide the active character within their AoE with the same ATK Bonus as given by the Elemental Skill, Tengu Stormcall. The ATK Bonus provided by various kinds of Tengu Juurai will not stack, and their effects and duration will be determined by the last Tengu Juurai to take effect.
// Has parameters: "wave_cluster_hits", which controls how many of the mini-clusters in each wave hit an opponent.
// Also has "waveAttackProcs", used to determine which waves proc the attack buff.
// Format for both is a digit of length 5 - rightmost value is the starting proc (titanbreaker hit), and it moves from right to left
// For example, if you want waves 3 and 4 only to proc the attack buff, set waveAttackProcs=11000
// For "wave_cluster_hits", use numbers in each slot to control the # of hits. So for center hit, then 3 hits from each wave, set wave_cluster_hits=33331
// Default for both is for the main titanbreaker and 1 wave to hit and also proc the buff
// Also implements C4
// The number of Tengu Juurai: Stormcluster released by Subjugation: Koukou Sendou is increased to 6.
func (c *char) Burst(p map[string]int) action.ActionInfo {
	waveClusterHits, ok := p["wave_cluster_hits"]
	if !ok {
		waveClusterHits = 41
		if c.Base.Cons >= 2 {
			waveClusterHits = 61
		}
	}
	waveAttackProcs, ok := p["waveAttackProcs"]
	if !ok {
		waveAttackProcs = 11
	}

	// Entire burst snapshots sometime after activation but before 1st hit.
	// For now, assume that it snapshots on cd delay
	// Flagged as no ICD since the stormclusters do not share ICD with the main hit
	// No ICD should not functionally matter as this only hits once

	//titan breaker
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Tengu Juurai: Titanbreaker",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstMain[c.TalentLvlBurst()],
	}
	// dTitanbreaker.Targets = combat.TargetAll

	// dStormcluster.Targets = combat.TargetAll

	var c1cb combat.AttackCBFunc
	if c.Base.Cons >= 1 {
		c1cb = func(a combat.AttackCB) {
			if a.Target.Type() != combat.TargettableEnemy {
				return
			}
			c.c1()
		}
	}

	if waveClusterHits%10 == 1 {
		// TODO: proper positioning
		// Actual hit procs after the full cast duration, or 50 frames
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6),
			burstStart,
			burstInitialHitmark,
			c1cb,
		)
	}
	if waveAttackProcs%10 == 1 {
		c.attackBuff(burstInitialHitmark)
	}

	//stormcluster
	// Each cluster wave hits ~50 frames after titanbreaker and each preceding wave
	// TODO: Replace with frame counts from KQM when those are available
	ai.Abil = "Tengu Juurai: Stormcluster"
	ai.ICDTag = combat.ICDTagElementalBurst
	ai.Mult = burstCluster[c.TalentLvlBurst()]
	for waveN := 0; waveN < 4; waveN++ {
		// Handles the potential manual user override through the input tags
		// For each wave, get the corresponding digit from the numeric sequence (e.g. for 4441, wave 2 = 4)
		waveHits := int((waveClusterHits % PowInt(10, waveN+2)) / PowInt(10, waveN+2-1))
		waveAttackProc := int((waveAttackProcs % PowInt(10, waveN+2)) / PowInt(10, waveN+2-1))
		if waveHits > 0 {
			for j := 0; j < waveHits; j++ {
				// TODO: proper positioning
				c.Core.QueueAttack(
					ai,
					combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3),
					burstStart,
					burstClusterHitmark+18*waveN,
					c1cb,
				)
			}
		}
		if waveAttackProc == 1 {
			c.attackBuff(burstClusterHitmark + 18*waveN)
		}
	}

	c.SetCDWithDelay(action.ActionBurst, 20*60, burstStart)
	c.ConsumeEnergy(50)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

// Get integer power - required for burst
func PowInt(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}
