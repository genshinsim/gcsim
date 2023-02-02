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
	burstHitmarks = 16 // initial 3 hits
	burstRadius   = []float64{2.5, 2.5, 3}
	burstDur      = 5 * 60
)

func init() {
	burstFrames = frames.InitAbilSlice(80)
	burstFrames[action.ActionSwap] = 79
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	//add cooldown to sim
	c.SetCDWithDelay(action.ActionBurst, 20*60, 18)
	//use up energy
	c.ConsumeEnergy(7)

	burstAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Moonjade Descent",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstDMG[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(burstAI, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 5), 16, 16)
	c.burstRadishAI = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Radish (Burst)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupYaoyaoRadishBurst,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstRadishDMG[c.TalentLvlBurst()],
	}
	c.QueueCharTask(c.newYueguiJump, 1*60+42)
	c.QueueCharTask(c.newYueguiJump, 2*60+42)
	c.QueueCharTask(c.newYueguiJump, 3*60+42)
	c.QueueCharTask(c.removeBurst, burstDur)
	c.AddStatus(burstKey, burstDur, true)

	// TODO: Yaoyao gains 15% movespeed and 50% dendro res
	// m := make([]float64, attributes.EndStatType)
	// m[attributes.DendroRes] = 0.50
	// m[attributes.Movespeed] = 0.15
	// c.AddStatMod(character.StatMod{
	// 		Base:         modifier.NewBaseWithHitlag(burstKey, 600),
	// 		AffectedStat: attributes.DendroRes,
	// 		Amount: func() ([]float64, bool) {
	// 			return m, true
	// 		}
	// },)

	if c.Base.Cons >= 4 {
		c.c4()
	}

	if c.Base.Ascension >= 1 {
		for i := 36; i <= burstDur; i += 36 { // 0.6*60 = 36
			c.QueueCharTask(c.a1ticker, i)
		}
	}
	c.ConsumeEnergy(0)
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) getBurstHealInfo() player.HealInfo {
	heal := burstRadishHealing[0][c.TalentLvlBurst()]*c.MaxHP() + burstRadishHealing[1][c.TalentLvlBurst()]
	return player.HealInfo{
		Caller:  c.Index,
		Target:  c.Core.Player.Active(),
		Message: "Yuegui burst heal",
		Src:     heal,
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
