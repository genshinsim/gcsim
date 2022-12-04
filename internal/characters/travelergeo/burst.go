package travelergeo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames [][]int

const burstStart = 35   // lines up with cooldown start
const burstHitmark = 51 // Initial Shockwave 1

func init() {
	burstFrames = make([][]int, 2)

	// Male
	burstFrames[0] = frames.InitAbilSlice(67) // Q -> N1/E
	burstFrames[0][action.ActionDash] = 42    // Q -> D
	burstFrames[0][action.ActionJump] = 42    // Q -> J
	burstFrames[0][action.ActionSwap] = 51    // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(64) // Q -> E
	burstFrames[1][action.ActionAttack] = 62  // Q -> N1
	burstFrames[1][action.ActionDash] = 42    // Q -> D
	burstFrames[1][action.ActionJump] = 42    // Q -> J
	burstFrames[1][action.ActionSwap] = 49    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	hits, ok := p["hits"]
	if !ok {
		// assume all 4 instances of shockwave dmg hit the enemy
		hits = 4
	}
	maxConstructCount, ok := p["construct_limit"]
	if !ok {
		// assume 2 barriers for zhongli pillar tick purposes (leads to 3 resonance ticks)
		// zhongli resonance limit is separate param on his skill (res_count)
		maxConstructCount = 2
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Wake of Earth",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagTravelerWakeOfEarth,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	//The shockwave triggered by Wake of Earth regenerates 5 Energy for every opponent hit.
	// A maximum of 25 Energy can be regenerated in this manner at any one time.
	src := c.Core.F
	var c4cb combat.AttackCBFunc
	if c.Base.Cons >= 4 {
		energyCount := 0
		c4cb = func(a combat.AttackCB) {
			t, ok := a.Target.(*enemy.Enemy)
			if !ok {
				return
			}
			// TODO: A bit of a cludge to deal with frame 0 casts. Will have to think about this behavior a bit more
			if t.GetTag("traveler-c4-src") == src && src > 0 {
				return
			}
			if energyCount >= 5 {
				return
			}
			t.SetTag("traveler-c4-src", src)
			c.AddEnergy("geo-traveler-c4", 5)
			energyCount++
		}
	}

	//1.1 sec duration, tick every .25
	for i := 0; i < hits; i++ {
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6),
			burstHitmark+(i+1)*15,
			c4cb,
		)
	}

	c.Core.Tasks.Add(func() {
		dur := 15 * 60
		if c.Base.Cons == 6 {
			dur += 300
		}
		c.Core.Constructs.NewNoLimitCons(c.newBarrier(dur, maxConstructCount), true)
		if c.Base.Cons >= 1 {
			c.Tags["wall"] = 1
		}
		// C1
		if c.Base.Cons >= 1 {
			c.Core.Tasks.Add(c.c1(1), 60) // start checking in 1s
		}
	}, burstStart)

	c.SetCDWithDelay(action.ActionBurst, 900, burstStart)
	c.ConsumeEnergy(37)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

type barrier struct {
	src    int
	expiry int
	char   *char
	count  int
}

func (c *char) newBarrier(dur, maxCount int) *barrier {
	return &barrier{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
		count:  maxCount,
	}
}

func (b *barrier) OnDestruct() {
	if b.char.Base.Cons >= 1 {
		b.char.Tags["wall"] = 0
	}
}

func (b *barrier) Key() int                         { return b.src }
func (b *barrier) Type() construct.GeoConstructType { return construct.GeoConstructTravellerBurst }
func (b *barrier) Expiry() int                      { return b.expiry }
func (b *barrier) IsLimited() bool                  { return true }
func (b *barrier) Count() int                       { return b.count }
