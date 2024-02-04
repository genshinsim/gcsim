package lyney

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type GrinMalkinHat struct {
	*gadget.Gadget
	char                *char
	pos                 geometry.Point
	pyrotechnicAI       combat.AttackInfo
	pyrotechnicSnapshot combat.Snapshot
	hpDrained           bool
	a1CB                combat.AttackCBFunc
}

func (c *char) newGrinMalkinHat(pos geometry.Point, hpDrained bool, duration int) *GrinMalkinHat {
	g := &GrinMalkinHat{}

	g.pos = pos

	// TODO: double check estimation of hitbox
	g.Gadget = gadget.New(c.Core, g.pos, 1, combat.GadgetTypGrinMalkinHat)
	g.char = c

	g.Duration = duration
	g.char.AddStatus(grinMalkinHatKey, g.Duration, false)

	g.pyrotechnicAI = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Pyrotechnic Strike",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagLyneyEndBoom,
		ICDGroup:   attacks.ICDGroupLyneyExtra,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       propPyrotechnic[c.TalentLvlAttack()],
	}
	g.pyrotechnicSnapshot = g.char.Snapshot(&g.pyrotechnicAI)
	g.char.addA1(&g.pyrotechnicAI, hpDrained)
	g.hpDrained = hpDrained
	g.a1CB = g.char.makeA1CB(hpDrained)

	g.OnExpiry = g.skillPyrotechnic("expiry")
	g.OnKill = g.skillPyrotechnic("kill")

	g.Core.Log.NewEvent("Lyney Grin-Malkin Hat added", glog.LogCharacterEvent, c.Index).Write("src", g.Src()).Write("hp_drained", g.hpDrained)

	return g
}

func (g *GrinMalkinHat) HandleAttack(atk *combat.AttackEvent) float64 {
	g.Core.Events.Emit(event.OnGadgetHit, g, atk)

	// TODO: gadget taking damage is not implemented

	return 0
}

func (g *GrinMalkinHat) skillPyrotechnic(reason string) func() {
	return func() {
		// needed for amos and slingshot to work correctly
		g.pyrotechnicSnapshot.SourceFrame = g.Core.F
		// TODO: snapshot timing
		g.Core.QueueAttackWithSnap(
			g.pyrotechnicAI,
			g.pyrotechnicSnapshot,
			combat.NewCircleHit(
				g.Core.Combat.Player(),
				g.Core.Combat.PrimaryTarget(),
				nil,
				1,
			),
			g.char.pyrotechnicTravel,
			g.a1CB,
			g.char.makeC4CB(),
		)
		g.updateHats(reason)
	}
}

func (g *GrinMalkinHat) skillExplode() {
	g.pyrotechnicAI.ICDTag = attacks.ICDTagLyneyEndBoomEnhanced
	g.pyrotechnicAI.StrikeType = attacks.StrikeTypeBlunt
	g.pyrotechnicAI.PoiseDMG = 90

	// needed for amos and slingshot to work correctly
	g.pyrotechnicSnapshot.SourceFrame = g.Core.F
	// TODO: snapshot timing
	g.Core.QueueAttackWithSnap(
		g.pyrotechnicAI,
		g.pyrotechnicSnapshot,
		combat.NewCircleHitOnTarget(g.pos, nil, 3.5),
		skillExplode,
		g.a1CB,
		g.char.makeC4CB(),
	)

	g.OnKill = nil // prevent additional pyrotechnic attack
	g.Kill()

	g.updateHats("skill explode")
}

func (g *GrinMalkinHat) updateHats(removeReason string) {
	for i := 0; i < len(g.char.hats); i++ {
		if g.char.hats[i] == g {
			g.char.hats = append(g.char.hats[:i], g.char.hats[i+1:]...)
			g.Core.Log.NewEvent("Lyney Grin-Malkin Hat removed", glog.LogCharacterEvent, g.char.Index).Write("src", g.Src()).Write("hp_drained", g.hpDrained).Write("remove_reason", removeReason)
		}
	}
}

func (g *GrinMalkinHat) SetDirection(trg geometry.Point) {}
func (g *GrinMalkinHat) SetDirectionToClosestEnemy()     {}
func (g *GrinMalkinHat) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}
