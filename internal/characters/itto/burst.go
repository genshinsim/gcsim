package itto

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstDuration  = 660 + 90 + 45 // barely cover basic combo
	burstBuffKey   = "itto-q"
	burstAtkSpdKey = "itto-q-atkspd"
)

func init() {
	burstFrames = frames.InitAbilSlice(91) // Q -> N1/CA0/CA1/CAF/E
	burstFrames[action.ActionDash] = 84    // Q -> D
	burstFrames[action.ActionJump] = 85    // Q -> J
	burstFrames[action.ActionSwap] = 90    // Q -> Swap
}

// Adapted from Noelle
// Burst:
// Time to show 'em the might of the Arataki Gang! For a time, Itto lets out his inner Raging Oni King, wielding his Oni King's Kanabou in battle.
// This state has the following special properties:
// - Converts Itto's Normal, Charged, and Plunging Attacks to Geo DMG. This cannot be overridden.
// - Increases Itto's Normal Attack SPD. Also increases his ATK based on his DEF.
// - On hit, the 1st and 3rd strikes of his attack combo will each grant Arataki Itto 1 stack of Superlative Superstrength.
// - Decreases Itto's Elemental and Physical RES by 20%.
// The Raging Oni King state will be cleared when Itto leaves the field.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	// N1 pre-stack tech. If Itto did N1 -> Q, then add a stack before Q def to atk conversion
	// https://library.keqingmains.com/evidence/characters/geo/itto#itto-n1-burst-cancel-ss-stack
	if p["prestack"] != 0 && c.Core.Player.CurrentState() == action.NormalAttackState && c.NormalCounter == 1 {
		c.addStrStack("n1-burst-cancel", 1)
	}

	// Generate a "fake" snapshot in order to show a listing of the applied mods in the debug
	aiSnapshot := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Royal Descent: Behold, Itto the Evil! (Stat Snapshot)",
	}
	c.Snapshot(&aiSnapshot)
	burstDefSnapshot := c.Base.Def*(1+c.NonExtraStat(attributes.DEFP)) + c.NonExtraStat(attributes.DEF)
	mult := defconv[c.TalentLvlBurst()]

	// TODO: Confirm exact timing of buff
	// Q def to atk conversion
	mATK := make([]float64, attributes.EndStatType)
	mATK[attributes.ATK] = mult * burstDefSnapshot
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstBuffKey, burstDuration),
		AffectedStat: attributes.ATK,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			return mATK, true
		},
	})

	// Q atk speed buff
	mAtkSpd := make([]float64, attributes.EndStatType)
	mAtkSpd[attributes.AtkSpd] = 0.10
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstAtkSpdKey, burstDuration),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.CurrentState() != action.NormalAttackState {
				return nil, false
			}
			return mAtkSpd, true
		},
	})

	c.Core.Log.NewEvent("itto burst", glog.LogSnapshotEvent, c.Index).
		Write("total def", burstDefSnapshot).
		Write("atk added", mATK[attributes.ATK]).
		Write("mult", mult)

	if c.Base.Cons >= 1 {
		// TODO: add link to itto-c1-mechanics tcl entry later
		// this is before Q animation is done, so no need for char queue
		// 75 is a rough count for when Itto gains the 2 stacks from C1
		c.Core.Tasks.Add(c.c1, 75)
	}

	if c.Base.Cons >= 2 {
		// TODO: check C2 delay, but it doesn't really matter
		// should apply after cd/energy delay
		// this is before Q animation is done, so no need for char queue
		c.Core.Tasks.Add(c.c2, 9)
	}

	// apply when burst ends
	c.burstCastF = c.Core.F
	if c.Base.Cons >= 4 {
		c.applyC4 = true
		src := c.burstCastF
		c.QueueCharTask(func() {
			if src == c.burstCastF {
				c.c4()
			}
		}, burstDuration)
	}

	c.SetCD(action.ActionBurst, 1080) // 18s * 60
	c.ConsumeEnergy(1)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		c.savedNormalCounter = 0
		prev := args[0].(int)
		if prev == c.Index && c.StatModIsActive(burstBuffKey) {
			c.DeleteStatMod(burstBuffKey)
			c.DeleteStatMod(burstAtkSpdKey)
			c.c4()
		}
		return false
	}, "itto-exit")
}
