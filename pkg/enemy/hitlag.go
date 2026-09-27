package enemy

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (e *Enemy) ApplyHitlag(factor, dur float64) {
	// TODO: extend all hitlag affected buff expiry by dur * (1 - factor) i think
	ext := int(math.Ceil(dur * (1 - factor)))
	e.frozenFrames += ext

	var logs []string
	var evt glog.Event
	if e.Core.Flags.LogDebug {
		logs = make([]string, 0, len(e.mods))
		evt = e.Core.Log.NewEvent("enemy hitlag - extending mods", glog.LogHitlagEvent, -1).
			Write("target", e.Key()).
			Write("duration", dur).
			Write("factor", factor).
			Write("frozen_frames", e.frozenFrames).
			SetEnded(e.Core.F + int(math.Ceil(dur)))
	}

	// check resist mods
	for i, v := range e.mods {
		if v.AffectedByHitlag() && v.Expiry() != -1 && v.Expiry() > e.Core.F {
			mod := e.mods[i]
			mod.Extend(mod.Key(), e.Core.Log, -1, ext)
			if e.Core.Flags.LogDebug {
				logs = append(logs, fmt.Sprintf("%v: %v", v.Key(), v.Expiry()))
			}
		}
	}

	if e.Core.Flags.LogDebug {
		evt.Write("mods affected", logs)
	}
}

func (e *Enemy) QueueEnemyTask(f func(), delay int) {
	if delay == 0 {
		f()
		return
	}
	e.queue.Add(f, delay)
}

func (e *Enemy) Tick() {
	if e.frozenFrames > 0 {
		// frozen for this frame, do nothing
		e.frozenFrames--
		e.Core.Log.NewEvent("enemy skipping tick", glog.LogHitlagEvent, -1).
			Write("target", e.Key()).
			Write("frozen_for", e.frozenFrames)
		return
	}
	e.timePassed++

	e.queue.Run()
	e.Reactable.Tick()
}
