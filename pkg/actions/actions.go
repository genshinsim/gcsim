package actions

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core"
)

type Ctrl struct {
	core *core.Core
}

func NewActionCtrl(c *core.Core) *Ctrl {
	return &Ctrl{
		core: c,
	}
}

//return frames, if executed, any errors
func (a *Ctrl) Exec(n core.Command) (int, bool, error) {

	//check the type of commands first;

	switch n.Type() {
	case core.CommandTypeSimCmd:
	case core.CommandTypeUseAbility:
		return a.execAbility(n.(*core.Ability))
	}

	return 0, false, errors.New("invalid command")
}

func (a *Ctrl) execAbility(abil *core.Ability) (int, bool, error) {
	//all ability should have a target char
	char := a.core.Chars[a.core.CharPos[abil.Character]]

	//check that the action is ready; if not ready skip the frame
	if !char.ActionReady(abil.Typ, abil.Param) {
		a.core.Log.Warnw("frame", a.core.F, "event", core.LogSimEvent, "queued action is not ready, should not happen; skipping frame")
		return 0, false, nil
	}

	return 0, true, nil
}
