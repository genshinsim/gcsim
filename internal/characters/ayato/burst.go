package ayato

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstStart   = 101
	burstMarkKey = "ayato-burst-mark"
)

func init() {
	burstFrames = frames.InitAbilSlice(123) // Q -> N1
	burstFrames[action.ActionSkill] = 122   // Q -> E
	burstFrames[action.ActionDash] = 122    // Q -> D
	burstFrames[action.ActionJump] = 122    // Q -> J
	burstFrames[action.ActionSwap] = 120    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Kamisato Art: Suiyuu",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	// snapshot when the circle forms (is this correct?)
	var snap combat.Snapshot
	c.Core.Tasks.Add(func() { snap = c.Snapshot(&ai) }, burstStart)

	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10)
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = burstatkp[c.TalentLvlBurst()]
	// tick every 0.5s from burstStart
	for i := 0; i < 18*60; i += 30 {
		c.Core.Tasks.Add(func() {
			// burst tick
			enemy := c.Core.Combat.RandomEnemyWithinArea(
				burstArea,
				func(e combat.Enemy) bool {
					return !e.StatusIsActive(burstMarkKey)
				},
			)
			var pos combat.Point
			if enemy != nil {
				pos = enemy.Pos()
				enemy.AddStatus(burstMarkKey, 1.45*60, true) // same enemy can't be targeted again for 1.45s
			} else {
				pos = combat.CalcRandomPointFromCenter(burstArea.Shape.Pos(), 1.5, 9.5, c.Core.Rand)
			}
			// deal dmg after a certain delay
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(pos, nil, 2.5), 38)

			// buff tick
			if !combat.TargetIsWithinArea(c.Core.Combat.Player(), burstArea) {
				return
			}
			active := c.Core.Player.ActiveChar()
			active.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("ayato-burst", 90),
				Amount: func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					return m, a.Info.AttackTag == combat.AttackTagNormal
				},
			})
		}, i+burstStart)
	}

	if c.Base.Cons >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.AtkSpd] = 0.15
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("ayato-c4", 15*60),
				AffectedStat: attributes.AtkSpd,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
	}
	//add cooldown to sim
	c.SetCD(action.ActionBurst, 20*60)
	//use up energy
	c.ConsumeEnergy(5)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
