package enemy

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Add.
func (e *Enemy) AddStatus(key string, dur int, hitlag bool) {
	mod := combat.Status{
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
	overwrote, oldEvt := modifier.Add[modifier.Mod](&e.mods, &mod, e.Core.F)
	modifier.LogAdd("status", -1, &mod, e.Core.Log, overwrote, oldEvt)
}

func (e *Enemy) AddResistMod(mod combat.ResistMod) {
	mod.SetExpiry(e.Core.F)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&e.mods, &mod, e.Core.F)
	modifier.LogAdd("enemy", -1, &mod, e.Core.Log, overwrote, oldEvt)
}

func (e *Enemy) AddDefMod(mod combat.DefMod) {
	mod.SetExpiry(e.Core.F)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&e.mods, &mod, e.Core.F)
	modifier.LogAdd("enemy", -1, &mod, e.Core.Log, overwrote, oldEvt)
}

// Delete.

func (e *Enemy) deleteMod(key string) {
	m := modifier.Delete(&e.mods, key)
	if m != nil {
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
func (e *Enemy) resist(ai *combat.AttackInfo, evt glog.Event) float64 {
	var logDetails []interface{}
	var sb strings.Builder

	if e.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, 5*len(e.mods))
	}

	r := e.resists[ai.Element]
	for _, v := range e.mods {
		m, ok := v.(*combat.ResistMod)
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
	var logDetails []interface{}
	var sb strings.Builder

	if e.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, 3*len(e.mods))
	}

	var r float64
	for _, v := range e.mods {
		m, ok := v.(*combat.DefMod)
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
