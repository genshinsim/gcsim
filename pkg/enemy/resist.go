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

func (e *Enemy) AddResistMod(key string, dur int, ele attributes.Element, val float64, hitlag bool) {
	mod := resistMod{
		Base:  modifier.NewBase(key, e.Core.F+dur, hitlag),
		ele:   ele,
		value: val,
	}
	var evt glog.Event
	overwrote, oldEvt := modifier.Add(&e.resistMods, &mod, e.Core.F)
	if overwrote {
		e.Core.Log.NewEvent(
			"enemy mod refreshed", glog.LogStatusEvent, -1,
			"overwrite", true,
			"key", mod.Key(),
			"expiry", mod.Expiry(),
		)
		evt = oldEvt
	} else {
		evt = e.Core.Log.NewEvent(
			"enemy mod added", glog.LogStatusEvent, -1,
			"overwrite", false,
			"key", mod.Key(),
			"expiry", mod.Expiry(),
		)
	}
	evt.SetEnded(mod.Expiry())
	mod.SetEvent(evt)
}

func (e *Enemy) DeleteResistMod(key string) {
	m := modifier.Delete(&e.resistMods, key)
	if m != nil {
		m.Event().SetEnded(e.Core.F)
		e.Core.Log.NewEvent("enemy mod deleted", glog.LogStatusEvent, -1, "key", key)
	}
}

func (c *Enemy) ResistModIsActive(key string) bool {
	_, ok := modifier.FindCheckExpiry(&c.resistMods, key, c.Core.F)
	return ok
}

func (e *Enemy) Resist(ai *combat.AttackInfo, evt glog.Event) float64 {
	var logDetails []interface{}
	var sb strings.Builder

	if e.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, 5*len(e.resistMods))
	}

	r := e.resist[ai.Element]
	for _, v := range e.resistMods {
		if v.Expiry() > e.Core.F && v.ele == ai.Element {
			if e.Core.Flags.LogDebug {
				sb.WriteString(v.Key())
				logDetails = append(logDetails, sb.String(), []string{
					"status: added",
					"expiry_frame: " + strconv.Itoa(v.Expiry()),
					"ele: " + v.ele.String(),
					"amount: " + strconv.FormatFloat(v.value, 'f', -1, 64),
				})
				sb.Reset()
			}
			r += v.value
		}
	}

	// No need to output if resist was not modified
	if e.Core.Flags.LogDebug && len(logDetails) > 1 {
		evt.Write("resist_mods", logDetails)
	}

	return r
}
