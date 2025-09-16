package core

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (c *Core) QueueAttackWithSnap(
	a info.AttackInfo,
	s info.Snapshot,
	p info.AttackPattern,
	dmgDelay int,
	callbacks ...info.AttackCBFunc,
) {
	if dmgDelay < 0 {
		panic("dmgDelay cannot be less than 0")
	}
	ae := info.AttackEvent{
		Info:        a,
		Pattern:     p,
		Snapshot:    s,
		SourceFrame: c.F,
	}
	// add callbacks only if not nil
	for _, f := range callbacks {
		if f != nil {
			ae.Callbacks = append(ae.Callbacks, f)
		}
	}
	c.queueDmg(&ae, dmgDelay)
}

func (c *Core) QueueAttackEvent(ae *info.AttackEvent, dmgDelay int) {
	c.queueDmg(ae, dmgDelay)
}

func (c *Core) QueueAttack(
	a info.AttackInfo,
	p info.AttackPattern,
	snapshotDelay int,
	dmgDelay int,
	callbacks ...info.AttackCBFunc,
) {
	// panic if dmgDelay < snapshotDelay; this should not happen. if it happens then there's something wrong with the
	// character's code
	if dmgDelay < snapshotDelay {
		panic("dmgDelay cannot be less than snapshotDelay")
	}
	if dmgDelay < 0 {
		panic("dmgDelay cannot be less than 0")
	}
	// create attackevent
	ae := info.AttackEvent{
		Info:        a,
		Pattern:     p,
		SourceFrame: c.F,
	}
	// add callbacks only if not nil
	for _, f := range callbacks {
		if f != nil {
			ae.Callbacks = append(ae.Callbacks, f)
		}
	}

	switch {
	case snapshotDelay < 0:
		// snapshotDelay < 0 means we don't need a snapshot; optimization for reaction
		// damage essentially
		c.queueDmg(&ae, dmgDelay)
	case snapshotDelay == 0:
		c.generateSnapshot(&ae)
		c.queueDmg(&ae, dmgDelay)
	default:
		// use add task ctrl to queue; no need to track here
		c.Tasks.Add(func() {
			c.generateSnapshot(&ae)
			c.queueDmg(&ae, dmgDelay-snapshotDelay)
		}, snapshotDelay)
	}
}

// This code here should probably be handled in player not core
// since it's a convenience function wrapped around queuedamage
//
// does it make sense for core to have any knowledge of teams? probably not??
func (c *Core) generateSnapshot(a *info.AttackEvent) {
	a.Snapshot = c.Player.ByIndex(a.Info.ActorIndex).Snapshot(&a.Info)
}

func (c *Core) queueDmg(a *info.AttackEvent, delay int) {
	if delay == 0 {
		c.Combat.ApplyAttack(a)
		return
	}
	c.Tasks.Add(func() {
		c.Combat.ApplyAttack(a)
	}, delay)
}
