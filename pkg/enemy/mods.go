package enemy

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type resistMod struct {
	ele   attributes.Element
	value float64
	modifier.Base
}

type defenseMod struct {
	value float64
	modifier.Base
}

// Add helpers

func (e *Enemy) AddResistMod(key string, dur int, ele attributes.Element, val float64, hitlag bool) {
	mod := resistMod{
		Base:  modifier.NewBase(key, e.Core.F+dur, hitlag),
		ele:   ele,
		value: val,
	}
	overwrote, oldEvt := modifier.Add[modifier.Mod](&e.mods, &mod, e.Core.F)
	modifier.LogAdd("enemy", &mod, e.Core.Log, overwrote, oldEvt)
}

func (e *Enemy) AddDefMod(key string, dur int, val float64, hitlag bool) {
	mod := defenseMod{
		Base:  modifier.NewBase(key, e.Core.F+dur, hitlag),
		value: val,
	}
	overwrote, oldEvt := modifier.Add[modifier.Mod](&e.mods, &mod, e.Core.F)
	modifier.LogAdd("enemy", &mod, e.Core.Log, overwrote, oldEvt)
}

// Delete helpers

func (e *Enemy) deleteMod(key string) {
	m := modifier.Delete(&e.mods, key)
	if m != nil {
		modifier.LogDelete("enemy", m, e.Core.Log, e.Core.F)
	}
}

func (e *Enemy) DeleteResistMod(key string) { e.deleteMod(key) }
func (e *Enemy) DeleteDefMod(key string)    { e.deleteMod(key) }

// Active check
func (e *Enemy) modIsActive(key string) bool {
	_, ok := modifier.FindCheckExpiry(&e.mods, key, e.Core.F)
	return ok
}
func (e *Enemy) ResistModIsActive(key string) bool { return e.modIsActive(key) }
func (e *Enemy) DefModIsActive(key string) bool    { return e.modIsActive(key) }

// Values

func (e *Enemy) Resist(ai *combat.AttackInfo, evt glog.Event) float64 {
	var logDetails []interface{}
	var sb strings.Builder

	if e.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, 5*len(e.mods))
	}

	r := e.resist[ai.Element]
	for _, v := range e.mods {
		m, ok := v.(*resistMod)
		if !ok {
			continue
		}
		if m.Expiry() > e.Core.F && m.ele == ai.Element {
			if e.Core.Flags.LogDebug {
				sb.WriteString(m.Key())
				logDetails = append(logDetails, sb.String(), []string{
					"status: added",
					"expiry_frame: " + strconv.Itoa(m.Expiry()),
					"ele: " + m.ele.String(),
					"amount: " + strconv.FormatFloat(m.value, 'f', -1, 64),
				})
				sb.Reset()
			}
			r += m.value
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
		m, ok := v.(*defenseMod)
		if !ok {
			continue
		}
		if m.Expiry() > t.Core.F {
			if t.Core.Flags.LogDebug {
				sb.WriteString(m.Key())
				logDetails = append(logDetails, sb.String(), []string{
					"status: added",
					"expiry_frame: " + strconv.Itoa(m.Expiry()),
					"amount: " + strconv.FormatFloat(m.value, 'f', -1, 64),
				})
				sb.Reset()
			}
			r += m.value
		}
	}

	// No need to output if def was not modified
	if t.Core.Flags.LogDebug && len(logDetails) > 1 {
		evt.Write("def_mods", logDetails)
	}

	return r
}
