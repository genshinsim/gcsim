package nilou

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type BountifulCore struct {
	srcFrame int
	*gadget.Gadget
}

// TODO: "and they have larger AoEs"
func newBountifulCore(c *core.Core, x float64, y float64, a *combat.AttackEvent) *BountifulCore {
	b := &BountifulCore{
		srcFrame: c.F,
	}

	b.Gadget = gadget.New(c, core.Coord{X: x, Y: y, R: 0.2}, combat.GadgetTypBountifulCore)
	b.Gadget.Duration = 0.2 * 60

	char := b.Core.Player.ByIndex(a.Info.ActorIndex)
	explode := func() {
		// c.Combat.Log.NewEvent("bountiful core boom", glog.LogCharacterEvent, char.Index)
		ai := reactable.NewBloomAttack(char, b)
		c.QueueAttack(ai, combat.NewCircleHit(b, 5, false, combat.TargettableEnemy), -1, 1)

		//self damage
		ai.Abil += " (self damage)"
		ai.FlatDmg = 0.05 * ai.FlatDmg
		c.QueueAttack(ai, combat.NewCircleHit(b.Gadget, 5, true, combat.TargettablePlayer), -1, 1)
	}
	//TODO: should bloom do damage if it blows up due to limit reached?
	b.Gadget.OnExpiry = explode
	b.Gadget.OnKill = explode

	return b
}

func (b *BountifulCore) Tick() {
	//this is needed since gadget tick
	b.Gadget.Tick()
}
func (b *BountifulCore) Attack(*combat.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (b *BountifulCore) ApplyDamage(*combat.AttackEvent, float64)               {}
