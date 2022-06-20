package kazuha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var plungeHoldFrames []int
var plungePressFrames []int

const (
	plungeHoldAnimation  = 60
	plungePressAnimation = 55
)

func init() {
	plungeHoldFrames = frames.InitAbilSlice(50)
	plungeHoldFrames[action.ActionAttack] = 60
	plungeHoldFrames[action.ActionBurst] = 60

	plungePressFrames = frames.InitAbilSlice(47)
	plungeHoldFrames[action.ActionAttack] = 55
	plungeHoldFrames[action.ActionBurst] = 55
}

func (c *char) HighPlungeAttack(p map[string]int) action.ActionInfo {
	ele := attributes.Physical
	//TODO: this really shouldn't be anything else since it should only be used after skill?
	if c.Core.Player.LastAction.Char == c.Index && c.Core.Player.LastAction.Type == action.ActionSkill {
		ele = attributes.Anemo
	}

	a := action.ActionInfo{
		State: action.PlungeAttackState,
	}
	//TODO: is this accurate?? these should be the hitmarks
	var f int
	if c.Core.Player.LastAction.Param["hold"] > 0 {
		f = 41
		a.Frames = frames.NewAbilFunc(plungeHoldFrames)
		a.AnimationLength = plungeHoldAnimation
		a.CanQueueAfter = f
	} else {
		f = 36
		a.Frames = frames.NewAbilFunc(plungePressFrames)
		a.AnimationLength = plungePressAnimation
		a.CanQueueAfter = f
	}

	_, ok := p["collide"]
	if ok {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Plunge (Collide)",
			AttackTag:      combat.AttackTagPlunge,
			ICDTag:         combat.ICDTagNone,
			ICDGroup:       combat.ICDGroupDefault,
			Element:        ele,
			Durability:     0,
			Mult:           plunge[c.TalentLvlAttack()],
			IgnoreInfusion: true,
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.3, false, combat.TargettableEnemy), f, f)
	}

	//aoe dmg
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Plunge",
		AttackTag:      combat.AttackTagPlunge,
		ICDTag:         combat.ICDTagNone,
		ICDGroup:       combat.ICDGroupDefault,
		StrikeType:     combat.StrikeTypeBlunt,
		Element:        ele,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
		IgnoreInfusion: true,
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), f, f)

	// a1 if applies
	if c.a1Ele != attributes.NoElement {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Kazuha A1",
			AttackTag:      combat.AttackTagPlunge,
			ICDTag:         combat.ICDTagNone,
			ICDGroup:       combat.ICDGroupDefault,
			StrikeType:     combat.StrikeTypeDefault,
			Element:        c.a1Ele,
			Durability:     25,
			Mult:           2,
			IgnoreInfusion: true,
		}

		c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), f-1, f-1)
		c.a1Ele = attributes.NoElement
	}

	return a
}
