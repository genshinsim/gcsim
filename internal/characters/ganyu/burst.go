package ganyu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstStart   = 130
	burstMarkKey = "ganyu-burst-mark"
)

func init() {
	burstFrames = frames.InitAbilSlice(125) // Q -> D/J
	burstFrames[action.ActionAttack] = 124  // Q -> N1
	burstFrames[action.ActionAim] = 124     // Q -> CA, assumed
	burstFrames[action.ActionSkill] = 124   // Q -> E
	burstFrames[action.ActionSwap] = 122    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Celestial Shower",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       shower[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	c.Core.Status.Add("ganyuburst", 15*60+burstStart)

	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10)
	// a4 related
	m := make([]float64, attributes.EndStatType)
	m[attributes.CryoP] = 0.2
	// tick every 0.3s from burstStart
	for i := 0; i < 15*60; i += 18 {
		// c4 related
		tick := i
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
				pos = combat.CalcRandomPointFromCenter(burstArea.Shape.Pos(), 0.5, 9.5, c.Core.Rand)
			}
			// deal dmg after a certain delay
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(pos, nil, 2.5), 10)

			// a4 buff tick
			if combat.TargetIsWithinArea(c.Core.Combat.Player(), burstArea) {
				active := c.Core.Player.ActiveChar()
				active.AddStatMod(character.StatMod{
					Base:         modifier.NewBase("ganyu-field", 60),
					AffectedStat: attributes.CryoP,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}
			// c4 debuff tick
			if c.Base.Cons >= 4 {
				enemies := c.Core.Combat.EnemiesWithinArea(burstArea, nil)
				// increase stacks every 3s but apply c4 status on every tick
				// c4 lingers for 3s
				increase := tick%180 == 0
				for _, e := range enemies {
					e.AddStatus(c4Key, c4Dur, true)
					if increase {
						c4Stacks := e.GetTag(c4Key) + 1
						if c4Stacks > 5 {
							c4Stacks = 5
						}
						e.SetTag(c4Key, c4Stacks)
						c.Core.Log.NewEvent(c4Key+" tick on enemy", glog.LogCharacterEvent, c.Index).
							Write("stacks", c4Stacks).
							Write("enemy key", e.Key())
					}
				}
			}
		}, i+burstStart)
	}

	//add cooldown to sim
	c.SetCD(action.ActionBurst, 15*60)
	//use up energy
	c.ConsumeEnergy(3)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
