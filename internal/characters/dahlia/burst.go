package dahlia

import (
	"strconv"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	burstFrames []int

	favonianFavorDuration int
	favonianFavorExpiry   int

	totalBenisonStacks   int
	currentBenisonStacks int

	normalAttackCount        int
	currentNormalAttackFrame int
	currentEnemyTarget       info.TargetKey
)

const (
	burstHitmark                 = 30
	burstShieldInitial           = 31 // First shield
	burstShieldAfterBenisonStack = 31 // No shield -> new shield once a Benison stack is added
	burstShieldRegenerated       = 26 // Broken shield -> new shield
	burstEnergyDrain             = 8
	burstCDStart                 = 0

	burstFavonianFavor    = "dahlia-favonian-favor"
	burstBenisonStacksKey = "dahlia-benison-stacks"
)

func init() {
	burstFrames = frames.InitAbilSlice(55) // Q -> W
	burstFrames[action.ActionAttack] = 53  // Q -> N1
	burstFrames[action.ActionSkill] = 52   // Q -> -> tE / hE (TO-DO: check if this includes both)
	burstFrames[action.ActionDash] = 51    // Q -> D
	burstFrames[action.ActionJump] = 51    // Q -> J
	burstFrames[action.ActionSwap] = 52    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	totalBenisonStacks = 0 // TO-DO: Should these be here?
	currentBenisonStacks = 0
	normalAttackCount = 0
	currentNormalAttackFrame = -1
	currentEnemyTarget = -1

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Radiant Psalter",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	snap := c.Snapshot(&ai) //?? TO-DO: Fix
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 0.5}, 4), //?? TO-DO: Fix values
		burstHitmark,
	)

	c.SetCDWithDelay(action.ActionBurst, 15*60, burstCDStart) // Delay is actually 0 but left for clarity
	c.ConsumeEnergy(burstEnergyDrain)

	// After a short delay, add the "Favonian Favor" status and all its effects
	c.Core.Tasks.Add(func() {
		// Increasing shield duration if C4+ (it lasts for exactly 15s, so shield remains after Q CD is over)
		favonianFavorDuration = 12
		if c.Base.Cons >= 4 {
			favonianFavorDuration = 15
		}
		favonianFavorExpiry = c.Core.F + favonianFavorDuration*60

		// Add "Favonian Favor" status for all party members
		for _, char := range c.Core.Player.Chars() {
			char.AddStatus(burstFavonianFavor, favonianFavorExpiry, true)
		}

		// Create shield
		c.genShield()

		// Apply Attack Speed buff to Dahlia
		c.addAttackSpeedbuff(c.Core.Player.ActiveChar())
	}, burstShieldInitial)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

// Benison stack generation, max 4 stacks
func (c *char) setupBurst() {
	// Add Benison stack when 4 hits from NAs occur
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		char := c.Core.Player.ActiveChar()

		if !char.StatusIsActive(burstFavonianFavor) {
			return false
		}

		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		ae := args[1].(*info.AttackEvent)
		if ae.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}

		// If already generated a total of 4 Benison stacks, do nothing
		if c.isMaxBenisonStacks() {
			return false
		}

		// If multiple enemies are hit by the same NA hit, don't increment the hit counter
		if ae.SourceFrame == currentNormalAttackFrame && e.Key() != currentEnemyTarget {
			return false
		}
		currentNormalAttackFrame = ae.SourceFrame
		currentEnemyTarget = e.Key()

		// If 4 hits from NAs occured, add Benison stack
		normalAttackCount++
		if normalAttackCount == 4 {
			c.addBenisonStack(1, ae.Info.ActorIndex)
			normalAttackCount = 0
		}

		return false
	}, burstBenisonStacksKey)
}

func (c *char) addBenisonStack(stacks int, charIndex int) {
	if c.isMaxBenisonStacks() { // TO-DO: Regarding setupBurst(), this is technically never true, but it's used elsewhere
		return
	}

	// If stacks to add exceed max, only add up to max
	stacksToAdd := stacks
	if totalBenisonStacks+stacks > c.maxBenisonStacks {
		stacksToAdd = c.maxBenisonStacks - totalBenisonStacks
	}
	totalBenisonStacks += stacksToAdd
	currentBenisonStacks += stacksToAdd

	// If C1, Dahlia restores 2.5 Energy per stack gained
	if c.Base.Cons >= 1 {
		c.AddEnergy(c1Key, 2.5*float64(stacksToAdd))
	}

	c.Core.Log.NewEvent(strconv.Itoa(stacksToAdd)+" "+burstBenisonStacksKey+" added", glog.LogShieldEvent, charIndex).
		Write("benison_stacks_remaining", currentBenisonStacks).
		Write("max_benison_stacks_reached", c.isMaxBenisonStacks())

	// If shield is already gone but new stacks got generated, create shield (after some delay)
	if !c.hasShield() && favonianFavorExpiry > c.Core.F+burstShieldAfterBenisonStack {
		c.QueueCharTask(func() {
			c.genShield()
			c.c2()
		}, burstShieldAfterBenisonStack)
	}
}

func (c *char) isMaxBenisonStacks() bool {
	return totalBenisonStacks >= c.maxBenisonStacks
}
