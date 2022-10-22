package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func WillCollide(p AttackPattern, t Target, key TargetKey) bool {
	//shape shouldn't be nil; panic here
	if p.Shape == nil {
		panic("unexpected nil shape")
	}

	//check if shape matches
	switch v := p.Shape.(type) {
	case *Circle:
		return t.Shape().IntersectCircle(*v)
	case *Rectangle:
		return t.Shape().IntersectRectangle(*v)
	case *SingleTarget:
		//only true if
		return v.Target == key
	default:
		return false
	}
}

func (c *Handler) AbsorbCheck(p AttackPattern, prio ...attributes.Element) attributes.Element {
	// check targets for collision first
	for _, e := range prio {
		for _, x := range c.enemies {
			t, ok := x.(TargetWithAura)
			if !ok {
				continue
			}
			if WillCollide(p, t, t.Key()) && t.AuraContains(e) {
				c.Log.NewEvent(
					"infusion check picked up "+e.String(),
					glog.LogElementEvent,
					-1,
				).
					Write("source", "enemy").
					Write("key", t.Key())
				return e
			}
		}
		for _, x := range c.gadgets {
			t, ok := x.(TargetWithAura)
			if !ok {
				continue
			}
			if WillCollide(p, t, t.Key()) && t.AuraContains(e) {
				c.Log.NewEvent(
					"infusion check picked up "+e.String(),
					glog.LogElementEvent,
					-1,
				).
					Write("source", "gadget").
					Write("key", t.Key())
				return e
			}
		}
		if t, ok := c.player.(TargetWithAura); ok {
			if WillCollide(p, t, 0) && t.AuraContains(e) {
				c.Log.NewEvent(
					"infusion check picked up "+e.String(),
					glog.LogElementEvent,
					-1,
				).
					Write("source", "player")
				return e
			}
		}

	}
	return attributes.NoElement
}
