package crystallize

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type Shard struct {
	*gadget.Gadget
	// earliest that a shard can be picked up after spawn
	EarliestPickup int
	// captures the shield because em snapshots
	Shield *Shield
	// for logging purposes
	src    int
	expiry int
}

func NewShard(c *core.Core, shp info.Shape, shd *Shield) *Shard {
	cs := &Shard{}

	circ, ok := shp.(*info.Circle)
	if !ok {
		panic("rectangle target hurtbox is not supported for crystallize shard spawning")
	}

	// for simplicity, crystallize shards spawn randomly at radius + 0.5
	r := circ.Radius() + 0.5
	// radius 2 is ok
	cs.Gadget = gadget.New(c, info.CalcRandomPointFromCenter(circ.Pos(), r, r, c.Rand), 2, info.GadgetTypCrystallizeShard)

	// shard lasts for 15s from shard spawn
	cs.Duration = 15 * 60
	// earliest shard pickup is 54f from crystallize text, so 31f from shard spawn
	cs.EarliestPickup = c.F + 31
	cs.Shield = shd
	cs.src = c.F
	cs.expiry = c.F + cs.Duration

	return cs
}

func (cs *Shard) AddShieldKillShard() bool {
	// don't pick up if shard is not available for pick up yet
	if cs.Core.F < cs.EarliestPickup {
		cs.Core.Log.NewEvent(
			fmt.Sprintf("%v crystallize shard could not be picked up", cs.Shield.Ele),
			glog.LogElementEvent,
			cs.Core.Player.Active(),
		).
			Write("src", cs.src).
			Write("expiry", cs.expiry).
			Write("earliest_pickup", cs.EarliestPickup)
		return false
	}
	cs.Core.Log.NewEvent(
		fmt.Sprintf("%v crystallize shard picked up", cs.Shield.Ele),
		glog.LogElementEvent,
		cs.Core.Player.Active(),
	).
		Write("src", cs.src).
		Write("expiry", cs.expiry).
		Write("earliest_pickup", cs.EarliestPickup)
	// add shield
	cs.Shield.Expires = cs.Core.F + 15.1*60 // shield lasts for 15.1s from shard pickup
	cs.Core.Player.Shields.Add(cs.Shield)
	// kill self
	cs.Kill()
	return true
}

func (cs *Shard) HandleAttack(atk *info.AttackEvent) float64 {
	cs.Core.Events.Emit(event.OnGadgetHit, cs, atk)
	return 0
}
func (cs *Shard) Attack(*info.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (cs *Shard) SetDirection(trg info.Point)                          {}
func (cs *Shard) SetDirectionToClosestEnemy()                          {}
func (cs *Shard) CalcTempDirection(trg info.Point) info.Point {
	return info.DefaultDirection()
}

func (cs *Shard) Src() int {
	return cs.src
}

func (cs *Shard) Expiry() int {
	return cs.expiry
}
