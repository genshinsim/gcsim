package qiqi

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

const burstHitmark = 82

func init() {
	burstFrames = frames.InitAbilSlice(115) // Q -> D
	burstFrames[action.ActionAttack] = 113  // Q -> N1
	burstFrames[action.ActionSkill] = 113   // Q -> E
	burstFrames[action.ActionJump] = 114    // Q -> J
	burstFrames[action.ActionSwap] = 112    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// Talisman is applied via a 0 dmg attack way before the damage is dealt
	talismanAi := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Fortune-Preserving Talisman (Talisman application)",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7)

	c.Core.QueueAttack(talismanAi, ap, 40, 40, c.talismanCB)

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Fortune-Preserving Talisman",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, ap, burstHitmark, burstHitmark)

	if c.revelation && c.getRadiance() == radianceStellarConduct {
		ai.Abil += stellarConductText
		ai.Mult = burstSSC[c.TalentLvlBurst()]
		ai.AttackTag = attacks.AttackTagDirectStellarConduct
		ai.IgnoreDefPercent = 1
		ai.Durability = 0
		ai.ICDTag = attacks.ICDTagNone
		c.Core.QueueAttack(ai, ap, burstHitmark, burstHitmark)
	}

	c.c6OnBurst()

	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(8)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) talismanCB(a info.AttackCB) {
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddStatus(talismanKey, 15*60, true)
}

func (c *char) burstInit() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		e, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		// do nothing if talisman expired
		if !e.StatusIsActive(talismanKey) {
			return
		}
		// do nothing if talisman still on icd
		if e.GetTag(talismanICDKey) >= c.Core.F {
			return
		}

		healAmt := c.healDynamic(burstHealPer, burstHealFlat, c.TalentLvlBurst())
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index(),
			Target:  atk.Info.ActorIndex,
			Message: "Fortune-Preserving Talisman",
			Src:     healAmt,
			Bonus:   c.Stat(attributes.Heal),
		})
		e.SetTag(talismanICDKey, c.Core.F+60)
		c.c4OnHeal()
	}, "talisman-heal-hook")
}
