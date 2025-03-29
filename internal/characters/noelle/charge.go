package noelle

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const chargeHitNum = 1

var (
	chargeFrames          []int
	chargeHitmarks        = []int{23, 26}
	chargePoiseDMG        = []float64{60, 120}
	chargeHitlagHaltFrame = []float64{0.03, 0.15}
	chargeHitboxes        = [][]float64{{3, 3.5}, {5, 5.5}}
	chargeOffsets         = []float64{0.3, 0}
	chargeFinalDelay      = []int{39 - 26}
)

type ChargeState struct {
	Counter int
	StartF  int
	Ended   bool
	Error   bool
}

func init() {
	// dash/jump/swap cancels are in frames func
	// charge -> x
	chargeFrames = frames.InitAbilSlice(108) // CA
	chargeFrames[action.ActionAttack] = 106
	chargeFrames[action.ActionSkill] = 63
	chargeFrames[action.ActionBurst] = 63
	chargeFrames[action.ActionDash] = 0
	chargeFrames[action.ActionJump] = 0
	chargeFrames[action.ActionSwap] = 0
	chargeFrames[action.ActionWalk] = 103
}

func (c *char) windupFrames() int {
	windup := 56
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		if c.NormalCounter == 3 {
			windup -= 11
		}
	case action.ChargeAttackState:
		windup -= 15
		if c.caState.StartF != 0 {
			windup = 0
		}
	case action.PlungeAttackState:
		windup -= 9
	case action.DashState:
		windup -= 12
	}
	return windup
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.CurrentState() != action.ChargeAttackState {
		c.caState = ChargeState{}
	}
	if c.caState.Error {
		return action.Info{}, errors.New("use of delays between charge attack hits is invalid")
	}

	final := false
	hitIndex := c.caState.Counter
	windup := c.windupFrames()
	if c.caState.Ended || p["final"] != 0 {
		final = true
		hitIndex = chargeHitNum
	}
	hitmark := chargeHitmarks[hitIndex]

	var spinFrames, transFrames int // transition frames for spin -> final
	if c.caState.StartF != 0 && final {
		transFrames = chargeFinalDelay[(c.caState.Counter-1)%chargeHitNum]
	} else if !final {
		transFrames = chargeFinalDelay[hitIndex]
	}

	if final {
		hitmark += transFrames
	} else {
		spinFrames = hitmark
	}

	atkspd := c.Stat(attributes.AtkSpd)
	act := action.Info{
		Frames: func(next action.Action) int {
			f := chargeFrames[next]

			switch {
			// spin/final -> dash/jump/swap
			case f == 0:
				f = hitmark
			// spin -> spin/final
			case !final && next == action.ActionCharge:
				f = hitmark
			// spin -> final -> x, final -> x
			default:
				f += spinFrames + transFrames
			}

			return windup + frames.AtkSpdAdjust(f, atkspd)
		},
		OnRemoved: func(next action.AnimationState) {
			if next != action.ChargeAttackState {
				c.caState = ChargeState{}
			}
		},
		AnimationLength: windup + spinFrames + transFrames + chargeFrames[action.InvalidAction],
		CanQueueAfter:   windup + hitmark,
		State:           action.ChargeAttackState,
	}
	act.QueueAction(func() { c.caState.Error = !final }, act.CanQueueAfter+1) // hitmark+1

	if c.caState.StartF == 0 && !final {
		src := c.Core.F
		c.caState.StartF = src

		if p["no_limit"] == 0 {
			c.Core.Tasks.Add(func() {
				if c.caState.StartF == src {
					c.caState.Ended = true
				}
			}, 60*5)
		}

		act.QueueAction(func() {
			c.caStaminaTask(src, &c.caState.StartF, &c.caState.Ended)
		}, windup)
	}

	// FIXME: doing 0f attack from action queue results in applied(anim) hitlag being 1f short. wrap it in c.Core.Tasks.Add(func() { ... }, 0) for now
	act.QueueAction(func() {
		c.Core.Tasks.Add(func() { c.queueChargeAttack(hitIndex) }, 0)
	}, windup+hitmark)

	finalStart := windup + spinFrames + transFrames
	act.QueueAction(func() { c.caState = ChargeState{} }, finalStart)
	if !final {
		act.QueueAction(func() {
			c.Core.Tasks.Add(func() { c.queueChargeAttack(chargeHitNum) }, 0)
		}, finalStart+chargeHitmarks[chargeHitNum])
	}

	defer func() {
		c.caState.Counter = (c.caState.Counter + 1) % chargeHitNum
	}()

	return act, nil
}

func (c *char) queueChargeAttack(hitIndex int) {
	burstIndex := 0
	if c.StatModIsActive(burstBuffKey) {
		burstIndex = 1
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           chargePoiseDMG[hitIndex],
		Element:            attributes.Physical,
		Durability:         25,
		HitlagFactor:       0.01,
		HitlagHaltFrames:   chargeHitlagHaltFrame[hitIndex] * 60,
		CanBeDefenseHalted: true,
	}
	if hitIndex != chargeHitNum {
		ai.Abil = "Charge Attack"
		ai.Mult = charge[c.TalentLvlAttack()]
	} else {
		ai.Abil = "Charge Attack (Finisher)"
		ai.Mult = chargeFinal[c.TalentLvlAttack()]
	}
	if burstIndex != 0 {
		ai.ICDTag = attacks.ICDTagNone
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: chargeOffsets[hitIndex]},
		chargeHitboxes[burstIndex][hitIndex],
	)

	c.Core.QueueAttack(ai, ap, 0, 0, c.skillHealCB(), c.makeA4CB())
}

func (c *char) caStaminaTask(src int, startF *int, ended *bool) {
	const tickInterval = .3
	c.QueueCharTask(func() {
		if src != *startF || *ended {
			return
		}
		if c.Core.Player.Stam == 0 {
			*ended = true
			return
		}

		r := 1 + c.Core.Player.StamPercentMod(action.ActionCharge)
		if r < 0 {
			r = 0
		}
		c.Core.Player.UseStam((40*r)*tickInterval, action.ActionCharge)

		c.caStaminaTask(src, startF, ended)
	}, 60*tickInterval)
}
