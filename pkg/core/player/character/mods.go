package character

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type (
	//Status is basic mod for keeping track Status; usually affected by hitlag
	Status struct {
		modifier.Base
	}
	AttackMod struct {
		Amount AttackModFunc
		modifier.Base
	}
	AttackModFunc func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool)

	CooldownMod struct {
		Amount CooldownModFunc
		modifier.Base
	}
	CooldownModFunc func(a action.Action) float64

	DamageReductionMod struct {
		Amount DamageReductionModFunc
		modifier.Base
	}
	DamageReductionModFunc func() (float64, bool)

	HealBonusMod struct {
		Amount HealBonusModFunc
		modifier.Base
	}
	HealBonusModFunc func() (float64, bool)

	ReactBonusMod struct {
		Amount ReactBonusModFunc
		modifier.Base
	}
	ReactBonusModFunc func(combat.AttackInfo) (float64, bool)

	StatMod struct {
		AffectedStat attributes.Stat
		Extra        bool
		Amount       StatModFunc
		modifier.Base
	}
	StatModFunc func() ([]float64, bool)
)

// Add.

func (c *CharWrapper) AddStatus(key string, dur int, hitlag bool) {
	mod := Status{
		Base: modifier.Base{
			ModKey: key,
			Dur:    dur,
			Hitlag: hitlag,
		},
	}
	if mod.Dur < 0 {
		mod.ModExpiry = -1
	} else {
		mod.ModExpiry = *c.f + mod.Dur
	}
	overwrote, oldEvt := modifier.Add[modifier.Mod](&c.mods, &mod, *c.f)
	modifier.LogAdd("status", c.Index, &mod, c.log, overwrote, oldEvt)
}

func (c *CharWrapper) AddAttackMod(mod AttackMod) {
	mod.SetExpiry(*c.f)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&c.mods, &mod, *c.f)
	modifier.LogAdd("attack", c.Index, &mod, c.log, overwrote, oldEvt)
}

func (c *CharWrapper) AddCooldownMod(mod CooldownMod) {
	mod.SetExpiry(*c.f)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&c.mods, &mod, *c.f)
	modifier.LogAdd("cd", c.Index, &mod, c.log, overwrote, oldEvt)
}

func (c *CharWrapper) AddDamageReductionMod(mod DamageReductionMod) {
	mod.SetExpiry(*c.f)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&c.mods, &mod, *c.f)
	modifier.LogAdd("dr", c.Index, &mod, c.log, overwrote, oldEvt)
}

func (c *CharWrapper) AddHealBonusMod(mod HealBonusMod) {
	mod.SetExpiry(*c.f)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&c.mods, &mod, *c.f)
	modifier.LogAdd("heal bonus", c.Index, &mod, c.log, overwrote, oldEvt)
}

func (c *CharWrapper) AddReactBonusMod(mod ReactBonusMod) {
	mod.SetExpiry(*c.f)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&c.mods, &mod, *c.f)
	modifier.LogAdd("react bonus", c.Index, &mod, c.log, overwrote, oldEvt)
}

func (c *CharWrapper) AddStatMod(mod StatMod) {
	mod.SetExpiry(*c.f)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&c.mods, &mod, *c.f)
	modifier.LogAdd("stat", c.Index, &mod, c.log, overwrote, oldEvt)
}

// Delete.

func (c *CharWrapper) deleteMod(key string) {
	m := modifier.Delete(&c.mods, key)
	if m != nil {
		m.Event().SetEnded(*c.f)
	}
}
func (c *CharWrapper) DeleteStatus(key string)             { c.deleteMod(key) }
func (c *CharWrapper) DeleteAttackMod(key string)          { c.deleteMod(key) }
func (c *CharWrapper) DeleteCooldownMod(key string)        { c.deleteMod(key) }
func (c *CharWrapper) DeleteDamageReductionMod(key string) { c.deleteMod(key) }
func (c *CharWrapper) DeleteHealBonusMod(key string)       { c.deleteMod(key) }
func (c *CharWrapper) DeleteReactBonusMod(key string)      { c.deleteMod(key) }
func (c *CharWrapper) DeleteStatMod(key string)            { c.deleteMod(key) }

// Active.

func (c *CharWrapper) modIsActive(key string) bool {
	_, ok := modifier.FindCheckExpiry(&c.mods, key, *c.f)
	return ok
}
func (e *CharWrapper) StatusIsActive(key string) bool             { return e.modIsActive(key) }
func (e *CharWrapper) CooldownModIsActive(key string) bool        { return e.modIsActive(key) }
func (e *CharWrapper) DamageReductionModIsActive(key string) bool { return e.modIsActive(key) }
func (e *CharWrapper) HealBonusModIsActive(key string) bool       { return e.modIsActive(key) }
func (e *CharWrapper) ReactBonusModIsActive(key string) bool      { return e.modIsActive(key) }
func (e *CharWrapper) StatModIsActive(key string) bool            { return e.modIsActive(key) }

// Expiry.

func (c *CharWrapper) getModExpiry(key string) int {
	m := modifier.Find(&c.mods, key)
	if m != -1 {
		return c.mods[m].Expiry()
	}
	//must be 0 if doesn't exist. avoid using -1 b/c that's infinite
	return 0
}
func (c *CharWrapper) StatusExpiry(key string) int { return c.getModExpiry(key) }

// Duration.

func (c *CharWrapper) getModDuration(key string) int {
	m := modifier.Find(&c.mods, key)
	if m == -1 {
		return 0
	}
	if c.mods[m].Expiry() > *c.f {
		return c.mods[m].Expiry() - *c.f
	}
	return 0
}
func (c *CharWrapper) StatusDuration(key string) int { return c.getModDuration(key) }

// Extend.

// extendMod returns true if mod is active and is extended
func (c *CharWrapper) extendMod(key string, ext int) bool {
	m, active := modifier.FindCheckExpiry(&c.mods, key, *c.f)
	if m == -1 {
		return false
	}
	if !active {
		return false //nothing to extend is not active
	}
	//other wise add to expiry
	mod := c.mods[m]
	mod.Extend(mod.Key(), c.log, c.Index, ext)
	return true
}

func (c *CharWrapper) ExtendStatus(key string, ext int) bool { return c.extendMod(key, ext) }

// Amount.

func (c *CharWrapper) ApplyAttackMods(a *combat.AttackEvent, t combat.Target) []interface{} {
	//skip if this is reaction damage
	if a.Info.AttackTag >= attacks.AttackTagNoneStat {
		return nil
	}

	var sb strings.Builder
	var logDetails []interface{}

	if c.debug {
		logDetails = make([]interface{}, 0, len(c.mods))
	}

	n := 0
	for _, v := range c.mods {
		m, ok := v.(*AttackMod)
		if !ok {
			c.mods[n] = v
			n++
			continue
		}
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			amt, ok := m.Amount(a, t)
			if ok {
				for k, v := range amt {
					a.Snapshot.Stats[k] += v
				}
			}
			c.mods[n] = v
			n++
			if c.debug {
				modStatus := make([]string, 0)
				if ok {
					sb.WriteString(m.Key())
					modStatus = append(
						modStatus,
						"status: added",
						"expiry_frame: "+strconv.Itoa(m.Expiry()),
					)
					modStatus = append(
						modStatus,
						attributes.PrettyPrintStatsSlice(amt)...,
					)
					logDetails = append(logDetails, sb.String(), modStatus)
					sb.Reset()
				} else {
					sb.WriteString(m.Key())
					modStatus = append(
						modStatus,
						"status: rejected",
						"reason: conditions not met",
					)
					logDetails = append(logDetails, sb.String(), modStatus)
					sb.Reset()
				}
			}
		}
	}
	c.mods = c.mods[:n]
	return logDetails
}

func (c *CharWrapper) CDReduction(a action.Action, dur int) int {
	var cd float64 = 1
	n := 0
	for _, v := range c.mods {
		m, ok := v.(*CooldownMod)
		if !ok {
			c.mods[n] = v
			n++
			continue
		}
		//if not expired
		if m.Expiry() == -1 || m.Expiry() > *c.f {
			amt := m.Amount(a)
			c.log.NewEvent("applying cooldown modifier", glog.LogActionEvent, c.Index).
				Write("key", m.Key()).
				Write("modifier", amt).
				Write("expiry", m.Expiry())
			cd += amt
			c.mods[n] = v
			n++
		}
	}
	c.mods = c.mods[:n]

	return int(float64(dur) * cd)
}

func (c *CharWrapper) DamageReduction(char int) (amt float64) {
	n := 0
	for _, v := range c.mods {
		m, ok := v.(*DamageReductionMod)
		if !ok {
			c.mods[n] = v
			n++
			continue
		}
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			a, done := m.Amount()
			amt += a
			if !done {
				c.mods[n] = v
				n++
			}
		}
	}
	c.mods = c.mods[:n]
	return amt
}

func (c *CharWrapper) HealBonus() (amt float64) {
	n := 0
	for _, v := range c.mods {
		m, ok := v.(*HealBonusMod)
		if !ok {
			c.mods[n] = v
			n++
			continue
		}
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			a, done := m.Amount()
			amt += a
			if !done {
				c.mods[n] = v
				n++
			}
		}
	}
	c.mods = c.mods[:n]
	return amt
}

// TODO: consider merging this with just attack mods? reaction bonus should
// maybe just be it's own stat instead of being a separate mod really
func (c *CharWrapper) ReactBonus(atk combat.AttackInfo) (amt float64) {
	n := 0
	for _, v := range c.mods {
		m, ok := v.(*ReactBonusMod)
		if !ok {
			c.mods[n] = v
			n++
			continue
		}
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			a, done := m.Amount(atk)
			amt += a
			if !done {
				c.mods[n] = v
				n++
			}
		}
	}
	c.mods = c.mods[:n]
	return amt
}
