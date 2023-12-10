package neuvillette

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
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

func (c *char) legalEvalFindDroplets() int {
	droplets := c.getSourcewaterDroplets()

	// TODO: If droplets time out before the "droplet check" it doesn't count.
	indices := c.Core.Combat.Rand.Perm(len(droplets))
	orbs := 0
	for _, ind := range indices {
		g := droplets[ind]
		c.consumeDroplet(g)
		orbs += 1
		if orbs >= 3 {
			break
		}
	}
	c.Core.Combat.Log.NewEvent(fmt.Sprint("Picked up ", orbs, " droplets"), glog.LogCharacterEvent, c.Index)
	return orbs
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

	return c.chargeAttackJudgement(p, windup)
}

func (c *char) chargeAttackJudgement(p map[string]int, windup int) (action.Info, error) {
	c.chargeJudgeDur = 0
	c.tickAnimLength = getChargeJudgementHitmarkDelay(0)
	// current framework doesn't really support actions getting shorter, so the legal eval is set to 0, but it may increase later
	chargeLegalEvalLeft := 0

	c.QueueCharTask(func() {
		chargeLegalEvalLeft = initialLegalEvalDur
		orbs := c.legalEvalFindDroplets()
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

		ticks, ok := p["ticks"]
		if !ok {
			ticks = -1
		} else if ticks < 0 {
			ticks = 0
		}

		// cannot use hitlag affected queue because the logic just does not work then
		// -> can't account for possible hitlag delaying the update of the anim length (sim moves on to next action, but ticks continue)

		// start counting at 1 for correct number of ticks when supplying ticks param
		c.Core.Tasks.Add(c.chargeJudgementTick(c.chargeJudgeStartF, 1, ticks, false), chargeLegalEvalLeft+getChargeJudgementHitmarkDelay(1))

		// He drains 5 times in 3s, on frame 40, 70, 100, 130, 160
		c.QueueCharTask(c.consumeHp(c.chargeJudgeStartF), chargeLegalEvalLeft+40)
	}, windup+3)

	return action.Info{
		Frames: func(next action.Action) int {
			return windup + 3 + chargeLegalEvalLeft + c.tickAnimLength + endLag[next]
		},
		AnimationLength: 1200, // there is no upper limit on the duration of the CA
		CanQueueAfter:   windup + 3 + endLag[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) chargeAttackShort(windup int) (action.Info, error) {
	// By releasing too fast it is possible to absorb 3 orbs but not do a big CA
	c.QueueCharTask(func() {
		c.legalEvalFindDroplets()
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
			ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), 3, 8)
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
	// calculated every hit since canqueueafter is after the first tick, so configs can change the primary target/entity positions while the CA happens
	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), 3.5, 15)
	if c.Base.Ascension >= 1 {
		c.chargeAi.FlatDmg = chargeJudgement[c.TalentLvlAttack()] * c.MaxHP() * a1Multipliers[c.countA1()]
	}
	c.Core.QueueAttack(c.chargeAi, ap, 0, 0)
}

func getChargeJudgementHitmarkDelay(tick int) int {
	// first tick happens 6f after start, second tick is 22f after first, then other frames are 25f after, then last tick is when the judgement wave ends.
	// TODO: check this is the case for c6
	switch tick {
	case 1:
		return 6
	case 2:
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
		// no longer in CA anim -> no tick
		if c.Core.F > c.chargeJudgeStartF+c.chargeJudgeDur {
			return
		}

		// last tick -> check for C6 extension
		if last {
			// C6 did not extend CA -> proc wave and stop queuing ticks
			if c.Core.F == c.chargeJudgeStartF+c.chargeJudgeDur {
				c.judgementWave()
			} else {
				// C6 extended the CA between when this tick was queued and when this tick was executed
				// -> allow the other non last queued task to execute by extend anim length to include that task
				c.tickAnimLength = c.tickAnimLengthC6Extend
			}
			return
		}

		// tick param supplied and hit the limit -> proc wave, enable early cancel flag for next action check and stop queuing ticks
		if tick == maxTick {
			c.judgementWave()
			c.chargeEarlyCancelled = true
			return
		}

		c.judgementWave()

		// next tick handling
		if maxTick == -1 || tick < maxTick {
			tickDelay := getChargeJudgementHitmarkDelay(tick + 1)
			// calc new animation length to be up until next tick happens
			nextTickAnimLength := c.Core.F - c.chargeJudgeStartF + tickDelay

			// always queue up non-final tick that will be executed in case C6 was proc'd
			c.Core.Tasks.Add(c.chargeJudgementTick(src, tick+1, maxTick, false), tickDelay)

			// queue up last tick if next tick would happen after CA duration ends
			if nextTickAnimLength > c.chargeJudgeDur {
				// queue up final tick to happen at end of CA duration
				c.Core.Tasks.Add(c.chargeJudgementTick(src, tick+1, maxTick, true), c.chargeJudgeDur-c.tickAnimLength)
				// update tickAnimLength to be equal to entire CA duration at the end
				c.tickAnimLength = c.chargeJudgeDur
				// if C6 is triggered, then tickAnimLength will be wrong so this var holds the actual tickAnimLength if ticks continued normally beyond the original final tick
				c.tickAnimLengthC6Extend = nextTickAnimLength
			} else {
				// next tick happens within CA duration -> update tickAnimLength as usual
				c.tickAnimLength = nextTickAnimLength
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

func (c *char) consumeDroplet(g *common.SourcewaterDroplet) {
	g.Kill()
	// the healing is slightly delayed by 8f
	c.QueueCharTask(func() {
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "Sourcewater Droplets Healing",
			Src:     c.MaxHP() * 0.16,
			Bonus:   c.Stat(attributes.Heal),
		})
	}, 8)
}
