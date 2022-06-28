package enemy

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type ResistMod struct {
	Ele   attributes.Element
	Value float64
	modifier.Base
}

type DefMod struct {
	Value float64
	Dur   int
	modifier.Base
}

// Add.

func (e *Enemy) AddResistMod(mod ResistMod) {
	mod.SetExpiry(e.Core.F)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&e.mods, &mod, e.Core.F)
	modifier.LogAdd("enemy", -1, &mod, e.Core.Log, overwrote, oldEvt)
}

func (e *Enemy) AddDefMod(mod DefMod) {
	mod.SetExpiry(e.Core.F)
	overwrote, oldEvt := modifier.Add[modifier.Mod](&e.mods, &mod, e.Core.F)
	modifier.LogAdd("enemy", -1, &mod, e.Core.Log, overwrote, oldEvt)
}

// Delete.

func (e *Enemy) deleteMod(key string) {
	m := modifier.Delete(&e.mods, key)
	if m != nil {
		modifier.LogDelete("enemy", -1, m, e.Core.Log, e.Core.F)
	}
}

func (e *Enemy) DeleteResistMod(key string) { e.deleteMod(key) }
func (e *Enemy) DeleteDefMod(key string)    { e.deleteMod(key) }

// Active.
func (e *Enemy) modIsActive(key string) bool {
	_, ok := modifier.FindCheckExpiry(&e.mods, key, e.Core.F)
	return ok
}
func (e *Enemy) ResistModIsActive(key string) bool { return e.modIsActive(key) }
func (e *Enemy) DefModIsActive(key string) bool    { return e.modIsActive(key) }

// Amount.

func (e *Enemy) Resist(ai *combat.AttackInfo, evt glog.Event) float64 {
	var logDetails []interface{}
	var sb strings.Builder

	if e.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, 5*len(e.mods))
	}

	r := e.resist[ai.Element]
	for _, v := range e.mods {
		m, ok := v.(*ResistMod)
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

func (t *Enemy) DefAdj(ai *combat.AttackInfo, evt glog.Event) float64 {
	var logDetails []interface{}
	var sb strings.Builder

	if t.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, 3*len(t.mods))
	}

	var r float64
	for _, v := range t.mods {
		m, ok := v.(*DefMod)
		if !ok {
			continue
		}
		if m.Expiry() > t.Core.F {
			if t.Core.Flags.LogDebug {
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
	if t.Core.Flags.LogDebug && len(logDetails) > 1 {
		evt.Write("def_mods", logDetails)
	}

	return r
}
