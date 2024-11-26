package mavuika

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	burstHitmark         = 99
	minFightingSpiritReq = 50

	crucibleOfDeathAndLifeStatus = "crucible-of-death-and-life"
	fightingSpiritGainIcd        = "fighting-spirit-gain-icd"
)

var (
	burstFrames []int
)

func init() {
	burstFrames = frames.InitAbilSlice(135)
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.fightingSpirit < minFightingSpiritReq {
		return action.Info{}, fmt.Errorf("%v: Cannot Burst with %v Fighting Spirit, should be at least %v",
			c.Base.Key, c.fightingSpirit, minFightingSpiritReq)
	}

	// I assume her burst does NOT change her stance (bike/no bike) whatsoever.
	// upon entering Nightsoul, she appears in the stance she was, when she exited Nightsoul.
	if c.nightsoulState.HasBlessing() {
		c.nightsoulState.GeneratePoints(10)
	} else {
		c.nightsoulState.EnterBlessing(10)
		c.c2BaseIncrease(true)
	}
	c.nightsoulPointReduceFunc(c.Core.F)

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Burst DMG",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Pyro,
		Durability:     25,
		FlatDmg:        c.TotalAtk() * 8.006,
	}
	c.SetCD(action.ActionBurst, 18*60)

	// consume fighting spirit
	c.consumedFightingSpirit = c.fightingSpirit
	c.fightingSpirit = 0

	ai.FlatDmg += (0.029*c.TotalAtk()*float64(c.consumedFightingSpirit) + c.c2FlatIncrease(attacks.AttackTagElementalBurst))
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 7), burstHitmark, burstHitmark)

	// Start countdown after initial hit
	c.QueueCharTask(func() {
		c.AddStatus(crucibleOfDeathAndLifeStatus, 7*60, true)
	}, burstHitmark+1)

	// Activate A4 without delay. TODO: confirm
	c.a4()

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionAttack], // change to earliest
		State:           action.BurstState,
	}, nil
}

func (c *char) burstInit() {
	c.maxFightingSpirit = 200

	c.Core.Events.Subscribe(event.OnNightsoulConsume, func(args ...interface{}) bool {
		amount := args[1].(float64)
		c.fightingSpirit = c.fightingSpiritMult * min(c.maxFightingSpirit, c.fightingSpirit+amount)
		c.c1Atk()
		return false
	}, "mavuika-fighting-spirit-on-ns-consumption")

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		_, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}
		if c.StatusIsActive(fightingSpiritGainIcd) {
			return false
		}
		c.fightingSpirit = c.fightingSpiritMult * min(c.maxFightingSpirit, c.fightingSpirit+1.5)
		c.AddStatus(fightingSpiritGainIcd, 0.1*60, false)
		c.c1Atk()
		return false
	}, "mavuika-fighting-spirit-on-na-hit")
}
