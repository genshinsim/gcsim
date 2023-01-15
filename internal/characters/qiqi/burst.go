package qiqi

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
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

// Only applies burst damage. Main Talisman functions are handled in qiqi.go
func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fortune-Preserving Talisman",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7)
	c.Core.QueueAttack(ai, ap, burstHitmark, burstHitmark)

	// Talisman is applied way before the damage is dealt
	c.Core.Tasks.Add(func() {
		for _, e := range c.Core.Combat.EnemiesWithinArea(ap, nil) {
			e.AddStatus(talismanKey, 15*60, true)
		}
	}, 40)

	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) talismanHealHook() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		e, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}

		//do nothing if talisman expired
		if !e.StatusIsActive(talismanKey) {
			return false
		}
		//do nothing if talisman still on icd
		if e.GetTag(talismanICDKey) >= c.Core.F {
			return false
		}

		healAmt := c.healDynamic(burstHealPer, burstHealFlat, c.TalentLvlBurst())
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  atk.Info.ActorIndex,
			Message: "Fortune-Preserving Talisman",
			Src:     healAmt,
			Bonus:   c.Stat(attributes.Heal),
		})
		e.SetTag(talismanICDKey, c.Core.F+60)

		return false
	}, "talisman-heal-hook")
}

// Handles C2, A4, and skill NA/CA on hit hooks
// Additionally handles burst Talisman hook - can't be done another way since Talisman is applied before the burst damage is dealt
func (c *char) onNACAHitHook() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		e, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		// All of the below only occur on Qiqi NA/CA hits
		switch atk.Info.AttackTag {
		case combat.AttackTagNormal:
		case combat.AttackTagExtra:
		default:
			return false
		}

		// A4
		// When Qiqi hits opponents with her Normal and Charged Attacks, she has a 50% chance to apply a Fortune-Preserving Talisman to them for 6s. This effect can only occur once every 30s.
		if !c.StatusIsActive(a4ICDKey) && (c.Core.Rand.Float64() < 0.5) {
			// Don't want to overwrite a longer burst duration talisman with a shorter duration one
			// TODO: Unclear how the interaction works if there is already a talisman on enemy
			// TODO: Being generous for now and not putting it on CD if there is a conflict
			if e.StatusExpiry(talismanKey) < c.Core.F+360 {
				e.AddStatus(talismanKey, 360, true)
				c.AddStatus(a4ICDKey, 1800, true) // 30s icd
				c.Core.Log.NewEvent(
					"Qiqi A4 Adding Talisman",
					glog.LogCharacterEvent,
					c.Index,
				).
					Write("target", e.Key()).
					Write("talisman_expiry", e.StatusExpiry(talismanKey))
			}
		}

		// Qiqi NA/CA healing proc in skill duration
		if c.StatusIsActive(skillBuffKey) {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Herald of Frost (Attack)",
				Src:     c.healSnapshot(&c.skillHealSnapshot, skillHealOnHitPer, skillHealOnHitFlat, c.TalentLvlSkill()),
				Bonus:   c.skillHealSnapshot.Stats[attributes.Heal],
			})
		}

		return false
	}, "qiqi-onhit-naca-hook")
}
