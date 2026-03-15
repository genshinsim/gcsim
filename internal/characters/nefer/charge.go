package nefer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	basicChargeWindup       = 20
	basicChargeHitmark      = 44
	phantasmConsumeDewFrame = 29
	phantasmHit1            = 30
	phantasmHit2            = 35
	phantasmHit3            = 43
	phantasmHit4            = 44
	phantasmHit5            = 45
)

var chargeFrames []int

func init() {
	chargeFrames = frames.InitAbilSlice(72)
	chargeFrames[action.ActionAttack] = 49
	chargeFrames[action.ActionSkill] = 48
	chargeFrames[action.ActionBurst] = 48
	chargeFrames[action.ActionDash] = 48
	chargeFrames[action.ActionJump] = 48
	chargeFrames[action.ActionSwap] = 48
	chargeFrames[action.ActionWalk] = 71
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a != action.ActionCharge {
		return c.Character.ActionStam(a, p)
	}
	if c.StatusIsActive(shadowDanceKey) && c.Core.Player.VerdantDew() > 0 {
		return 0
	}
	return 50
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(shadowDanceKey) && c.Core.Player.VerdantDew() > 0 {
		return c.phantasmChargeAttack()
	}
	return c.basicChargeAttack()
}

func (c *char) basicChargeAttack() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       charge[0][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), info.Point{Y: -2}, 3, 9),
		basicChargeHitmark,
		basicChargeHitmark,
	)
	c.QueueCharTask(c.absorbSeedsOfDeceit, basicChargeHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionAttack],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) phantasmChargeAttack() (action.Info, error) {
	consumeFrame := phantasmConsumeDewFrame
	c.QueueCharTask(func() {
		c.Core.Player.ConsumeVerdantDew(1)
		c.absorbSeedsOfDeceit()
	}, consumeFrame)

	neferHit1 := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Phantasm Performance (Nefer 1)",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       phantasm[0][c.TalentLvlSkill()],
		FlatDmg:    c.Stat(attributes.EM) * phantasm[1][c.TalentLvlSkill()],
	}
	neferHit2 := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Phantasm Performance (Nefer 2)",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       phantasm[2][c.TalentLvlSkill()],
		FlatDmg:    c.Stat(attributes.EM) * phantasm[3][c.TalentLvlSkill()],
	}
	shadeHit1 := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Phantasm Performance (Shade 1)",
		AttackTag:        attacks.AttackTagDirectLunarBloom,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		UseEM:            true,
		IgnoreDefPercent: 1,
		Mult:             phantasm[4][c.TalentLvlSkill()],
	}
	shadeHit2 := shadeHit1
	shadeHit2.Abil = "Phantasm Performance (Shade 2)"
	shadeHit2.Mult = phantasm[5][c.TalentLvlSkill()]
	shadeHit3 := shadeHit1
	shadeHit3.Abil = "Phantasm Performance (Shade 3)"
	shadeHit3.Mult = phantasm[6][c.TalentLvlSkill()]

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	c.Core.QueueAttack(neferHit1, ap, phantasmHit1, phantasmHit1)
	c.Core.QueueAttack(shadeHit1, ap, phantasmHit2, phantasmHit2)
	c.Core.QueueAttack(shadeHit2, ap, phantasmHit3, phantasmHit3)
	c.Core.QueueAttack(neferHit2, ap, phantasmHit4, phantasmHit4)
	c.Core.QueueAttack(shadeHit3, ap, phantasmHit5, phantasmHit5)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: 106,
		CanQueueAfter:   89,
		State:           action.ChargeAttackState,
	}, nil
}
