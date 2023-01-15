package yelan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int
var burstDiceHitmarks = []int{25, 30, 36, 41} //c2 hitmark not framecounted

// initial hit
const burstHitmark = 76

func init() {
	burstFrames = frames.InitAbilSlice(93) // Q -> N1/CA/D
	burstFrames[action.ActionSkill] = 92   // Q -> E
	burstFrames[action.ActionJump] = 91    // Q -> J
	burstFrames[action.ActionSwap] = 90    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Depth-Clarion Dice",
		AttackTag:        combat.AttackTagElementalBurst,
		ICDTag:           combat.ICDTagNone,
		ICDGroup:         combat.ICDGroupDefault,
		StrikeType:       combat.StrikeTypePierce,
		Element:          attributes.Hydro,
		Durability:       50,
		Mult:             0,
		FlatDmg:          burst[c.TalentLvlBurst()] * c.MaxHP(),
		HitlagHaltFrames: 0.05 * 60,
		HitlagFactor:     0.05,
		IsDeployable:     true,
	}
	//apply hydro every 3rd hit
	//triggered on normal attack or yelan's skill

	//Initial hit
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{X: -1.5, Y: -1.7}, 6),
		burstHitmark,
		burstHitmark,
	)

	//TODO: check if we need to add f to this
	c.Core.Tasks.Add(func() {
		c.AddStatus(burstKey, 15*60, false)
		c.a4() //TODO: does this call need to be delayed?
	}, burstHitmark)
	if c.Base.Cons >= 6 { //C6 passive, lasts 20 seconds
		c.Core.Status.Add(c6Status, 20*60)
		c.c6count = 0
	}
	c.Core.Log.NewEvent("burst activated", glog.LogCharacterEvent, c.Index).
		Write("expiry", c.Core.F+15*60)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(6)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) summonExquisiteThrow() {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Exquisite Throw",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagYelanBurst,
		ICDGroup:   combat.ICDGroupYelanBurst,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    burstDice[c.TalentLvlBurst()] * c.MaxHP(),
	}
	for i := 0; i < 3; i++ {
		//TODO: probably snapshots before hitmark
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				0.5,
			),
			burstDiceHitmarks[i],
			burstDiceHitmarks[i],
		)
	}
	if c.Base.Cons >= 2 && c.c2icd <= c.Core.F {
		ai.Abil = "Yelan C2 Proc"
		ai.FlatDmg = 14.0 / 100 * c.MaxHP()
		c.c2icd = c.Core.F + 1.8*60
		//TODO: frames timing on this?
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				0.5,
			),
			burstDiceHitmarks[3],
			burstDiceHitmarks[3],
		)
	}
}
