package sethos

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var aimedFrames [][]int

var aimedHitmarks = []int{12, 83, 372}
var startCharge = aimedHitmarks[0]

const shadowPierceShotAil = "Shadowpiercing Shot"

func init() {
	// outside of E status
	aimedFrames = make([][]int, 3)

	// Aimed Shot
	aimedFrames[0] = frames.InitAbilSlice(23)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Fully-Charged Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(91)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Shadowpiercing Shot
	aimedFrames[2] = frames.InitAbilSlice(380)
	aimedFrames[2][action.ActionDash] = aimedHitmarks[2]
	aimedFrames[2][action.ActionJump] = aimedHitmarks[2]
}

func (c *char) Aimed(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(burstBuffKey) {
		return action.Info{}, fmt.Errorf("%v: Cannot aim while in burst", c.CharWrapper.Base.Key)
	}

	hold, ok := p["hold"]
	if !ok {
		// is this a good default? it's gonna take 6s to do without energy
		hold = attacks.AimParamLv2
	}
	switch hold {
	case attacks.AimParamPhys:
	case attacks.AimParamLv1:
	case attacks.AimParamLv2:
		return c.ShadowPierce(p)
	default:
		return action.Info{}, fmt.Errorf("invalid hold param supplied, got %v", hold)
	}

	skip, energy := c.a1Calc()
	if skip > aimedHitmarks[hold]-startCharge {
		skip = aimedHitmarks[hold] - startCharge
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	c.QueueCharTask(func() {
		ai := combat.AttackInfo{
			ActorIndex:           c.Index,
			Abil:                 "Fully-Charged Aimed Shot",
			AttackTag:            attacks.AttackTagExtra,
			ICDTag:               attacks.ICDTagNone,
			ICDGroup:             attacks.ICDGroupDefault,
			StrikeType:           attacks.StrikeTypePierce,
			Element:              attributes.Electro,
			Durability:           25,
			Mult:                 fullaim[c.TalentLvlAttack()],
			HitWeakPoint:         weakspot == 1,
			HitlagHaltFrames:     0.12 * 60,
			HitlagFactor:         0.01,
			HitlagOnHeadshotOnly: true,
			IsDeployable:         true,
		}
		if hold < attacks.AimParamLv1 {
			ai.Abil = "Aimed Shot"
			ai.Element = attributes.Physical
			ai.Mult = aim[c.TalentLvlAttack()]
		}

		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			0,
			travel,
		)
		c.a1Consume(energy, hold)
	}, aimedHitmarks[hold]-skip)

	return action.Info{
		Frames:          func(next action.Action) int { return aimedFrames[hold][next] - skip },
		AnimationLength: aimedFrames[hold][action.InvalidAction] - skip,
		CanQueueAfter:   aimedHitmarks[hold] - skip,
		State:           action.AimState,
	}, nil
}

func (c *char) ShadowPierce(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	skip, energy := c.a1Calc()
	if skip > aimedHitmarks[2]-startCharge {
		skip = aimedHitmarks[2] - startCharge
	}
	hitHaltFrames := 0.0
	if weakspot == 1 {
		hitHaltFrames = 0.12 * 60
	}

	c.QueueCharTask(func() {
		em := c.Stat(attributes.EM)
		ai := combat.AttackInfo{
			ActorIndex:           c.Index,
			Abil:                 shadowPierceShotAil,
			AttackTag:            attacks.AttackTagExtra,
			ICDTag:               attacks.ICDTagNone,
			ICDGroup:             attacks.ICDGroupDefault,
			StrikeType:           attacks.StrikeTypePierce,
			Element:              attributes.Electro,
			Durability:           50,
			Mult:                 shadowpierceAtk[c.TalentLvlAttack()],
			HitWeakPoint:         weakspot == 1,
			HitlagHaltFrames:     hitHaltFrames,
			HitlagFactor:         0.01,
			HitlagOnHeadshotOnly: true,
			IsDeployable:         true,
			FlatDmg:              shadowpierceEM[c.TalentLvlAttack()] * em,
		}

		if c.StatusIsActive(a4Key) {
			ai.FlatDmg += 7 * em
			c.Core.Log.NewEvent("Sethos A4 proc dmg add", glog.LogPreDamageMod, c.Index).
				Write("em", em).
				Write("ratio", 7.0).
				Write("addition", 7*em)
		}

		deltaPos := c.Core.Combat.Player().Pos().Sub(c.Core.Combat.PrimaryTarget().Pos())
		dist := deltaPos.Magnitude()

		// simulate piercing. Extends 15 units from player
		ap := combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -dist},
			0.1,
			15,
		)
		c.a1Consume(energy, attacks.AimParamLv2)
		c.Core.QueueAttack(
			ai,
			ap,
			0,
			travel,
			c.makeA4cb(),
			c.makeC4cb(),
			c.makeC6cb(energy),
		)
	}, aimedHitmarks[2]-skip)

	return action.Info{
		Frames:          func(next action.Action) int { return aimedFrames[2][next] - skip },
		AnimationLength: aimedFrames[2][action.InvalidAction] - skip,
		CanQueueAfter:   aimedHitmarks[2] - skip,
		State:           action.AimState,
	}, nil
}
