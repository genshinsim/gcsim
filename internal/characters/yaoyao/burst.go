package yaoyao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const burstKey = "yaoyaoburst"

var (
	burstFrames         []int
	burstInitialHitmark = 16
	burstDur            = 6 * 60
)

func init() {
	burstFrames = frames.InitAbilSlice(63)
	burstFrames[action.ActionAttack] = 58
	burstFrames[action.ActionSkill] = 57
	burstFrames[action.ActionDash] = 58
	burstFrames[action.ActionJump] = 57
	burstFrames[action.ActionSwap] = 56
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// add cooldown to sim
	c.SetCD(action.ActionBurst, 20*60)
	// use up energy
	c.ConsumeEnergy(7)

	burstAI := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Moonjade Descent",
		AttackTag:        attacks.AttackTagElementalBurst,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		Durability:       25,
		Mult:             burstDMG[c.TalentLvlBurst()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.05,
	}
	c.Core.QueueAttack(burstAI, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3), burstInitialHitmark, burstInitialHitmark)

	c.Core.Tasks.Add(c.newYueguiJump, 104)
	c.Core.Tasks.Add(c.newYueguiJump, 162)
	c.Core.Tasks.Add(c.newYueguiJump, 221)
	c.Core.Tasks.Add(c.removeBurst, burstDur)
	c.AddStatus(burstKey, burstDur, false)

	if c.Base.Cons >= 4 {
		c.c4()
	}

	if c.Base.Ascension >= 1 {
		c.QueueCharTask(c.a1Ticker, 0.6*60)
	}
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) getBurstHealInfo(snap *combat.Snapshot) player.HealInfo {
	maxhp := snap.BaseHP*(1+snap.Stats[attributes.HPP]) + snap.Stats[attributes.HP]
	heal := burstRadishHealing[0][c.TalentLvlBurst()]*maxhp + burstRadishHealing[1][c.TalentLvlBurst()]
	return player.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "Yuegui Burst Heal",
		Src:     heal,
		Bonus:   snap.Stats[attributes.Heal],
	}
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.StatusIsActive(burstKey) {
			c.removeBurst()
		}
		return false
	}, "yaoyao-exit")
}

func (c *char) removeBurst() {
	c.DeleteStatus(burstKey)
	// remove all jumping yuegui
	for i, yg := range c.yueguiJumping {
		if yg != nil {
			yg.Kill()
			c.yueguiJumping[i] = nil
		}
	}
	c.numYueguiJumping = 0
}
