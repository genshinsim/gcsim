package kazuha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var plungePressFrames []int
var plungeHoldFrames []int

const plungePressHitmark = 36
const plungeHoldHitmark = 41

// TODO: missing plunge -> skill
func init() {
	// skill (press) -> high plunge -> x
	plungePressFrames = frames.InitAbilSlice(55)
	plungePressFrames[action.ActionDash] = 48
	plungePressFrames[action.ActionJump] = 48
	plungePressFrames[action.ActionSwap] = 49

	// skill (hold) -> high plunge -> x
	plungeHoldFrames = frames.InitAbilSlice(61)
	plungeHoldFrames[action.ActionSkill] = 60 // uses burst frames
	plungeHoldFrames[action.ActionBurst] = 60
	plungeHoldFrames[action.ActionDash] = 48
	plungeHoldFrames[action.ActionJump] = 48
	plungeHoldFrames[action.ActionSwap] = 53
}

func (c *char) HighPlungeAttack(p map[string]int) action.ActionInfo {
	ele := attributes.Physical
	//TODO: this really shouldn't be anything else since it should only be used after skill?
	if c.Core.Player.LastAction.Char == c.Index && c.Core.Player.LastAction.Type == action.ActionSkill {
		ele = attributes.Anemo
	}

	act := action.ActionInfo{
		State: action.PlungeAttackState,
	}

	//TODO: is this accurate?? these should be the hitmarks
	var hitmark int
	if c.Core.Player.LastAction.Param["hold"] == 0 {
		hitmark = plungePressHitmark
		act.Frames = frames.NewAbilFunc(plungePressFrames)
		act.AnimationLength = plungePressFrames[action.InvalidAction]
		act.CanQueueAfter = plungePressFrames[action.ActionDash] // earliest cancel
	} else {
		hitmark = plungeHoldHitmark
		act.Frames = frames.NewAbilFunc(plungeHoldFrames)
		act.AnimationLength = plungeHoldFrames[action.InvalidAction]
		act.CanQueueAfter = plungeHoldFrames[action.ActionDash] // earliest cancel
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
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.3, false, combat.TargettableEnemy), hitmark, hitmark)
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

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), hitmark, hitmark)

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

		c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), hitmark-1, hitmark-1)
		c.a1Ele = attributes.NoElement
	}

	return act
}
