package sandrone

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	chargeEndFrames      []int
	chargeFirstTick      = 18
	overdriveFirstTick   = 21
	sweepingFireTickRate = []int{24, 18, 25, 15, 18, 24, 18, 18, 22, 18, 18, 24, 18, 18, 24, 18, 18, 24, 18, 18, 24, 18, 18, 24}
	beamTickRate         = 60
)

func init() {
	chargeEndFrames = frames.InitAbilSlice(20) // CA -> N1/CA
	chargeEndFrames[action.ActionAttack] = 2
	chargeEndFrames[action.ActionCharge] = 28
	chargeEndFrames[action.ActionSkill] = 0
	chargeEndFrames[action.ActionBurst] = 0
	chargeEndFrames[action.ActionDash] = 0
	chargeEndFrames[action.ActionJump] = 0
	chargeEndFrames[action.ActionWalk] = 16
	chargeEndFrames[action.ActionSwap] = 1
}

// Charged Attack
// Summons the Facilitatrice Automatisee Generaliste Interactive et Omnitache, or "Fagio," and
// switches it to Decoding mode: Unleashes a sweeping fire attack on enemies in front, as well as
// periodically fires condensed beams, dealing AoE Cryo DMG.
//
// Radiance: Stellar Glimmer: The condensed beams will instead deal AoE Cryo DMG that is considered
// the corresponding Stellar Glimmer reaction DMG.
//
// While in Decoding mode, Fagio's Decoding Power will continuously increase. Firing condensed beams
// will also further increase its Decoding Power.
//
// When Fagio's Decoding Power reaches the maximum of 100 points, it switches to Power Overdrive
// mode, in which it fires at opponents with longer intervals, dealing Cryo DMG. It will be unable
// to switch to Decoding mode at this time.
//
// When not in Decoding mode, Fagio's Decoding Power gradually decreases. When Sandrone is
// off-field, she will focus on maintaining Fagio, causing its Decoding Power to decrease at 300%
// the original rate.
//
// Fagio exits Power Overdrive mode when its Decoding Power has decreased to under 50.
func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	windup := 0
	switch c.Core.Player.CurrentState() {
	// other states have windup included in the endlag animation
	case action.Idle, action.JumpState, action.DashState, action.SwapState:
		windup = 25
	}
	c.caLength = 0
	c.currSweepingFireTick = 0
	c.caBeamsRemaining = c.caMaxBeams
	c.caToBeams = true

	beams, okBeams := p["beams"]
	frames, okFrames := p["frames"]

	if okBeams && okFrames {
		return action.Info{}, fmt.Errorf("cannot use both 'beams' and 'frames'")
	}

	if okBeams {
		if beams > c.caMaxBeams {
			return action.Info{}, fmt.Errorf("cannot set than %v beams", c.caMaxBeams)
		}
		c.caBeamsRemaining = beams
	}

	if okFrames {
		c.caToBeams = false
		c.caLength = frames
		c.caBeamsRemaining = math.MaxInt
		c.Core.Tasks.Add(c.exitCA, windup+chargeFirstTick+c.caLength)
	}

	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEvent(fmt.Sprintf("CA started with Decode (%.2f)", c.decode), glog.LogCharacterEvent, c.Index())
	}

	c.c2OnCAStart()
	c.c6OnCaStart()
	if c.currFagio == stateOverdriveIdle {
		c.currFagio = stateOverdriveCharge
		c.QueueCharTask(func() {
			src := c.Core.F
			c.caSrc = src
			c.decodeReduceSrc = -1
			c.overdriveTicker(src)
		}, windup+overdriveFirstTick)
	} else {
		// heat starts ticking from windup + first tick
		c.currFagio = stateDecoding
		c.QueueCharTask(func() {
			src := c.Core.F
			c.caSrc = src
			c.decodeReduceSrc = -1
			c.gainHeatTicker(src)
			c.sweepingFireTicker(src)
			c.Core.Tasks.Add(func() { c.beamTicker(src) }, 60)
		}, windup+chargeFirstTick)
	}

	return action.Info{
		Frames: func(next action.Action) int {
			return windup + chargeFirstTick + c.caLength + chargeEndFrames[next] + 1
		},
		AnimationLength: math.MaxInt, // there is no upper limit on the duration of the CA
		CanQueueAfter:   windup + chargeFirstTick + chargeEndFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) exitCA() {
	// exit the CA at the end of this frame
	c.Core.Tasks.Add(func() {
		switch c.currFagio {
		case stateDecoding:
			c.currFagio = stateIdle
		case stateOverdriveCharge:
			c.currFagio = stateOverdriveIdle
		}
		c.caLength = c.Core.F - c.caSrc
		c.caSrc = -1
		src := c.Core.F
		c.decodeReduceSrc = src
		c.Core.Tasks.Add(func() { c.heatReduceTicker(src) }, 0.5*60)
	}, 0)
}

func (c *char) sweepingFireTicker(src int) {
	if c.caSrc != src {
		return
	}

	if c.currFagio != stateDecoding {
		return
	}

	if c.caLength+c.caSrc < c.Core.F {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Sweeping Fire",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagSandroneExtraAttackSweepingFire,
		ICDGroup:   attacks.ICDGroupSandroneSweepingFire,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       chargeSweep[c.TalentLvlAttack()],
	}

	// approximate the gadget hitbox as a long rectangle
	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: -7.5},
		3,
		15,
	)

	c.Core.QueueAttack(ai, ap, 0, 0)

	delay := sweepingFireTickRate[c.currSweepingFireTick]
	c.currSweepingFireTick += 1
	c.currSweepingFireTick %= len(sweepingFireTickRate)
	c.Core.Tasks.Add(func() { c.sweepingFireTicker(src) }, delay)
}

func (c *char) overdriveTicker(src int) {
	if c.caSrc != src {
		return
	}

	if c.currFagio != stateOverdriveCharge {
		return
	}

	if c.caLength+c.caSrc < c.Core.F {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Power Overdrive",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagSandroneExtraAttackSweepingFire,
		ICDGroup:   attacks.ICDGroupSandroneSweepingFire,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       chargeOverdrive[c.TalentLvlAttack()],
	}

	// approximate the gadget hitbox as a long rectangle
	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: -7.5},
		3,
		15,
	)
	c.Core.QueueAttack(ai, ap, 0, 10)

	delay := 24
	c.Core.Tasks.Add(func() { c.overdriveTicker(src) }, delay)
}

func (c *char) beamTicker(src int) {
	if c.caSrc != src {
		return
	}

	if c.currFagio != stateDecoding {
		return
	}

	if c.caLength+c.caSrc < c.Core.F {
		return
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Condensed Beam",
		AttackTag:      attacks.AttackTagExtra,
		AdditionalTags: []attacks.AttackTag{attacks.AttackTagSandroneBeam},
		ICDTag:         attacks.ICDTagSandroneExtraAttackLaser,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       45,
		Element:        attributes.Cryo,
		Durability:     25,
		Mult:           chargeBeam[c.TalentLvlAttack()],
	}
	c.toStellar(&ai, chargeBeamSSC[c.TalentLvlAttack()], chargeBeamSSw[c.TalentLvlAttack()])

	// approximate the gadget hitbox as a long rectangle
	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: -7.5},
		3,
		15,
	)

	c.c2OnBeam()
	c.c6OnBeam()
	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)

	c.addDecode(17)

	if c.currFagio != stateDecoding {
		return
	}

	c.caBeamsRemaining -= 1
	if c.caBeamsRemaining <= 0 {
		c.exitCA()
		return
	}

	c.Core.Tasks.Add(func() { c.beamTicker(src) }, beamTickRate)
}

func (c *char) heatReduceTicker(src int) {
	// goes down by 2.5 per 0.5s, upgrade to 7.5 per 0.5s when off field
	// ignores hitlag

	if c.decodeReduceSrc != src {
		return
	}

	switch c.currFagio {
	case stateIdle:
	case stateOverdriveIdle:
	default:
		return
	}

	amt := 2.5
	if c.Core.Player.Active() != c.Index() {
		amt *= 3
	}
	c.reduceDecode(amt)
	c.Core.Tasks.Add(func() { c.heatReduceTicker(src) }, 0.5*60)
}

func (c *char) gainHeatTicker(src int) {
	// goes up by 3.2 per 0.2s
	// ignores hitlag
	// big beam adds 17 heat

	if c.caSrc != src {
		return
	}

	if c.currFagio != stateDecoding {
		return
	}

	c.addDecode(3.2)

	if c.currFagio != stateDecoding {
		return
	}

	var delay int = 0.2 * 60

	if c.caToBeams {
		c.caLength = c.Core.F - c.caSrc + delay
	}

	c.Core.Tasks.Add(func() { c.gainHeatTicker(src) }, delay)
}

func (c *char) enterOverdriveCharge() {
	// switch CA to slower ticks
	// start ticker for slower CA ticks
	if c.currFagio != stateOverdriveCharge {
		return
	}
	c.Core.Tasks.Add(func() { c.overdriveTicker(c.caSrc) }, overdriveFirstTick)
}

func (c *char) addDecode(amt float64) {
	amt *= c.c1DecoderGainMult()
	c.decode += amt

	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEvent(fmt.Sprintf("Gained %.2f Decode (%.2f)", amt, min(c.decode, maxDecode)), glog.LogCharacterEvent, c.Index())
	}

	if c.decode < maxDecode {
		return
	}
	c.decode = maxDecode

	reachedMaxDecode := func() {
		if c.decode < maxDecode {
			return
		}
		if c.currFagio == stateDecoding {
			c.currFagio = stateOverdriveCharge
			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Enter Overdrive", glog.LogCharacterEvent, c.Index())
			}
			if c.caToBeams {
				c.exitCA()
			} else {
				c.enterOverdriveCharge()
			}
		}

		if c.currFagio == stateIdle {
			c.currFagio = stateOverdriveIdle
		}
	}
	// delay on entering overdrive state
	c.Core.Tasks.Add(reachedMaxDecode, 19)
}

func (c *char) reduceDecode(amt float64) {
	if c.decode <= 0 {
		return
	}

	start := c.decode
	c.decode -= amt

	if c.decode < 50 {
		c.currFagio = stateIdle
	}

	if c.decode < 0 {
		c.decode = 0
	}

	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEvent(fmt.Sprintf("Reduced %.2f Decode (%.2f)", amt, c.decode), glog.LogCharacterEvent, c.Index())
	}

	actualAmt := start - c.decode
	c.a1OnDecreaseDecode(actualAmt)
}
