package dahlia

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstHitmark                 = 30
	burstShieldInitial           = 31 // First shield
	burstShieldAfterBenisonStack = 31 // No shield -> new shield once a Benison stack is added
	burstShieldRegenerated       = 26 // Broken shield -> new shield
	burstEnergyDrain             = 8
	burstCDStart                 = 0

	benisonMaxGenerate  = 4
	maxFavonianFavorExt = 5 * 60 // Checked using AS buffs and on-field Dahlia full NA string

	burstFavonianFavor    = "dahlia-favonian-favor"
	burstBenisonStacksKey = "dahlia-benison-stacks"

	normalAttackStackIcd = "dahlia-na-stack-icd"
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
	c.benisonGenStackLimit = benisonMaxGenerate
	c.currentBenisonStacks = 0
	c.normalAttackCount = 0

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

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 0.8}, 5.5),
		burstHitmark, // TODO: assumed snapshot at burstHitmark, check if true
		burstHitmark,
	)

	c.SetCDWithDelay(action.ActionBurst, 15*60, burstCDStart) // Delay is actually 0 but left for clarity
	c.ConsumeEnergy(burstEnergyDrain)

	// After a short delay, add the "Favonian Favor" status and all its effects
	c.Core.Tasks.Add(func() {
		// Increasing Favonisn Favor duration if C4+
		// It lasts for exactly 15s (pre-hitlag), so shield + Attack Speed remain after Q CD is over
		favonianFavorDuration := 12*60 + c.c4FavonianFavorBonusDur()
		c.favonianFavorMaxExpiry = c.Core.F + favonianFavorDuration + maxFavonianFavorExt

		// Add "Favonian Favor" status to Dahlia
		// (team members should technically get it but only he can extend it with hitlag)
		c.AddStatus(burstFavonianFavor, favonianFavorDuration, true)
		c.QueueCharTask(func() {
			c.removeShield()
		}, favonianFavorDuration)

		// Create shield
		c.genShield()

		// Start calculating Attack Speed buff
		c.updateSpeedBuff(c.Core.Player.Active())()
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
		if !c.StatusIsActive(burstFavonianFavor) {
			return false
		}

		char := c.Core.Player.ActiveChar()
		ae := args[1].(*info.AttackEvent)
		if char.Index() != ae.Info.ActorIndex {
			return false
		}
		if ae.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}

		if c.StatusIsActive(normalAttackStackIcd) {
			return false
		}

		// If already generated a total of 4 Benison stacks, do nothing
		if c.benisonGenStackLimit == 0 {
			return false
		}

		// If 4 hits from NAs occured, add Benison stack
		c.normalAttackCount++
		if c.normalAttackCount == 4 {
			c.addBenisonStack(1, ae.Info.ActorIndex)
			c.normalAttackCount = 0
		}
		c.AddStatus(normalAttackStackIcd, 0.05*60, true)

		return false
	}, burstBenisonStacksKey)
}

func (c *char) addBenisonStack(stacks, charIndex int) {
	if c.benisonGenStackLimit <= 0 {
		return
	}
	stacks = min(stacks, c.benisonGenStackLimit)
	c.benisonGenStackLimit -= stacks
	c.currentBenisonStacks += stacks

	// If C1, Dahlia restores 2.5 Energy per stack gained
	c.c1OnBenisonEnergy(stacks)

	c.Core.Log.NewEvent(fmt.Sprintf("%v %v added", stacks, burstBenisonStacksKey), glog.LogShieldEvent, charIndex).
		Write("benison_stacks_remaining", c.currentBenisonStacks).
		Write("benison_stacks_generated", benisonMaxGenerate-c.benisonGenStackLimit)

	// If shield is already gone but new stacks got generated, create shield (after some delay)
	if !c.hasShield() && c.StatusExpiry(burstFavonianFavor) > c.Core.F+burstShieldAfterBenisonStack {
		c.Core.Tasks.Add(func() {
			c.currentBenisonStacks--
			c.genShield()
			c.c2()
		}, burstShieldAfterBenisonStack)
	}
}
