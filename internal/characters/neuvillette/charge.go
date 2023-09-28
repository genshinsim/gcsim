package neuvillette

import (
	"sort"

	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// Hopefully chargeFrames of smol CA will end up equal to the endLag of the big CA?
var chargeFrames []int
var endLag []int

const shortChargeHitmark = 55
const startCharge = 30
const chargeJudgementName = "Charged Attack: Equitable Judgment"

func init() {
	chargeFrames = frames.InitAbilSlice(60)
	endLag = frames.InitAbilSlice(30)
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	// skip CA windup if we're in NA/CA animation
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState, action.ChargeAttackState:
		windup = 14
	}

	playerPos := c.Core.Combat.Player().Pos()

	// TODO: Find the base duration
	chargeLegalEvalLeft := 240
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

	for _, g := range droplets {
		g.Kill()
		c.healWithDroplets()
		chargeLegalEvalLeft -= 90
		if chargeLegalEvalLeft <= 0 {
			chargeLegalEvalLeft = 0
			break
		}
	}
	if p["short"] != 0 {
		// By releasing too fast it is possible to absorb 3 orbs but not do a big CA
		// Apparently this is the same input as doing a fast CA cancel so it might be random
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

			// TODO: Not sure of snapshot timing
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 5}, 3),
				shortChargeHitmark-windup,
				shortChargeHitmark-windup,
			)
		}

		return action.Info{
			Frames:          func(next action.Action) int { return startCharge - windup + chargeFrames[next] },
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

	chargeJudgementStart := startCharge - windup + chargeLegalEvalLeft
	c.chargeJudgementDur = 180

	if c.Base.Cons >= 6 {
		// the c6 droplet check has to happen immediately because otherwise we don't know how long this action will take
		c.c6DropletCheck()
	}

	// TODO: param for letting the user not do the full channel?
	c.QueueCharTask(c.consumeHp, chargeJudgementStart+0)
	c.QueueCharTask(c.judgementWave, chargeJudgementStart+0)

	return action.Info{
		Frames: func(next action.Action) int {
			return startCharge - windup + chargeLegalEvalLeft + c.chargeJudgementDur + endLag[next]
		},
		AnimationLength: startCharge - windup + chargeLegalEvalLeft + c.chargeJudgementDur + endLag[action.InvalidAction],
		CanQueueAfter:   startCharge - windup + chargeLegalEvalLeft + c.chargeJudgementDur + endLag[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) judgementWave() {
	if c.Core.F-c.chargeSrc >= c.chargeJudgementDur {
		return
	}
	// Maybe we can optimize the Attack Pattern to not be recalculated every hit
	// since sim changing position and/or primary target during the CA is not supported?
	ap := combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{}, 3, 8)
	if c.Base.Ascension >= 1 {
		c.chargeAi.FlatDmg = chargeJudgement[c.TalentLvlAttack()] * c.MaxHP() * a1Multipliers[c.countA1()]
	}
	c.Core.QueueAttack(c.chargeAi, ap, 0, 0, c.c6cb)

	// Seems to be 0.42s interval?
	c.QueueCharTask(c.judgementWave, 25)
}

func (c *char) consumeHp() {
	// He only drains 5 times in 3s, on frame 0, 30, 60, 90, 120, 150?
	if c.Core.F-c.chargeSrc >= c.chargeJudgementDur {
		return
	}

	if c.CurrentHPRatio() <= 0.5 {
		return
	}

	hpDrain := 0.08 * c.MaxHP()

	c.Core.Player.Drain(player.DrainInfo{
		ActorIndex: c.Index,
		Abil:       "Charged Attack: Equitable Judgment",
		Amount:     hpDrain,
	})
	c.QueueCharTask(c.consumeHp, 30)
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
