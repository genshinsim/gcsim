package chiori

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const (
	burstHitmark        = 100 //TODO: i made this up
	burstSnapshotTiming = 10  //TODO: i made this up
)

func init() {
	burstFrames = frames.InitAbilSlice(103) // TODO: i made this up
}

// Twin swords leave their sheaths as Chiori slices with the clean cuts
// of a master tailor, dealing AoE Geo DMG based on her ATK and DEF.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Twin Blades: In Flight",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone, //TODO: to check
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		PoiseDMG:   200,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burstAtkScaling[c.TalentLvlBurst()],
	}

	c.Core.Tasks.Add(func() {
		//TODO: snapshot timing?
		snap := c.Snapshot(&ai)
		// flat dmg for def scaling portion
		ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		ai.FlatDmg *= burstDefScaling[c.TalentLvlBurst()]
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2), burstHitmark-burstSnapshotTiming)
	}, burstSnapshotTiming)

	c.ConsumeEnergy(10)                //TODO: delay??
	c.SetCD(action.ActionBurst, 60*15) //TODO: delay??

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], //TOOD: i made this up
		State:           action.BurstState,
	}, nil
}
