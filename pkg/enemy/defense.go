package enemy

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type defenseMod struct {
	value float64
	modifier.Base
}

func (e *Enemy) AddDefMod(key string, dur int, val float64, hitlag bool) {
	mod := defenseMod{
		Base:  modifier.NewBase(key, e.Core.F+dur, hitlag),
		value: val,
	}
	var evt glog.Event
	overwrote, oldEvt := modifier.Add(&e.defenseMods, &mod, e.Core.F)
	if overwrote {
		e.Core.Log.NewEvent(
			"mod refreshed", glog.LogStatusEvent, -1,
			"overwrite", true,
			"key", mod.Key(),
			"expiry", mod.Expiry(),
		)
		evt = oldEvt
	} else {
		evt = e.Core.Log.NewEvent(
			"mod added", glog.LogStatusEvent, -1,
			"overwrite", false,
			"key", mod.Key(),
			"expiry", mod.Expiry(),
		)
	}
	evt.SetEnded(mod.Expiry())
	mod.SetEvent(evt)
}

func (e *Enemy) DeleteDefMod(key string) {
	m := modifier.Delete(&e.defenseMods, key)
	if m != nil {
		m.Event().SetEnded(e.Core.F)
		e.Core.Log.NewEvent("enemy mod deleted", glog.LogStatusEvent, -1, "key", key)
	}
}

func (c *Enemy) DefModIsActive(key string) bool {
	_, ok := modifier.FindCheckExpiry(&c.defenseMods, key, c.Core.F)
	return ok
}

func (t *Enemy) DefAdj(ai *combat.AttackInfo, evt glog.Event) float64 {
	var logDetails []interface{}
	var sb strings.Builder

	if t.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, 3*len(t.defenseMods))
	}

	var r float64
	for _, v := range t.defenseMods {
		if v.Expiry() > t.Core.F {
			if t.Core.Flags.LogDebug {
				sb.WriteString(v.Key())
				logDetails = append(logDetails, sb.String(), []string{
					"status: added",
					"expiry_frame: " + strconv.Itoa(v.Expiry()),
					"amount: " + strconv.FormatFloat(v.value, 'f', -1, 64),
				})
				sb.Reset()
			}
			r += v.value
		}
	}

	// No need to output if def was not modified
	if t.Core.Flags.LogDebug && len(logDetails) > 1 {
		evt.Write("def_mods", logDetails)
	}

	return r
}
