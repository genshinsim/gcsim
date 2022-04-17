package character

import "github.com/genshinsim/gcsim/pkg/core"

func (t *Tmpl) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	t.Core.Log.NewEvent("ActionFrames not implemented", core.LogActionEvent, t.Index)
	return 0, 0
}

func (t *Tmpl) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {

	// check if next is skill hold
	if next == core.ActionSkill && p["hold"] != 0 {
		next = core.ActionSkillHoldFramesOnly
	}

	switch prev := t.Core.LastAction.Typ; prev {
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
	case core.ActionSkill:
		// check if prev was hold
		if t.Core.LastAction.Param["hold"] != 0 {
			prev = core.ActionSkillHoldFramesOnly
		}
		n, ok := t.cancelFrames[prev]
		if !ok {
			return 0
		}
		return n[next]
	default:
		n, ok := t.cancelFrames[prev]
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
