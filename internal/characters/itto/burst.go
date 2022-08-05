package itto

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(92) // Q -> N1/CA0/CA1/CAF
	burstFrames[action.ActionDash] = 86    // Q -> D
	burstFrames[action.ActionJump] = 84    // Q -> J
	burstFrames[action.ActionSwap] = 90    // Q -> Swap
}

const burstBuffKey = "itto-q"

// Adapted from Noelle
// Burst:
// Time to show 'em the might of the Arataki Gang! For a time, Itto lets out his inner Raging Oni King, wielding his Oni King's Kanabou in battle.
// This state has the following special properties:
// - Converts Itto's Normal, Charged, and Plunging Attacks to Geo DMG. This cannot be overridden.
// - Increases Itto's Normal Attack SPD. Also increases his ATK based on his DEF.
// - On hit, the 1st and 3rd strikes of his attack combo will each grant Arataki Itto 1 stack of Superlative Superstrength.
// - Decreases Itto's Elemental and Physical RES by 20%.
// The Raging Oni King state will be cleared when Itto leaves the field.
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// N1 pre-stack tech
	lastWasItto := c.Core.Player.LastAction.Char == c.Index
	lastAction := c.Core.Player.LastAction.Type
	if lastWasItto && lastAction == action.ActionAttack && c.NormalCounter == 1 {
		// If Itto did N1 -> Q, then add a stack before Q def to atk conversion
		c.changeStacks(1)
		c.Core.Log.NewEvent("itto n1 pre-stack added", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.Tags[c.stackKey])
	}

	// Add mod for def to attack burst conversion
	val := make([]float64, attributes.EndStatType)

	// Generate a "fake" snapshot in order to show a listing of the applied mods in the debug
	aiSnapshot := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Royal Descent: Behold, Itto the Evil! (Stat Snapshot)",
	}
	snapshot := c.Snapshot(&aiSnapshot)
	burstDefSnapshot := snapshot.BaseDef*(1+snapshot.Stats[attributes.DEFP]) + snapshot.Stats[attributes.DEF]
	mult := defconv[c.TalentLvlBurst()]
	fa := mult * burstDefSnapshot
	val[attributes.ATK] = fa

	// TODO: Confirm exact timing of buff
	// Q def to atk conversion
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c.burstBuffKey, c.burstBuffDuration),
		AffectedStat: attributes.ATK,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})

	// Q atk speed buff
	mAtkSpd := make([]float64, attributes.EndStatType)
	mAtkSpd[attributes.AtkSpd] = 0.10
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c.burstBuffKey+"-atkspd", c.burstBuffDuration),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.CurrentState() != action.NormalAttackState {
				return nil, false
			}
			return mAtkSpd, true
		},
	})

	c.Core.Log.NewEvent("itto burst", glog.LogSnapshotEvent, c.Index).
		Write("frame", c.Core.F).
		Write("total def", burstDefSnapshot).
		Write("atk added", fa).
		Write("mult", mult)

	if c.Base.Cons >= 1 {
		// TODO: add link to itto-c1-mechanics tcl entry later
		// this is before Q animation is done, so no need for char queue
		// 75 is a rough count for when Itto gains the 2 stacks from C1
		c.Core.Tasks.Add(c.c1(), 75)
	}

	if c.Base.Cons >= 2 {
		// TODO: check C2 delay, but it doesn't really matter
		// this is before Q animation is done, so no need for char queue
		c.Core.Tasks.Add(c.c2(), 9)
	}

	if c.Base.Cons >= 4 {
		c.c4Applied = false
		c.QueueCharTask(c.c4(), c.burstBuffDuration)
	}

	c.SetCD(action.ActionBurst, 660) // 11s * 60
	c.ConsumeEnergy(1)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}
}
