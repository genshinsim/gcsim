package sandrone

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type fagioMode int

const (
	fagioNone fagioMode = iota
	fagioDecoding
	fagioOverdrive
)

func (m fagioMode) String() string {
	switch m {
	case fagioDecoding:
		return "decoding"
	case fagioOverdrive:
		return "overdrive"
	default:
		return "none"
	}
}

// Decoding Power rates are not published; they are fitted to footage:
// starting from 0, continuous holding reaches 100 right at the 3rd condensed beam (~236f).
// TODO: rates estimated from in-game footage.
const (
	decodingMax       = 100.0
	overdriveExit     = 50.0
	decodeRate        = 18.0 // per second while in Decoding mode
	decodeBeamGain    = 10.0 // per condensed beam fired
	decayRate         = 4.0  // per second while not in Decoding mode, Sandrone on-field
	decayOffFieldMult = 3.0
	// per second while on the Hovermech ("swiftly decrease"); fast enough that the usual
	// CA -> E -> CA loop starts every Charged Attack from a reset gauge
	hovermechDecay = 100.0

	fagioTickInterval = 6 // frames

	overdriveFirstShot = 60
	overdriveInterval  = 258 // ~4.3s between Power Overdrive shots
)

func (c *char) decodeMult() float64 {
	if c.Base.Cons >= 1 {
		return 0.5
	}
	return 1
}

func (c *char) setFagioMode(m fagioMode) {
	if c.fagioMode == m {
		return
	}
	c.Core.Log.NewEvent("fagio mode changed", glog.LogCharacterEvent, c.Index()).
		Write("from", c.fagioMode.String()).
		Write("to", m.String()).
		Write("decoding_power", c.decodingPower)
	c.fagioMode = m
	if m == fagioOverdrive {
		c.startOverdrive()
	}
	c.ensureFagioTicker()
}

func (c *char) ensureFagioTicker() {
	if c.tickerActive {
		return
	}
	c.tickerActive = true
	c.Core.Tasks.Add(c.fagioTick, fagioTickInterval)
}

// fagioTick runs while Fagio is decoding or still has Decoding Power to lose. It is not affected by
// Sandrone's hitlag since Fagio keeps working off-field.
func (c *char) fagioTick() {
	const dt = float64(fagioTickInterval) / 60

	if c.fagioMode == fagioDecoding {
		c.addDecodingPower(decodeRate * dt * c.decodeMult())
	} else {
		rate := decayRate
		switch {
		case c.StatusIsActive(hovermechKey):
			rate = hovermechDecay
		case c.Core.Player.Active() != c.Index():
			rate *= decayOffFieldMult
		}
		c.loseDecodingPower(rate * dt)
	}

	if c.fagioMode != fagioDecoding && c.decodingPower <= 0 {
		c.tickerActive = false
		return
	}
	c.Core.Tasks.Add(c.fagioTick, fagioTickInterval)
}

func (c *char) addDecodingPower(amt float64) {
	if c.fagioMode != fagioDecoding {
		return
	}
	c.decodingPower = min(decodingMax, c.decodingPower+amt)
	if c.decodingPower >= decodingMax {
		c.setFagioMode(fagioOverdrive)
	}
}

func (c *char) loseDecodingPower(amt float64) {
	lost := min(amt, c.decodingPower)
	if lost <= 0 {
		return
	}
	c.decodingPower -= lost
	c.lostPower += lost
	for c.lostPower >= 10 {
		c.lostPower -= 10
		c.a1AddTactics()
	}
	if c.fagioMode == fagioOverdrive && c.decodingPower < overdriveExit {
		c.setFagioMode(fagioNone)
	}
}

// When Fagio's Decoding Power reaches the maximum of 100 points, it switches to Power Overdrive
// mode, in which it fires at opponents with longer intervals, dealing Cryo DMG.
func (c *char) startOverdrive() {
	src := c.Core.F
	c.overdriveSrc = src
	c.Core.Tasks.Add(func() { c.overdriveShot(src) }, overdriveFirstShot)
}

func (c *char) overdriveShot(src int) {
	if c.overdriveSrc != src || c.fagioMode != fagioOverdrive {
		return
	}
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Fagio (Power Overdrive)",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       overdrive[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3), 0, 0)
	c.Core.Tasks.Add(func() { c.overdriveShot(src) }, overdriveInterval)
}
