package enemy

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Add.
func (e *Enemy) AddStatus(key string, dur int, hitlag bool) {
	mod := info.Status{
		Base: modifier.Base{
			ModKey: key,
			Dur:    dur,
			Hitlag: hitlag,
		},
	}
	if mod.Dur < 0 {
		mod.ModExpiry = -1
	} else {
		mod.ModExpiry = e.Core.F + mod.Dur
	}
	addMod(e, "status", key, &mod)
}

// Add a ResistMod
//
//	e.AddResistMod(info.ResistMod{
//		Base:  modifier.NewBaseWithHitlag("debuff", 10*60),
//		Ele:   attributes.Cryo,
//		Value: -0.15,
//	})
func (e *Enemy) AddResistMod(mod info.ResistMod) {
	mod.SetExpiry(e.Core.F)
	addMod(e, "enemy", mod.ModKey, &mod)
}

// Add a DefMod
//
//	e.AddDefMod(info.DefMod{
//		Base:  modifier.NewBaseWithHitlag("debuff", 10*60),
//		Value: -0.3,
//	})
func (e *Enemy) AddDefMod(mod info.DefMod) {
	mod.SetExpiry(e.Core.F)
	addMod(e, "enemy", mod.ModKey, &mod)
}

// addMod adds a copy of mod, whose key is key, the way modifier.Add adds a mod. When a mod of
// the same key and type is already there, it is overwritten in place instead of being replaced
// by a new allocation. For enemies that is the same thing: no pointer to an enemy mod is kept
// outside e.mods, and nothing that ranges over e.mods can add a mod.
func addMod[T any, PT interface {
	*T
	modifier.Mod
}](e *Enemy, prefix, key string, mod *T) {
	ind := modifier.Find(&e.mods, key)
	if ind != -1 {
		if old, ok := e.mods[ind].(PT); ok {
			overwrote, oldEvt := modifier.Replaced(old, e.Core.F)
			*old = *mod
			modifier.LogAdd(prefix, -1, old, e.Core.Log, overwrote, oldEvt)
			return
		}
	}
	m := PT(new(T))
	*m = *mod
	overwrote := false
	var oldEvt glog.Event
	if ind == -1 {
		e.mods = append(e.mods, m)
	} else {
		overwrote, oldEvt = modifier.Replaced(e.mods[ind], e.Core.F)
		e.mods[ind] = m
	}
	modifier.LogAdd(prefix, -1, m, e.Core.Log, overwrote, oldEvt)
}

// Delete.

func (e *Enemy) deleteMod(key string) {
	m := modifier.Delete(&e.mods, key)
	if m != nil && (m.Expiry() > e.Core.F || m.Expiry() == -1) {
		m.Event().SetEnded(e.Core.F)
	}
}

func (e *Enemy) DeleteStatus(key string)    { e.deleteMod(key) }
func (e *Enemy) DeleteResistMod(key string) { e.deleteMod(key) }
func (e *Enemy) DeleteDefMod(key string)    { e.deleteMod(key) }

// Active.
func (e *Enemy) modIsActive(key string) bool {
	_, ok := modifier.FindCheckExpiry(&e.mods, key, e.Core.F)
	return ok
}
func (e *Enemy) StatusIsActive(key string) bool    { return e.modIsActive(key) }
func (e *Enemy) ResistModIsActive(key string) bool { return e.modIsActive(key) }
func (e *Enemy) DefModIsActive(key string) bool    { return e.modIsActive(key) }

// Expiry

func (e *Enemy) getModExpiry(key string) int {
	m := modifier.Find(&e.mods, key)
	if m != -1 {
		return e.mods[m].Expiry()
	}
	// must be 0 if doesn't exist. avoid using -1 b/c that's infinite
	return 0
}
func (e *Enemy) StatusExpiry(key string) int { return e.getModExpiry(key) }

// Amount.

// TODO: this needs to purge if done?
func (e *Enemy) resist(ai *info.AttackInfo, evt glog.Event) float64 {
	var logDetails []any
	var sb strings.Builder

	if e.Core.Flags.LogDebug {
		logDetails = make([]any, 0, 5*len(e.mods))
	}

	r := e.resists[ai.Element]
	for _, v := range e.mods {
		m, ok := v.(*info.ResistMod)
		if !ok {
			continue
		}
		if m.Expiry() > e.Core.F && m.Ele == ai.Element {
			if e.Core.Flags.LogDebug {
				sb.WriteString(m.Key())
				logDetails = append(logDetails, sb.String(), []string{
					"status: added",
					"expiry_frame: " + strconv.Itoa(m.Expiry()),
					"ele: " + m.Ele.String(),
					"amount: " + strconv.FormatFloat(m.Value, 'f', -1, 64),
				})
				sb.Reset()
			}
			r += m.Value
		}
	}

	// No need to output if resist was not modified
	if e.Core.Flags.LogDebug && len(logDetails) > 1 {
		evt.Write("resist_mods", logDetails)
	}

	return r
}

func (e *Enemy) defAdj(evt glog.Event) float64 {
	var logDetails []any
	var sb strings.Builder

	if e.Core.Flags.LogDebug {
		logDetails = make([]any, 0, 3*len(e.mods))
	}

	var r float64
	for _, v := range e.mods {
		m, ok := v.(*info.DefMod)
		if !ok {
			continue
		}
		if m.Expiry() > e.Core.F {
			if e.Core.Flags.LogDebug {
				sb.WriteString(m.Key())
				logDetails = append(logDetails, sb.String(), []string{
					"status: added",
					"expiry_frame: " + strconv.Itoa(m.Expiry()),
					"amount: " + strconv.FormatFloat(m.Value, 'f', -1, 64),
				})
				sb.Reset()
			}
			r += m.Value
		}
	}

	// No need to output if def was not modified
	if e.Core.Flags.LogDebug && len(logDetails) > 1 {
		evt.Write("def_mods", logDetails)
	}

	return r
}
