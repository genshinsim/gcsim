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
	modifier.Add(e.Core.F, e.Core.Log, &e.defenseMods, &mod)
}

func (e *Enemy) DeleteDefMod(key string) {
	modifier.Delete(e.Core.F, e.Core.Log, &e.defenseMods, key)
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
