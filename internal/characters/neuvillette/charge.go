package neuvillette

import (
	"fmt"
	"sort"

	"github.com/genshinsim/gcsim/internal/common"
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

// Hopefully chargeFrames of smol CA will end up equal to the endLag of the big CA?
var chargeFrames []int
var endLag []int
var earlyCancelEndLag []int

const initialLegalEvalDur = 212

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
	// there is a windup out of dash/jump/walk/swap. Otherwise it is rolled into the Q/E/CA/NA -> CA frames
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.Idle, action.DashState, action.JumpState, action.WalkState, action.SwapState:
		windup = 14
	}

	playerPos := c.Core.Combat.Player().Pos()

	chargeLegalEvalLeft := initialLegalEvalDur

	droplets := make([]*common.SourcewaterDroplet, 0)
	for _, g := range c.Core.Combat.Gadgets() {
		droplet, ok := g.(*common.SourcewaterDroplet)
		if !ok {
			continue
		}
		if droplet.Pos().Distance(playerPos) <= 8 {
			droplets = append(droplets, droplet)
		}
	}

	// TODO: If droplets time out before the "droplet check" it doesn't count.
	// However, this check needs to happen before c6 check, which needs to happen when this function is called.

	// TODO: Apparently it's semi random? I don't know how Neuv prioritizes his droplets
	sort.Slice(droplets, func(i, j int) bool {
		// return droplets[i].Pos().Distance(playerPos) < droplets[j].Pos().Distance(playerPos)
		return droplets[i].Duration < droplets[j].Duration
	})

	// TODO this should happen 3 frames into his CA but modifying action length
	// during the action is unsupported in the current framework
	orbs := 0
	for _, g := range droplets {
		g.Kill()
		// the healing seems to be slightly delayed by 10f
		c.QueueCharTask(c.healWithDroplets, 7)
		orbs += 1
		if orbs >= 3 {
			break
		}
	}
	c.Core.Combat.Log.NewEvent(fmt.Sprint("Picked up ", orbs, " droplets"), glog.LogCharacterEvent, c.Index)
	chargeLegalEvalLeft -= dropletLegalEvalReduction[orbs]

	if p["short"] != 0 {
		// By releasing too fast it is possible to absorb 3 orbs but not do a big CA
		r := 1 + c.Core.Player.StamPercentMod(action.ActionCharge)
		if r < 0 {
			r = 0
		}

		// If there is not enough stamina to CA, nothing happens and he floats back down
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
				shortChargeHitmark-windup,
				shortChargeHitmark-windup,
			)
		}

		return action.Info{
			Frames:          func(next action.Action) int { return windup + chargeFrames[next] },
			AnimationLength: chargeFrames[action.InvalidAction],
			CanQueueAfter:   chargeFrames[action.ActionDash],
			State:           action.ChargeAttackState,
		}, nil
	}

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

	chargeJudgementStart := windup + chargeLegalEvalLeft
	chargeJudgementDur := 173

	if c.Base.Cons >= 6 {
		// the c6 droplet check has to happen immediately because otherwise we don't know how long this action will take
		chargeJudgementDur += c.c6DropletCheck()
	}
	if p["ticks"] > 0 {
		// param for letting the user not do the full channel
		// calculate how long the judgement duration should be based on their tick count
		// additionally modify the frames so that only D/J/Q/E can follow. Otherwise sim errors on action
		// also need to verify it for c6
		maxTicks := p["ticks"]
		ticksDone := 0
		delay := getChargeJudgementHitmarkDelay(ticksDone)
		for delay < chargeJudgementDur && ticksDone < maxTicks {
			if ticksDone == maxTicks {
				for d := 40; d < delay; d += 30 {
					c.QueueCharTask(c.consumeHp, chargeJudgementStart+d)
				}
				return action.Info{
					Frames: func(next action.Action) int {
						return chargeJudgementStart + delay + earlyCancelEndLag[next]
					},
					AnimationLength: chargeJudgementStart + delay + earlyCancelEndLag[action.InvalidAction],
					CanQueueAfter:   chargeJudgementStart + delay + earlyCancelEndLag[action.ActionDash],
					State:           action.ChargeAttackState,
				}, nil
			}
			c.QueueCharTask(c.judgementWave, chargeJudgementStart+delay)

			ticksDone += 1
			delay = getChargeJudgementHitmarkDelay(ticksDone)
		}
		if maxTicks == ticksDone+1 {
			c.QueueCharTask(c.judgementWave, chargeJudgementStart+chargeJudgementDur)
			// He drains 5 times in 3s, on frame 40, 70, 100, 130, 160
			for d := 40; d < chargeJudgementDur; d += 30 {
				c.QueueCharTask(c.consumeHp, chargeJudgementStart+d)
			}
			return action.Info{
				Frames: func(next action.Action) int {
					return chargeJudgementStart + chargeJudgementDur + endLag[next]
				},
				AnimationLength: chargeJudgementStart + chargeJudgementDur + endLag[action.InvalidAction],
				CanQueueAfter:   chargeJudgementStart + chargeJudgementDur + endLag[action.ActionDash],
				State:           action.ChargeAttackState,
			}, nil
		}

		return action.Info{}, fmt.Errorf("%v: Cannot execute %d CA Judgement Ticks. Max executed %d", c.CharWrapper.Base.Key, maxTicks, ticksDone)
	}

	for tick := 0; getChargeJudgementHitmarkDelay(tick) < chargeJudgementDur; tick += 1 {
		c.QueueCharTask(c.judgementWave, chargeJudgementStart+getChargeJudgementHitmarkDelay(tick))
	}
	c.QueueCharTask(c.judgementWave, chargeJudgementStart+chargeJudgementDur)

	// He drains 5 times in 3s, on frame 40, 70, 100, 130, 160
	for d := 40; d < chargeJudgementDur; d += 30 {
		c.QueueCharTask(c.consumeHp, chargeJudgementStart+d)
	}

	return action.Info{
		Frames: func(next action.Action) int {
			return chargeJudgementStart + chargeJudgementDur + endLag[next]
		},
		AnimationLength: chargeJudgementStart + chargeJudgementDur + endLag[action.InvalidAction],
		CanQueueAfter:   chargeJudgementStart + chargeJudgementDur + endLag[action.ActionDash],
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
	if c.Base.Cons >= 6 {
		c.Core.QueueAttack(c.chargeAi, ap, 0, 0, c.c6cb)
		return
	}
	c.Core.QueueAttack(c.chargeAi, ap, 0, 0)
}

func getChargeJudgementHitmarkDelay(tick int) int {
	// first tick happens 6f after start, second tick is 22f after first, then other frames are 25f after, then last tick is when the judgement wave ends.
	// TODO: check this is the case for c6
	switch tick {
	case 0:
		return 6
	default:
		return 6 + 22 + 25*(tick-1)
	}
}

func (c *char) consumeHp() {
	if c.CurrentHPRatio() <= 0.5 {
		return
	}

	hpDrain := 0.08 * c.MaxHP()

	c.Core.Player.Drain(player.DrainInfo{
		ActorIndex: c.Index,
		Abil:       "Charged Attack: Equitable Judgment",
		Amount:     hpDrain,
	})
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
