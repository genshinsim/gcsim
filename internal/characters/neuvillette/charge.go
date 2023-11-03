package neuvillette

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var chargeFrames []int
var endLag []int
var earlyCancelEndLag []int

const initialLegalEvalDur = 209

var dropletLegalEvalReduction = []int{0, 57, 57 + 54, 57 + 54 + 98}

const shortChargeHitmark = 55

const chargeJudgementName = "Charged Attack: Equitable Judgment"

func init() {
	chargeFrames = frames.InitAbilSlice(87)
	chargeFrames[action.ActionCharge] = 69
	chargeFrames[action.ActionSkill] = 26
	chargeFrames[action.ActionBurst] = 27
	chargeFrames[action.ActionDash] = 25
	chargeFrames[action.ActionJump] = 26
	chargeFrames[action.ActionWalk] = 61
	chargeFrames[action.ActionSwap] = 58

	endLag = frames.InitAbilSlice(51)
	endLag[action.ActionWalk] = 36
	endLag[action.ActionCharge] = 30
	endLag[action.ActionSwap] = 27
	endLag[action.ActionBurst] = 0
	endLag[action.ActionSkill] = 0
	endLag[action.ActionDash] = 0
	endLag[action.ActionJump] = 0

	earlyCancelEndLag = frames.InitAbilSlice(5000)
	earlyCancelEndLag[action.ActionBurst] = 0
	earlyCancelEndLag[action.ActionSkill] = 0
	earlyCancelEndLag[action.ActionDash] = 0
	earlyCancelEndLag[action.ActionJump] = 0
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.chargeEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Charged Attack: Equitable Judgement with Charged Attack", c.CharWrapper.Base.Key)
	}
	// there is a windup out of dash/jump/walk/swap. Otherwise it is rolled into the Q/E/CA/NA -> CA frames
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.Idle, action.DashState, action.JumpState, action.WalkState, action.SwapState:
		windup = 14
	}

	if p["short"] != 0 {
		return c.chargeAttackShort(windup)
	}

	return c.chargeAttackJudegement(p, windup)
}

func (c *char) chargeAttackJudegement(p map[string]int, windup int) (action.Info, error) {
	c.chargeJudgeDur = 0
	c.nextTickTime = getChargeJudgementHitmarkDelay(0)
	// current framework doesn't really support actions getting shorter, so the charge left is set to 0, but it will be increase later, when the ball absorb thing happens
	chargeLegalEvalLeft := 0

	c.QueueCharTask(func() {
		chargeLegalEvalLeft = initialLegalEvalDur
		droplets := c.getSourcewaterDroplets()

		// TODO: If droplets time out before the "droplet check" it doesn't count.
		indices := c.Core.Combat.Rand.Perm(len(droplets))

		// this should happen 3 frames into his CA
		orbs := 0
		for _, ind := range indices {
			g := droplets[ind]
			g.Kill()
			// the healing seems to be slightly delayed by 7f
			c.QueueCharTask(c.healWithDroplets, 10)
			orbs += 1
			if orbs >= 3 {
				break
			}
		}
		c.Core.Combat.Log.NewEvent(fmt.Sprint("Picked up ", orbs, " droplets"), glog.LogCharacterEvent, c.Index)
		chargeLegalEvalLeft -= dropletLegalEvalReduction[orbs]

		c.chargeAi = combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       chargeJudgementName,
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagExtraAttack,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    chargeJudgement[c.TalentLvlAttack()] * c.MaxHP(),
		}

		c.chargeJudgeStartF = c.Core.F + chargeLegalEvalLeft
		c.chargeJudgeDur = 173

		if c.Base.Cons >= 6 {
			c.QueueCharTask(c.c6DropletCheck(c.chargeJudgeStartF), chargeLegalEvalLeft)
			c.QueueCharTask(c.c6(c.chargeJudgeStartF), chargeLegalEvalLeft)
		}
		if p["ticks"] > 0 {
			c.QueueCharTask(c.chargeJudgementTick(c.chargeJudgeStartF, 0, p["ticks"], false), chargeLegalEvalLeft+getChargeJudgementHitmarkDelay(0))
		} else {
			c.QueueCharTask(c.chargeJudgementTick(c.chargeJudgeStartF, 0, -1, false), chargeLegalEvalLeft+getChargeJudgementHitmarkDelay(0))
		}
		// He drains 5 times in 3s, on frame 40, 70, 100, 130, 160
		c.QueueCharTask(c.consumeHp(c.chargeJudgeStartF), chargeLegalEvalLeft+30)
	}, windup+3)

	return action.Info{
		Frames: func(next action.Action) int {
			return windup + 3 + chargeLegalEvalLeft + c.nextTickTime + endLag[next]
		},
		AnimationLength: 1200, // there is no upper limit on the duration of the CA. We don't know when it ends until it ends
		CanQueueAfter:   windup + 3 + endLag[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) chargeAttackShort(windup int) (action.Info, error) {
	// By releasing too fast it is possible to absorb 3 orbs but not do a big CA
	c.QueueCharTask(func() {
		droplets := c.getSourcewaterDroplets()
		indices := c.Core.Combat.Rand.Perm(len(droplets))
		orbs := 0
		for _, ind := range indices {
			g := droplets[ind]
			g.Kill()
			// the healing seems to be slightly delayed by 7f
			c.QueueCharTask(c.healWithDroplets, 7)
			orbs += 1
			if orbs >= 3 {
				break
			}
		}
		c.Core.Combat.Log.NewEvent(fmt.Sprint("Picked up ", orbs, " droplets"), glog.LogCharacterEvent, c.Index)
		// If there is not enough stamina to CA, nothing happens and he floats back down
		r := 1 + c.Core.Player.StamPercentMod(action.ActionCharge)
		if r < 0 {
			r = 0
		}
		if c.Core.Player.Stam > 50*r {
			// use stam
			c.Core.Player.Stam -= 50 * r
			c.Core.Player.LastStamUse = c.Core.F
			c.Core.Player.Events.Emit(event.OnStamUse, action.ActionCharge)
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Charge Attack",
				AttackTag:  attacks.AttackTagExtra,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Hydro,
				Durability: 25,
				Mult:       charge[c.TalentLvlAttack()],
			}
			ap := combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{}, 3, 8)
			// TODO: Not sure of snapshot timing
			c.Core.QueueAttack(
				ai,
				ap,
				shortChargeHitmark+windup,
				shortChargeHitmark+windup,
			)
		}
	}, windup+3)

	return action.Info{
		Frames:          func(next action.Action) int { return windup + chargeFrames[next] },
		AnimationLength: windup + chargeFrames[action.InvalidAction],
		CanQueueAfter:   windup + chargeFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) judgementWave() {
	// Maybe we can optimize the Attack Pattern to not be recalculated every hit
	// since sim changing position and/or primary target during the CA is not supported?
	ap := combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{}, 3.5, 15)
	if c.Base.Ascension >= 1 {
		c.chargeAi.FlatDmg = chargeJudgement[c.TalentLvlAttack()] * c.MaxHP() * a1Multipliers[c.countA1()]
	}
	c.Core.QueueAttack(c.chargeAi, ap, 0, 0)
}

func getChargeJudgementHitmarkDelay(tick int) int {
	// first tick happens 6f after start, second tick is 22f after first, then other frames are 25f after, then last tick is when the judgement wave ends.
	// TODO: check this is the case for c6
	switch tick {
	case 0:
		return 6
	case 1:
		return 22
	default:
		return 25
	}
}

func (c *char) chargeJudgementTick(src, tick, maxTick int, last bool) func() {
	return func() {
		if c.chargeJudgeStartF != src {
			return
		}
		if last {
			if c.Core.F == c.chargeJudgeStartF+c.chargeJudgeDur {
				c.judgementWave()
			} else {
				c.nextTickTime = c.nextTickTime2
			}
			return
		}
		if c.Core.F > c.chargeJudgeStartF+c.chargeJudgeDur {
			return
		}
		c.judgementWave()
		if maxTick == -1 || tick < maxTick {
			nextF := c.Core.F + getChargeJudgementHitmarkDelay(tick+1)

			// this one is queued even if the above queues in case the c6 extension extends the charge attack between this tick and the final tick
			c.QueueCharTask(c.chargeJudgementTick(src, tick+1, maxTick, false), nextF-c.Core.F)
			c.nextTickTime = nextF - c.chargeJudgeStartF
			if nextF > c.chargeJudgeStartF+c.chargeJudgeDur {
				c.QueueCharTask(c.chargeJudgementTick(src, tick+1, maxTick, true), c.chargeJudgeStartF+c.chargeJudgeDur-c.Core.F)
				c.nextTickTime = c.chargeJudgeDur
				c.nextTickTime2 = nextF - c.chargeJudgeStartF
			}
		}
	}
}

func (c *char) consumeHp(src int) func() {
	return func() {
		if c.chargeJudgeStartF != src {
			return
		}
		if c.Core.F > c.chargeJudgeStartF+c.chargeJudgeDur {
			return
		}
		if c.CurrentHPRatio() > 0.5 {
			hpDrain := 0.08 * c.MaxHP()

			c.Core.Player.Drain(player.DrainInfo{
				ActorIndex: c.Index,
				Abil:       "Charged Attack: Equitable Judgment",
				Amount:     hpDrain,
			})
		}
		c.QueueCharTask(c.consumeHp(src), 30)
	}
}

func (c *char) healWithDroplets() {
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: "Sourcewater Droplets Healing",
		Src:     c.MaxHP() * 0.16,
		Bonus:   c.Stat(attributes.Heal),
	})
}
