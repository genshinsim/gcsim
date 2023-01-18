package yaoyao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const burstKey = "yaoyaoburst"

var (
	burstFrames   []int
	burstHitmarks = []int{18, 33, 56} // initial 3 hits
	burstRadius   = []float64{2.5, 2.5, 3}
	skillStart    = 40
)

func init() {
	burstFrames = frames.InitAbilSlice(80)
	burstFrames[action.ActionSwap] = 79
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.AddStatus(burstKey, 5*60, true)
	//add cooldown to sim
	c.SetCDWithDelay(action.ActionBurst, 20*60, 18)
	//use up energy
	c.ConsumeEnergy(24)

	c.burstAI = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Radish (Burst)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupYaoyaoRadishBurst,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		// Mult:       burstDmg[c.TalentLvlBurst()],
	}
	c.QueueCharTask(c.newYueguiJump, 1*60+skillStart)
	c.QueueCharTask(c.newYueguiJump, 2*60+skillStart)
	c.QueueCharTask(c.newYueguiJump, 3*60+skillStart)
	c.QueueCharTask(c.removeBurst, 5*60+skillStart)
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) getBurstHealInfo() player.HealInfo {
	return player.HealInfo{
		Caller:  c.Index,
		Target:  c.Core.Player.Active(),
		Message: "Yuegui burst heal",
		// Src:     burstHeal[c.TalentLvlBurst()],
	}
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.StatModIsActive(burstKey) {
			c.removeBurst()
		}
		return false
	}, "yaoyao-exit")
}

func (c *char) removeBurst() {
	c.DeleteStatMod(burstKey)
	// remove all jumping yuegui
	for i, yg := range c.yueguiJumping {
		yg.Kill()
		c.yueguiJumping[i] = nil

	}
	c.numYueguiJumping = 0
}
