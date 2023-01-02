package ayato

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstStart = 101

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

	rad, ok := p["radius"]
	if !ok {
		rad = 1
	}

	r := 2.5 + float64(rad)
	prob := r * r / 90.25

	lastHit := make(map[combat.Target]int)
	// ccc := 0
	//tick every .5 sec, every fourth hit is targetted i.e. 1, 0, 0, 0, 1
	dur := 18
	for delay := 0; delay < dur*60; delay += 30 {
		c.Core.Tasks.Add(func() {
			//check if this hits first
			target := -1
			for i, t := range c.Core.Combat.Enemies() {
				// skip non-enemy targets
				if _, ok := t.(*enemy.Enemy); !ok {
					continue
				}
				if lastHit[t] < c.Core.F {
					target = i
					lastHit[t] = c.Core.F + 117 //cannot be targetted again for 1.95s
					break
				}
			}
			// log.Println(target)
			//[1:14 PM] Aluminum | Harbinger of Jank: assuming uniform distribution and enemy at center:
			//(radius_droplet + radius_enemy)^2 / radius_burst^2
			trg := c.Core.Combat.Enemy(target)
			if target == -1 {
				if c.Core.Rand.Float64() > prob {
					// no one getting hit
					return
				} else {
					// droplet is not targeted but randomly clips enemy
					// TODO: enemies within radius?
					trg = c.Core.Combat.Enemy(c.Core.Combat.RandomEnemyTarget())
				}
			}
			//deal dmg
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(trg, nil, 2.5), 0)
		}, delay+139)
	}

	c.Core.Status.Add("ayatoburst", dur*60+burstStart)

	// NA buff starts after cast, ticks every 0.5s and last for 1.5s
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = burstatkp[c.TalentLvlBurst()]
	for i := burstStart; i < burstStart+dur*60; i += 30 {
		c.Core.Tasks.Add(func() {
			active := c.Core.Player.ActiveChar()
			active.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("ayato-burst", 90),
				Amount: func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					return m, a.Info.AttackTag == combat.AttackTagNormal
				},
			})
		}, i)
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
