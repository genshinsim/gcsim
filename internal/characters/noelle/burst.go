package noelle

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

// TODO: not sure about this
const burstStart = 80

func init() {
	burstFrames = frames.InitAbilSlice(121)
	burstFrames[action.ActionAttack] = 83
	burstFrames[action.ActionSkill] = 82
	burstFrames[action.ActionDash] = 81
	burstFrames[action.ActionJump] = 81
	burstFrames[action.ActionWalk] = 90
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// TODO: Assume snapshot happens immediately upon cast since the conversion buffs the two burst hits
	// Generate a "fake" snapshot in order to show a listing of the applied mods in the debug
	aiSnapshot := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sweeping Time (Stat Snapshot)",
	}
	snapshot := c.Snapshot(&aiSnapshot)
	burstDefSnapshot := snapshot.BaseDef*(1+snapshot.Stats[attributes.DEFP]) + snapshot.Stats[attributes.DEF]
	mult := defconv[c.TalentLvlBurst()]
	if c.Base.Cons >= 6 {
		mult += 0.5
	}
	// Add mod for def to attack burst conversion
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATK] = mult * burstDefSnapshot

	// TODO: Confirm exact timing of buff - for now matched to status duration previously set, which is 900 + animation frames
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("noelle-burst", 900+burstStart),
		AffectedStat: attributes.ATK,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
	c.Core.Log.NewEvent("noelle burst", glog.LogSnapshotEvent, c.Index).
		Write("total def", burstDefSnapshot).
		Write("atk added", m[attributes.ATK]).
		Write("mult", mult)

	c.Core.Status.Add("noelleq", 900+burstStart)
	// Queue up task for Noelle burst extension
	// https://library.keqingmains.com/evidence/characters/geo/noelle#noelle-c6-burst-extension
	if c.Base.Cons >= 6 {
		c.Core.Tasks.Add(func() {
			if c.Core.Player.Active() == c.Index {
				return
			}
			// Adding the mod again with the same key replaces it
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBase("noelle-burst", 600),
				AffectedStat: attributes.ATK,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
			c.Core.Log.NewEvent("noelle max burst extension activated", glog.LogCharacterEvent, c.Index).
				Write("new_expiry", c.Core.F+600)
			c.Core.Status.Add("noelleq", 600)
		}, 900+burstStart)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sweeping Time (Burst)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(6.5, false, combat.TargettableEnemy), 24, 24)

	ai.Abil = "Sweeping Time (Skill)"
	ai.Mult = burstskill[c.TalentLvlBurst()]
	c.Core.QueueAttack(ai, combat.NewDefCircHit(4.5, false, combat.TargettableEnemy), 65, 65)

	c.SetCD(action.ActionBurst, 900)
	c.ConsumeEnergy(8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
