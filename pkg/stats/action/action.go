package action

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(stats.Config{
		Name: "action",
		New:  NewStat,
	})
}

type buffer struct {
	energySpent []float64
	failures    [][]stats.ActionFailInterval
	// per char, the actions that failed and haven't been executed since, in the order they
	// first failed. A char waits on one action at a time, so this rarely holds more than one
	// and a scan is cheaper than a map lookup on every frame an action waits
	activeFailures [][]activeFailure
	actionEvents   [][]stats.ActionEvent
}

type activeFailure struct {
	action action.Action
	start  int
	reason action.Failure
}

// findFailure returns the index of the active failure of e, or -1
func findFailure(active []activeFailure, e action.Action) int {
	for i := range active {
		if active[i].action == e {
			return i
		}
	}
	return -1
}

func (b buffer) addFailure(core *core.Core, char int, active activeFailure) {
	interval := stats.ActionFailInterval{
		Start:  active.start,
		End:    core.F,
		Reason: active.reason.String(),
	}

	// TODO: limit intervals to be at least length x (5?)
	b.failures[char] = append(b.failures[char], interval)
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		energySpent:    make([]float64, len(core.Player.Chars())),
		failures:       make([][]stats.ActionFailInterval, len(core.Player.Chars())),
		activeFailures: make([][]activeFailure, len(core.Player.Chars())),
		actionEvents:   make([][]stats.ActionEvent, len(core.Player.Chars())),
	}

	core.Events.Subscribe(event.OnActionExec, func(args ...any) {
		char := args[0].(int)
		e := args[1].(action.Action)

		if e == action.ActionBurst {
			out.energySpent[char] += core.Player.Chars()[char].EnergyMax
		}

		// TODO: ActionId population
		event := stats.ActionEvent{
			Frame:  core.F,
			Action: e.String(),
		}
		out.actionEvents[char] = append(out.actionEvents[char], event)

		if i := findFailure(out.activeFailures[char], e); i != -1 {
			out.addFailure(core, char, out.activeFailures[char][i])
			out.activeFailures[char] = slices.Delete(out.activeFailures[char], i, i+1)
		}
	}, "stats-action-exec-log")

	core.Events.Subscribe(event.OnActionFailed, func(args ...any) {
		char := args[0].(int)
		e := args[1].(action.Action)
		reason := args[3].(action.Failure)

		// Assumes we will continue trying an action until it succeeds.
		// If we ever give up trying actions, this will no longer be accurate
		// TODO: track by action id to handle this edge case?
		if findFailure(out.activeFailures[char], e) == -1 {
			out.activeFailures[char] = append(out.activeFailures[char], activeFailure{
				action: e,
				start:  core.F,
				reason: reason,
			})
		}
	}, "stats-action-failed-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(core.Player.Chars()); c++ {
		for _, active := range b.activeFailures[c] {
			b.addFailure(core, c, active)
		}

		result.Characters[c].FailedActions = b.failures[c]
		result.Characters[c].EnergySpent = b.energySpent[c]
		result.Characters[c].ActionEvents = b.actionEvents[c]
	}
}
