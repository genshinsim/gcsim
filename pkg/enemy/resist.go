package enemy

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type resistMod struct {
	ele   attributes.Element
	value float64
	tmpl
}

func (e *Enemy) AddResistMod(key string, dur int, ele attributes.Element, val float64) {
	mod := resistMod{
		tmpl: tmpl{
			key:    key,
			expiry: e.Core.F + dur,
		},
		ele:   ele,
		value: val,
	}
	addMod(e, &e.resistMods, &mod)
}

func (e *Enemy) DeleteResistMod(key string) {
	deleteMod(e, &e.resistMods, key)
}

func (c *Enemy) ResistModIsActive(key string) bool {
	_, ok := findModCheckExpiry(&c.resistMods, key, c.Core.F)
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
		if v.expiry > e.Core.F && v.ele == ai.Element {
			if e.Core.Flags.LogDebug {
				sb.WriteString(v.key)
				logDetails = append(logDetails, sb.String(), []string{
					"status: added",
					"expiry_frame: " + strconv.Itoa(v.expiry),
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
