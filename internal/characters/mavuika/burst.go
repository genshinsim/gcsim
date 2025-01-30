package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	burstKey       = "mavuika-burst"
	energyNAICDKey = "mavuika-fighting-spirit-na-icd"
	burstDuration  = 7.0 * 60
	burstHitmark   = 106
)

var (
	burstFrames []int
)

func (c *char) nightsoulConsumptionMul() float64 {
	if c.StatusIsActive(burstKey) {
		return 0.0
	}
	return 1.0
}

func init() {
	burstFrames = frames.InitAbilSlice(116) // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.burstStacks = c.fightingSpirit
	c.fightingSpirit = 0
	c.enterBike()
	c.QueueCharTask(func() {
		c.enterNightsoulOrRegenerate(10)
	}, 87)
	c.QueueCharTask(func() {
		c.AddStatus(burstKey, burstDuration, true)
	}, burstHitmark-1)
	c.QueueCharTask(func() {
		c.a4()

		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Sunfell Slice",
			AttackTag:      attacks.AttackTagElementalBurst,
			ICDTag:         attacks.ICDTagNone,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeBlunt,
			PoiseDMG:       150,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           burst[c.TalentLvlBurst()],
			FlatDmg:        c.burstBuffSunfell() + c.c2BikeQ(),
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: 1.0},
			6,
		)
		c.Core.QueueAttack(ai, ap, 0, 0)
	}, burstHitmark)

	c.SetCDWithDelay(action.ActionBurst, 18*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstBuffCA() float64 {
	if !c.StatusIsActive(burstKey) {
		return 0.0
	}
	return c.burstStacks * burstCABonus[c.TalentLvlBurst()] * c.TotalAtk()
}

func (c *char) burstBuffNA() float64 {
	if !c.StatusIsActive(burstKey) {
		return 0.0
	}
	return c.burstStacks * burstNABonus[c.TalentLvlBurst()] * c.TotalAtk()
}

func (c *char) burstBuffSunfell() float64 {
	if !c.StatusIsActive(burstKey) {
		return 0.0
	}
	return c.burstStacks * burstQBonus[c.TalentLvlBurst()] * c.TotalAtk()
}

func (c *char) gainFightingSpirit(val float64) {
	c.fightingSpirit += val * c.c1FightingSpiritEff()
	if c.fightingSpirit > 200 {
		c.fightingSpirit = 200
	}
	c.c1OnFightingSpirit()
}

func (c *char) burstInit() {
	c.fightingSpirit = 200
	c.Core.Events.Subscribe(event.OnNightsoulConsume, func(args ...interface{}) bool {
		amount := args[1].(float64)
		if amount < 0.0000001 {
			return false
		}
		c.gainFightingSpirit(amount)
		return false
	}, "mavuika-fighting-spirit-ns")

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if ae.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}
		if c.StatusIsActive(energyNAICDKey) {
			return false
		}
		c.AddStatus(energyNAICDKey, 0.1*60, true)
		c.gainFightingSpirit(1.5)
		return false
	}, "mavuika-fighting-spirit-na")
}
