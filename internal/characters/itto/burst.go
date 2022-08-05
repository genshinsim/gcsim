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
	burstAnimation = 91
	burstBuffKey   = "itto-burst"
)

func init() {
	burstFrames = frames.InitAbilSlice(92)
	burstFrames[action.ActionCharge] = 91
	burstFrames[action.ActionSkill] = 91
	burstFrames[action.ActionDash] = 86
	burstFrames[action.ActionJump] = 84
	burstFrames[action.ActionSwap] = 90
}

// Adapted from Noelle
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// Generate a "fake" snapshot in order to show a listing of the applied mods in the debug
	aiSnapshot := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Royal Descent: Behold, Itto the Evil! (Stat Snapshot)",
	}
	snapshot := c.Snapshot(&aiSnapshot)
	burstDefSnapshot := snapshot.BaseDef*(1+snapshot.Stats[attributes.DEFP]) + snapshot.Stats[attributes.DEF]
	mult := defconv[c.TalentLvlBurst()]

	// TODO: Confirm exact timing of buff - for now matched to status duration previously set, which is 900 + animation frames
	// padded to cover basic combo
	burstDur := burstAnimation + 840 // ~15.5s.

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATK] = mult * burstDefSnapshot
	m[attributes.AtkSpd] = .10
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstBuffKey, burstDur),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.CurrentState() == action.NormalAttackState {
				m[attributes.AtkSpd] = .10
				return m, true
			}
			m[attributes.AtkSpd] = 0
			return m, true
		},
	})
	c.Core.Log.NewEvent("itto burst", glog.LogSnapshotEvent, c.Index).
		Write("total def", burstDefSnapshot).
		Write("atk added", m[attributes.ATK]).
		Write("mult", mult)

	if c.Base.Cons >= 1 {
		// TODO: add link to itto-c1-mechanics tcl entry later
		// this is before Q animation is done, so no need for char queue
		// 75 is a rough count for when Itto gains the 2 stacks from C1
		c.Core.Tasks.Add(c.c1, 75)
	}

	// apply when burst ends
	c.burstCastF = c.Core.F
	if c.Base.Cons >= 4 {
		c.applyC4 = true
		src := c.burstCastF
		c.Core.Tasks.Add(func() {
			if src == c.burstCastF {
				c.c4()
			}
		}, burstDur)
	}

	// handle energy and c2
	cd := 1080
	c.ConsumeEnergy(0)
	if c.Base.Cons >= 2 {
		cd -= c.geoCharCount * (1.5 * 60)
		c.AddEnergy("itto-c2", float64(c.geoCharCount)*6)
	}

	c.SetCD(action.ActionBurst, cd)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev == c.Index && c.StatModIsActive(burstBuffKey) {
			c.DeleteStatMod(burstBuffKey)
			c.c4()
		}
		return false
	}, "itto-exit")
}
