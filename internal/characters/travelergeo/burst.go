package travelergeo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
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
		hits = 4 // assume all 4 instances of shockwave dmg hit the enemy
	}
	maxConstructCount, ok := p["construct_limit"]
	if !ok {
		// assume all 4 walls actually spawn
		// going lower than 4 starts not spawning walls from top left, going counterclockwise
		maxConstructCount = 4
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Wake of Earth",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagTravelerWakeOfEarth,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	// C4
	// The shockwave triggered by Wake of Earth regenerates 5 Energy for every opponent hit.
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
	player := c.Core.Combat.Player()
	c.burstArea = combat.NewCircleHitOnTarget(player, nil, 7)
	// 1.1 sec duration, tick every .25
	for i := 0; i < hits; i++ {
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(c.burstArea.Shape.Pos(), nil, 6),
			burstHitmark+(i+1)*15,
			c4cb,
		)
	}

	// 4 walls spawn at +-2.75, +-6.67 assuming default viewing direction
	// spawning starts top right, and goes clockwise
	// if you rotate (2.75, 6.67) counterclockwise until ending up with (0, x), then the angle is around 22.5
	// this angle gets used for determining the wall's viewing direction
	angles := []float64{22.5, 112.5, 202.5, 292.5}
	offsets := []combat.Point{{X: 2.75, Y: 6.67}, {X: 2.75, Y: -6.67}, {X: -2.75, Y: -6.67}, {X: -2.75, Y: 6.67}}
	c.Core.Tasks.Add(func() {
		// C1
		// Party members within the radius of Wake of Earth have their CRIT Rate increased by 10% and have increased resistance against interruption.
		if c.Base.Cons >= 1 {
			c.Tags["wall"] = 1
		}
		if c.Base.Cons >= 1 {
			c.Core.Tasks.Add(c.c1(1), 60) // start checking in 1s
		}
		// C6
		// The barrier created by Wake of Earth lasts 5s longer.
		// The meteorite created by Starfell Sword lasts 10s longer.
		dur := 15 * 60
		if c.Base.Cons >= 6 {
			dur += 300
		}
		// spawn walls up until the specified limit is reached
		for i := 0; i < maxConstructCount; i++ {
			dir := combat.DegreesToDirection(angles[i]).Rotate(player.Direction())
			pos := combat.CalcOffsetPoint(player.Pos(), offsets[i], player.Direction())
			c.Core.Constructs.NewNoLimitCons(c.newWall(dur, dir, pos), false)
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

type wall struct {
	src    int
	expiry int
	char   *char
	dir    combat.Point
	pos    combat.Point
}

func (c *char) newWall(dur int, dir, pos combat.Point) *wall {
	return &wall{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
		dir:    dir,
		pos:    pos,
	}
}

func (w *wall) OnDestruct() {
	if w.char.Base.Cons >= 1 {
		w.char.Tags["wall"] = 0
	}
}

func (w *wall) Key() int                         { return w.src }
func (w *wall) Type() construct.GeoConstructType { return construct.GeoConstructTravellerBurst }
func (w *wall) Expiry() int                      { return w.expiry }
func (w *wall) IsLimited() bool                  { return true }
func (w *wall) Count() int                       { return 1 }
func (w *wall) Direction() combat.Point          { return w.dir }
func (w *wall) Pos() combat.Point                { return w.pos }
