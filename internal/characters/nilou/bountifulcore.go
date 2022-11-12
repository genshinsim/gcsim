package nilou

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type BountifulCore struct {
	srcFrame int
	*gadget.Gadget
}

func newBountifulCore(c *core.Core, x float64, y float64, a *combat.AttackEvent) *BountifulCore {
	b := &BountifulCore{
		srcFrame: c.F,
	}

	b.Gadget = gadget.New(c, core.Coord{X: x, Y: y, R: 0.2}, combat.GadgetTypDendroCore)
	b.Gadget.Duration = 0.4 * 60

	char := b.Core.Player.ByIndex(a.Info.ActorIndex)
	explode := func() {
		ai, snap := reactable.NewBloomAttack(char, b)
		ap := combat.NewCircleHit(b.Gadget, 6.5)
		c.QueueAttackWithSnap(ai, snap, ap, 1)

		//self damage
		ai.Abil += " (self damage)"
		ai.FlatDmg = 0.05 * ai.FlatDmg
		ap.SkipTargets[combat.TargettablePlayer] = false
		ap.SkipTargets[combat.TargettableEnemy] = true
		ap.SkipTargets[combat.TargettableGadget] = true
		c.QueueAttackWithSnap(ai, snap, ap, 1)
	}
	b.Gadget.OnExpiry = explode
	b.Gadget.OnKill = explode

	return b
}

func (b *BountifulCore) Tick() {
	//this is needed since gadget tick
	b.Gadget.Tick()
}

func (b *BountifulCore) HandleAttack(atk *combat.AttackEvent) float64 {
	b.Core.Events.Emit(event.OnGadgetHit, b, atk)
	return 0
}
func (b *BountifulCore) Attack(*combat.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (b *BountifulCore) ApplyDamage(*combat.AttackEvent, float64)               {}
