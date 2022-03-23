package combat

import (
	"log"

	"github.com/genshinsim/gcsim/pkg/core"
)

type Ctrl struct {
	core *core.Core
}

func NewCtrl(c *core.Core) *Ctrl {
	return &Ctrl{
		core: c,
	}
}

func (c *Ctrl) QueueAttackWithSnap(a core.AttackInfo, s core.Snapshot, p core.AttackPattern, dmgDelay int, callbacks ...core.AttackCBFunc) {
	if dmgDelay < 0 {
		panic("dmgDelay cannot be less than 0")
	}
	ae := core.AttackEvent{
		Info:        a,
		Pattern:     p,
		Snapshot:    s,
		SourceFrame: c.core.F,
	}
	//add callbacks only if not nil
	for _, f := range callbacks {
		if f != nil {
			ae.Callbacks = append(ae.Callbacks, f)
		}
	}
	c.queueDmg(&ae, dmgDelay)
}

func (c *Ctrl) QueueAttackEvent(ae *core.AttackEvent, dmgDelay int) {
	c.queueDmg(ae, dmgDelay)
}

func (c *Ctrl) QueueAttack(a core.AttackInfo, p core.AttackPattern, snapshotDelay int, dmgDelay int, callbacks ...core.AttackCBFunc) {
	//panic if dmgDelay > snapshotDelay; this should not happen. if it happens then there's something wrong with the
	//character's code
	if dmgDelay < snapshotDelay {
		panic("dmgDelay cannot be less than snapshotDelay")
	}
	if dmgDelay < 0 {
		panic("dmgDelay cannot be less than 0")
	}
	//create attackevent
	ae := core.AttackEvent{
		Info:        a,
		Pattern:     p,
		SourceFrame: c.core.F,
	}
	//add callbacks only if not nil
	for _, f := range callbacks {
		if f != nil {
			ae.Callbacks = append(ae.Callbacks, f)
		}
	}

	switch {
	case snapshotDelay < 0:
		//snapshotDelay < 0 means we don't need a snapshot; optimization for reaction
		//damage essentially
		c.queueDmg(&ae, dmgDelay)
	case snapshotDelay == 0:
		c.generateSnapshot(&ae)
		c.queueDmg(&ae, dmgDelay)
	default:
		//use add task ctrl to queue; no need to track here
		c.core.Tasks.Add(func() {
			c.generateSnapshot(&ae)
			c.queueDmg(&ae, dmgDelay-snapshotDelay)
		}, snapshotDelay)
	}

}

func (c *Ctrl) generateSnapshot(a *core.AttackEvent) {
	a.Snapshot = c.core.Chars[a.Info.ActorIndex].Snapshot(&a.Info)
}

func (c *Ctrl) queueDmg(a *core.AttackEvent, delay int) {
	if delay == 0 {
		c.ApplyDamage(a)
		return
	}
	c.core.Tasks.Add(func() {
		c.ApplyDamage(a)
	}, delay)
}

func willAttackLand(a *core.AttackEvent, t core.Target, index int) (bool, string) {
	//shape shouldn't be nil; panic here
	if a.Pattern.Shape == nil {
		panic("unexpected nil shape")
	}
	//shape can't be nil now, check if type matches
	if !a.Pattern.Targets[t.Type()] {
		return false, "wrong type"
	}
	//skip if self harm is false and dmg src == i
	if !a.Pattern.SelfHarm && a.Info.DamageSrc == index {
		return false, "no self harm"
	}

	//check if shape matches
	switch v := a.Pattern.Shape.(type) {
	case *core.Circle:
		return t.Shape().IntersectCircle(*v), "intersect circle"
	case *core.Rectangle:
		return t.Shape().IntersectRectangle(*v), "intersect rectangle"
	case *core.SingleTarget:
		//only true if
		return v.Target == index, "target"
	default:
		return false, "unknown shape"
	}
}

func (c *Ctrl) ApplyDamage(a *core.AttackEvent) float64 {
	// died := false
	var total float64
	for i, t := range c.core.Targets {
		//skip nil targets; we don't want to reindex...
		if t == nil {
			continue
		}

		willHit, reason := willAttackLand(a, t, i)
		if !willHit {
			// Move target logs into the "Sim" event log to avoid cluttering main display for stuff like Guoba
			// And obvious things like "Fischl A4 is single target so it didn't hit targets 2-4"
			// TODO: Maybe want to add a separate set of log events for this?
			//don't log this for target 0
			if c.core.Flags.LogDebug && i > 0 {
				c.core.Log.NewEvent(
					"skipped "+a.Info.Abil+" "+reason,
					core.LogSimEvent,
					a.Info.ActorIndex,
					"attack_tag", a.Info.AttackTag,
					"applied_ele", a.Info.Element,
					"dur", a.Info.Durability,
					"target", i,
					"shape", a.Pattern.Shape.String(),
					// "type", fmt.Sprintf("%T", a.Pattern.Shape),
				)
			}
			continue
		}

		//make a copy first
		cpy := *a

		//at this point attack will land
		c.core.Events.Emit(core.OnAttackWillLand, t, &cpy)

		//check to make sure it's not cancelled for w/e reason
		if a.Cancelled {
			continue
		}

		var evt core.LogEvent = nil
		var amp string
		var dmg float64
		var crit bool

		if c.core.Flags.LogDebug {
			evt = c.core.Log.NewEvent(
				cpy.Info.Abil,
				core.LogDamageEvent,
				cpy.Info.ActorIndex,
				"target", i,
				"attack-tag", cpy.Info.AttackTag,
				"ele", cpy.Info.Element.String(),
				"damage", &dmg,
				"crit", &crit,
				"amp", &amp,
				"abil", cpy.Info.Abil,
				"source_frame", cpy.SourceFrame,
			)
			evt.Write(cpy.Snapshot.Logs...)
		}

		if !cpy.Info.SourceIsSim {
			if cpy.Info.ActorIndex < 0 {
				log.Println(cpy)
			}
			char := c.core.Chars[cpy.Info.ActorIndex]
			preDmgModDebug := char.PreDamageSnapshotAdjust(&cpy, t)
			if c.core.Flags.LogDebug {
				evt.Write("pre_damage_mods", preDmgModDebug)
			}
		}

		dmg, crit = t.Attack(&cpy, evt)
		total += dmg

		c.core.Events.Emit(core.OnDamage, t, &cpy, dmg, crit)

		//callbacks
		cb := core.AttackCB{
			Target:      t,
			AttackEvent: &cpy,
			Damage:      dmg,
			IsCrit:      crit,
		}
		for _, f := range cpy.Callbacks {
			f(cb)
		}

		//check if target is dead; skip this for i = 0 since we don't want to
		//delete the player by accident
		if c.core.Flags.DamageMode && t.HP() <= 0 {
			log.Println("died")
			// died = true
			t.Kill()
			c.core.Events.Emit(core.OnTargetDied, t, cpy)
			//this should be ok for stuff like guoba since they won't take damage
			c.core.Targets[i] = nil
			// log.Println("target died", i, dmg)
		}

		// this works because string in golang is a slice underneath, so the &amp points to the slice info
		// that's why when the underlying string in amp changes (has to be reallocated) the pointer doesn't
		// change since it's just pointing to the slice "header"
		if cpy.Info.Amped {
			amp = string(cpy.Info.AmpType)
		}

	}
	// if died {
	// 	c.core.ReindexTargets()
	// }
	c.core.TotalDamage += total
	return total
}

func (c *Ctrl) TargetHasResMod(key string, param int) bool {
	if param >= len(c.core.Targets) {
		return false
	}
	return c.core.Targets[param].HasResMod(key)
}
func (c *Ctrl) TargetHasDefMod(key string, param int) bool {
	if param >= len(c.core.Targets) {
		return false
	}
	return c.core.Targets[param].HasDefMod(key)
}

func (c *Ctrl) TargetHasElement(ele core.EleType, param int) bool {
	if param >= len(c.core.Targets) {
		return false
	}
	return c.core.Targets[param].AuraContains(ele)
}
