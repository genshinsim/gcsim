package yelan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int
var burstTravel = 20

// initial hit
const burstHitmark = 76
const c2Hitmark = 17

func init() {
	burstFrames = frames.InitAbilSlice(93) // Q -> N1/CA/D
	burstFrames[action.ActionSkill] = 92   // Q -> E
	burstFrames[action.ActionJump] = 91    // Q -> J
	burstFrames[action.ActionSwap] = 90    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Depth-Clarion Dice",
		AttackTag:        attacks.AttackTagElementalBurst,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypePierce,
		Element:          attributes.Hydro,
		Durability:       50,
		Mult:             0,
		FlatDmg:          burst[c.TalentLvlBurst()] * c.MaxHP(),
		HitlagHaltFrames: 0.05 * 60,
		HitlagFactor:     0.05,
		IsDeployable:     true,
	}

	travel, ok := p["travel"]
	if ok {
		burstTravel = travel
	}

	// apply hydro every 3rd hit
	// triggered on normal attack or yelan's skill

	// Initial hit
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{X: -1.5, Y: -1.7}, 6),
		burstHitmark,
		burstHitmark,
	)

	//TODO: check if we need to add f to this
	c.Core.Tasks.Add(func() {
		c.AddStatus(burstKey, 15*60, false)
		c.a4() //TODO: does this call need to be delayed?
	}, burstHitmark)
	if c.Base.Cons >= 6 { // C6 passive, lasts 20 seconds
		c.Core.Status.Add(c6Status, 20*60)
		c.c6count = 0
	}
	c.Core.Log.NewEvent("burst activated", glog.LogCharacterEvent, c.Index).
		Write("expiry", c.Core.F+15*60)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(6)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstWaveWrapper() {
	c.summonExquisiteThrow()
	c.AddStatus(burstICDKey, 60, true)
}

func (c *char) summonExquisiteThrow() {
	hp := c.MaxHP()
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Exquisite Throw",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagYelanBurst,
		ICDGroup:   attacks.ICDGroupYelanBurst,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    burstDice[c.TalentLvlBurst()] * hp,
	}
	snap := c.Snapshot(&ai)
	for i := 0; i < 3; i++ {
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				0.5,
			),
			burstTravel+i*6,
		)
	}
	if c.Base.Cons >= 2 && c.c2icd <= c.Core.F {
		ai.Abil = "Yelan C2 Proc"
		ai.ICDTag = attacks.ICDTagNone
		ai.ICDGroup = attacks.ICDGroupDefault
		ai.FlatDmg = 14.0 / 100 * hp
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
			0,
			c2Hitmark,
		)
	}
}
