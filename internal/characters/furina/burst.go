package furina

import (
	"math"

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
	burstDur     = 18.2 * 60
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

func (c *char) addFanfare(amt float64) {
	if c.Base.Cons >= 2 {
		amt *= 3.5
	}
	c.curFanfare = math.Min(c.maxC2Fanfare, c.curFanfare+amt)
}

func (c *char) burstInit() {
	c.maxC2Fanfare = 300
	c.maxQFanfare = 300
	if c.Base.Cons >= 1 {
		c.maxQFanfare = 400
		c.maxC2Fanfare = 400
	}
	if c.Base.Cons >= 2 {
		// 400 + 140/0.35 = 800
		c.maxC2Fanfare = 800
	}
	c.burstBuff = make([]float64, attributes.EndStatType)

	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}

		di := args[0].(player.DrainInfo)

		if di.Amount <= 0 {
			return false
		}

		char := c.Core.Player.ByIndex(di.ActorIndex)
		amt := di.Amount / char.MaxHP() * 100
		c.addFanfare(amt)

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
		amt := (amount - overheal) / char.MaxHP() * 100

		c.addFanfare(amt)

		return false
	}, "furina-fanfare-on-heal")

	burstDMGRatio := burstFanfareDMGRatio[c.TalentLvlBurst()]
	burstHealRatio := burstFanfareHBRatio[c.TalentLvlBurst()]
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("furina-burst-damage-buff", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if c.StatusIsActive(burstKey) {
					c.burstBuff[attributes.DmgP] = math.Min(c.curFanfare, c.maxQFanfare) * burstDMGRatio
				} else {
					c.burstBuff[attributes.DmgP] = 0
				}
				return c.burstBuff, true
			},
		})

		char.AddHealBonusMod(character.HealBonusMod{
			Base: modifier.NewBase("furina-burst-heal-buff", -1),
			Amount: func() (float64, bool) {
				if c.StatusIsActive(burstKey) {
					return math.Min(c.curFanfare, c.maxQFanfare) * burstHealRatio, false
				}
				return 0, false
			},
		})
	}
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

	c.curFanfare = 0
	c.DeleteStatus(burstKey)

	c.QueueCharTask(func() {
		if c.Base.Cons >= 1 {
			c.curFanfare = 150
		}
		c.AddStatus(burstKey, burstDur, true)
	}, 95)

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), burstHitmark, burstHitmark)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
