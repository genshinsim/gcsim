package hutao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var chargeFrames []int
var paramitaChargeFrames []int

const chargeHitmark = 19
const paramitaChargeHitmark = 6

func init() {
	// charge -> x
	chargeFrames = frames.InitAbilSlice(62)
	chargeFrames[action.ActionAttack] = 57
	chargeFrames[action.ActionSkill] = 57
	chargeFrames[action.ActionSkill] = 60
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark

	// charge (paramita) -> x
	paramitaChargeFrames = frames.InitAbilSlice(44)
	paramitaChargeFrames[action.ActionBurst] = 35
	paramitaChargeFrames[action.ActionDash] = paramitaChargeHitmark
	paramitaChargeFrames[action.ActionJump] = paramitaChargeHitmark
	paramitaChargeFrames[action.ActionSwap] = 42
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {

	if c.StatModIsActive(paramitaBuff) {
		return c.ppChargeAttack(p)
	}

	//check for particles
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupPole,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy), 0, chargeHitmark, c.ppParticles, c.applyBB)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}

func (c *char) ppChargeAttack(p map[string]int) action.ActionInfo {
	//TODO: currently assuming snapshot is on cast since it's a bullet and nothing implemented re "pp slide"
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupPole,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy), 0, paramitaChargeHitmark, c.ppParticles, c.applyBB)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(paramitaChargeFrames),
		AnimationLength: paramitaChargeFrames[action.InvalidAction],
		CanQueueAfter:   paramitaChargeHitmark,
		State:           action.ChargeAttackState,
	}
}

func (c *char) applyBB(a combat.AttackCB) {
	if !c.StatModIsActive(paramitaBuff) {
		return
	}
	trg, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	if !trg.StatusIsActive(bbDebuff) {
		//start ticks
		trg.QueueEnemyTask(c.bbtickfunc(c.Core.F, trg), 240)
		trg.SetTag(bbDebuff, c.Core.F) //to track current bb source
	}

	trg.AddStatus(bbDebuff, 570, true) //lasts 8s + 1.5s
}

func (c *char) bbtickfunc(src int, trg *enemy.Enemy) func() {
	return func() {
		//do nothing if source changed
		if trg.Tags[bbDebuff] != src {
			return
		}
		if !trg.StatusIsActive(bbDebuff) {
			return
		}
		c.Core.Log.NewEvent("Blood Blossom checking for tick", glog.LogCharacterEvent, c.Index).
			Write("src", src)

		//queue up one damage instance
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Blood Blossom",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       bb[c.TalentLvlSkill()],
		}
		//if cons 2, add flat dmg
		if c.Base.Cons >= 2 {
			ai.FlatDmg += c.MaxHP() * 0.1
		}
		c.Core.QueueAttack(ai, combat.NewDefSingleTarget(trg.Index(), combat.TargettableEnemy), 0, 0)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Blood Blossom ticked", glog.LogCharacterEvent, c.Index).
				Write("next expected tick", c.Core.F+240).
				Write("dur", trg.StatusExpiry(bbDebuff)).
				Write("src", src)
		}
		//queue up next instance
		c.Core.Tasks.Add(c.bbtickfunc(src, trg), 240)
	}
}
