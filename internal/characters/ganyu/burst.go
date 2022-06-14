package ganyu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

const burstStart = 130

func init() {
	burstFrames = frames.InitAbilSlice(124)
	burstFrames[action.ActionSwap] = 122
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

	rad, ok := p["radius"]
	if !ok {
		rad = 1
	}
	r := 2.5 + float64(rad)
	prob := r * r / 90.25

	//tick every .3 sec, every fifth hit is targetted i.e. 1, 0, 0, 0, 0, 1
	//first hit at 148
	//duration is 15 seconds
	//starts from end of cast
	m := make([]float64, attributes.EndStatType)
	m[attributes.CryoP] = 0.2
	lastHit := make(map[combat.Target]int)
	for delay := burstStart; delay < 900+burstStart; delay += 18 {
		tick := delay
		c.Core.Tasks.Add(func() {

			// A4 every .3 seconds for the duration of the burst, add ice dmg up to active char for 1sec
			active := c.Core.Player.ActiveChar()
			active.AddStatMod("ganyu-field", 60, attributes.CryoP, func() ([]float64, bool) {
				return m, true
			})
			if tick >= 900+burstStart-18 {
				c.Core.Log.NewEvent("a4 last tick", glog.LogCharacterEvent, c.Index, "ends_on", c.Core.F+60)
			}

			// increase C4 stacks at 3s interval
			// assume this lasts for the full duration since no one moves...
			if c.Base.Cons >= 4 && (tick-burstStart)%180 == 0 {
				c.c4Stacks++
				if c.c4Stacks > 5 {
					c.c4Stacks = 5
				}
				c.Core.Log.NewEvent("ganyu c4 tick", glog.LogCharacterEvent, c.Index, "stacks", c.c4Stacks)
			}

			// damage ticks and C4
			//check if this hits first
			target := -1
			for i, t := range c.Core.Combat.Targets() {
				// skip non-enemy targets
				x, ok := t.(*enemy.Enemy)
				if !ok {
					continue
				}

				// C4 lingers for 3s
				if c.Base.Cons >= 4 {
					x.SetTag("ganuyc4", c.Core.F+60*3)
				}

				if lastHit[t] < c.Core.F {
					target = i
					lastHit[t] = c.Core.F + 87 //cannot be targetted again for 1.45s
					break
				}
			}
			// log.Println(target)
			//[1:14 PM] Aluminum | Harbinger of Jank: assuming uniform distribution and enemy at center:
			//(radius_icicle + radius_enemy)^2 / radius_burst^2
			if target == -1 && c.Core.Rand.Float64() > prob {
				//no one getting hit
				return
			}
			//deal dmg
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(9, false, combat.TargettableEnemy), 0)
		}, delay)

	}

	//add cooldown to sim
	c.SetCD(action.ActionBurst, 15*60)
	//use up energy
	c.ConsumeEnergy(3)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		Post:            burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
