package enemy

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (e *Enemy) ApplyHitlag(factor, dur float64) {
	//TODO: extend all hitlag affected buff expiry by dur * (1 - factor) i think
	ext := dur * (1 - factor)

	var logs []string
	var evt glog.Event
	if e.Core.Flags.LogDebug {
		logs = make([]string, 0, len(e.mods))
		evt = e.Core.Log.NewEvent("enemy hitlag - extending mods", glog.LogHitlagEvent, -1, "target", e.Index())
	}

	//check resist mods
	for i, v := range e.mods {
		if v.AffectedByHitlag() {
			e.mods[i].Extend(ext)
			if e.Core.Flags.LogDebug {
				logs = append(logs, fmt.Sprintf("%v: %v", v.Key(), v.Expiry()))
			}
		}
	}

	if e.Core.Flags.LogDebug {
		evt.Write("mods affected", logs)
	}
}
