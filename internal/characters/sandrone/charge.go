package sandrone

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

// TODO: frames estimated from in-game footage (60 FPS), +-2-4 frames.
const (
	chargeDefaultHold = 240
	chargeMinHold     = 60
	chargeMaxHold     = 900
	chargeEndLag      = 20

	fagioSummonFrame = 24
	sweepFirstHit    = 70
	sweepInterval    = 24
	beamFirstHit     = 120
	beamInterval     = 58

	beamAbil = "Condensed Beam"
)

// Summons Fagio and switches it to Decoding mode: Unleashes a sweeping fire attack on enemies in
// front, as well as periodically fires condensed beams, dealing AoE Cryo DMG.
//
// Use the "hold" param (frames, default 240) to control how long the Charged Attack is held.
func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	c.endHovermech()

	hold, ok := p["hold"]
	if !ok {
		hold = chargeDefaultHold
	}
	hold = max(chargeMinHold, min(chargeMaxHold, hold))

	src := c.Core.F
	c.chargeSrc = src

	// Fagio cannot switch to Decoding mode while in Power Overdrive; only the sweeping fire happens
	c.QueueCharTask(func() {
		if c.chargeSrc != src || c.fagioMode == fagioOverdrive {
			return
		}
		c.beamCount = 0
		c.setFagioMode(fagioDecoding)
	}, fagioSummonFrame)

	c.QueueCharTask(func() {
		if c.chargeSrc == src && c.fagioMode == fagioDecoding {
			c.setFagioMode(fagioNone)
		}
	}, hold)

	for t := sweepFirstHit; t < hold; t += sweepInterval {
		c.QueueCharTask(func() {
			if c.chargeSrc != src {
				return
			}
			c.sweepingFire()
		}, t)
	}

	for t := beamFirstHit; t < hold; t += beamInterval {
		c.QueueCharTask(func() {
			if c.chargeSrc != src || c.fagioMode != fagioDecoding {
				return
			}
			c.condensedBeam()
		}, t)
	}

	chargeFrames := frames.InitAbilSlice(hold + chargeEndLag)
	chargeFrames[action.ActionDash] = hold
	chargeFrames[action.ActionSwap] = hold

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   hold,
		State:           action.ChargeAttackState,
		OnRemoved: func(next action.AnimationState) {
			// releasing / interrupting the hold ends Decoding mode
			if c.chargeSrc == src && c.fagioMode == fagioDecoding {
				c.setFagioMode(fagioNone)
			}
			if c.chargeSrc == src {
				c.chargeSrc = -1
			}
		},
	}, nil
}

func (c *char) sweepingFire() {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Sweeping Fire",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       chargeSweep[c.TalentLvlAttack()],
	}
	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1}, 3, 6)
	c.Core.QueueAttack(ai, ap, 0, 0)
}

// Periodically fires condensed beams, dealing AoE Cryo DMG.
// Radiance: Stellar Glimmer: The condensed beams will instead deal AoE Cryo DMG that is considered
// the corresponding Stellar Glimmer reaction DMG.
func (c *char) condensedBeam() {
	c.beamCount++

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       beamAbil,
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       chargeBeam[c.TalentLvlAttack()],
	}
	lvl := c.TalentLvlAttack()
	c.toStellar(&ai, chargeBeamConduct[lvl], chargeBeamSwirl[lvl])

	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1}, 4, 10)
	c.Core.QueueAttack(ai, ap, 0, 0)

	c.c6ClusterBeam(ap)

	// firing condensed beams also further increases Decoding Power
	c.addDecodingPower(decodeBeamGain * c.decodeMult())
}
