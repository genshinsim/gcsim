package character

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (t *Tmpl) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	t.coretype.Log.NewEvent("ActionFrames not implemented", coretype.LogActionEvent, t.Index)
	return 0, 0
}

func (t *Tmpl) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	switch t.Core.LastAction.Typ {
	case core.ActionSwap:
		//if not same character then there should be no delay
		return 0
	case core.ActionAttack:
		//check our hit counter; should be hit counter - 1
		lastHit := t.NormalCounter - 1
		if lastHit < 0 {
			lastHit = t.NormalHitNum - 1
		}
		n, ok := t.normalCancelFrames[lastHit]
		if !ok {
			//sanity check just in case no frames are set
			return 0
		}
		return n[next]
	default:
		n, ok := t.cancelFrames[t.Core.LastAction.Typ]
		if !ok {
			return 0
		}
		return n[next]
	}
}

func (t *Tmpl) SetNormalCancelFrames(normalHitNum int, nextAbil core.ActionType, frames int) {

	if _, ok := t.normalCancelFrames[normalHitNum]; !ok {
		t.normalCancelFrames[normalHitNum] = make(map[core.ActionType]int)
	}

	t.normalCancelFrames[normalHitNum][nextAbil] = frames
}

func (t *Tmpl) SetAbilCancelFrames(prevAbil core.ActionType, nextAbil core.ActionType, frames int) {

	if _, ok := t.cancelFrames[prevAbil]; !ok {
		t.cancelFrames[prevAbil] = make(map[core.ActionType]int)
	}

	t.cancelFrames[prevAbil][nextAbil] = frames
}
