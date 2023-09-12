package noelle

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

// TODO: not sure about this
const (
	burstStart   = 80
	burstBuffKey = "noelle-burst"
)

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
	c.Snapshot(&aiSnapshot)
	burstDefSnapshot := c.Base.Def*(1+c.NonExtraStat(attributes.DEFP)) + c.NonExtraStat(attributes.DEF)
	mult := defconv[c.TalentLvlBurst()]
	if c.Base.Cons >= 6 {
		mult += 0.5
	}
	// Add mod for def to attack burst conversion
	c.burstBuff[attributes.ATK] = mult * burstDefSnapshot

	dur := 900 + burstStart // default duration
	if c.Base.Cons >= 6 {
		// https://library.keqingmains.com/evidence/characters/geo/noelle#noelle-c6-burst-extension
		// check extension
		ext, ok := p["extend"]
		if ok {
			if ext < 0 {
				ext = 0
			}
			if ext > 10 {
				ext = 10
			}
		} else {
			ext = 10 // to maintain prev default behaviour of full extension
		}

		dur += ext * 60
		c.Core.Log.NewEvent("noelle c6 extension applied", glog.LogCharacterEvent, c.Index).
			Write("total_dur", dur).
			Write("ext", ext)
	}
	// TODO: Confirm exact timing of buff - for now matched to status duration previously set, which is 900 + animation frames
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("noelle-burst", dur),
		AffectedStat: attributes.ATK,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			return c.burstBuff, true
		},
	})
	c.Core.Log.NewEvent("noelle burst", glog.LogSnapshotEvent, c.Index).
		Write("total def", burstDefSnapshot).
		Write("atk added", c.burstBuff[attributes.ATK]).
		Write("mult", mult)

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Sweeping Time (Burst)",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagElementalBurst,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Geo,
		Durability:         25,
		Mult:               burst[c.TalentLvlBurst()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.15 * 60,
		CanBeDefenseHalted: true,
	}

	// Burst part
	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.5),
			0,
			0,
			c.skillHealCB(),
		)
	}, 24)

	// Skill part
	// Burst and Skill part of Q have the same hitlag values and both can heal
	c.QueueCharTask(func() {
		ai.Abil = "Sweeping Time (Skill)"
		ai.Mult = burstskill[c.TalentLvlBurst()]
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4),
			0,
			0,
			c.skillHealCB(),
		)
	}, 65)

	c.SetCD(action.ActionBurst, 900)
	c.ConsumeEnergy(8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
