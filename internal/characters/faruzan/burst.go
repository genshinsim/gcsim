package faruzan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
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
