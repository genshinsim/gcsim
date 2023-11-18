package furina

import (
	"math"

	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstHitmark = 98
	burstKey     = "furina-burst"
)

func init() {
	burstFrames = frames.InitAbilSlice(121)
	burstFrames[action.ActionAttack] = 113 // Q -> N1
	burstFrames[action.ActionCharge] = 113 // Q -> CA
	burstFrames[action.ActionSkill] = 114  // Q -> E
	burstFrames[action.ActionDash] = 115   // Q -> D
	burstFrames[action.ActionJump] = 115   // Q -> J
	burstFrames[action.ActionSwap] = 111   // Q -> Swap
}

func (c *char) burstInit() {
	c.maxFanfare = 300
	c.maxQFanfare = 300
	if c.Base.Cons >= 1 {
		c.maxQFanfare = 400
		c.maxFanfare = 400
	}
	if c.Base.Cons >= 2 {
		// 400 + 140/0.35 = 800
		c.maxFanfare = 800
	}
	c.curFanfare = 0
	c.burstBuff = make([]float64, attributes.EndStatType)

	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("furina-burst-damage-buff", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				c.burstBuff[attributes.DmgP] = common.Min(c.curFanfare, c.maxQFanfare) * burstFanfareDMGRatio[c.TalentLvlBurst()]
				return c.burstBuff, c.StatusIsActive(burstKey)
			},
		})

		char.AddHealBonusMod(character.HealBonusMod{
			Base: modifier.NewBase("furina-burst-heal-buff", -1),
			Amount: func() (float64, bool) {
				return common.Min(c.curFanfare, c.maxQFanfare) * burstFanfareHBRatio[c.TalentLvlBurst()], c.StatusIsActive(burstKey)
			},
		})
	}

	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}

		di := args[0].(player.DrainInfo)

		if di.Amount <= 0 {
			return false
		}

		char := c.Core.Player.ByIndex(di.ActorIndex)
		stacksAmount := di.Amount / char.MaxHP() * 100

		if c.Base.Cons >= 2 {
			stacksAmount *= 3.5
		}

		c.curFanfare = common.Min(c.maxFanfare, c.curFanfare+stacksAmount)

		return false
	}, "furina-fanfare-on-hp-drain")

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}

		target := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)

		if amount <= 0 {
			return false
		}

		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}

		char := c.Core.Player.ByIndex(target)
		stacksAmount := (amount - overheal) / char.MaxHP() * 100

		if c.Base.Cons >= 2 {
			stacksAmount *= 3.5
		}

		c.curFanfare = common.Min(c.maxFanfare, c.curFanfare+stacksAmount)

		return false
	}, "furina-fanfare-on-heal")
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Let the People Rejoice",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    c.MaxHP() * burstDMG[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), burstHitmark, burstHitmark)
	c.QueueCharTask(func() {
		c.curFanfare = 0

		if c.Base.Cons >= 1 {
			c.curFanfare = 150
		}

		c.AddStatus(burstKey, 18*60, false)
	}, burstHitmark+1)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
