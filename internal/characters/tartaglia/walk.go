package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var walkFrames []int

func init() {
	walkFrames = frames.InitAbilSlice(1)
	walkFrames[action.ActionSkill] = 4
}

func (c *char) Walk(p map[string]int) action.Info {
	f, ok := p["f"]
	if !ok {
		f = 1
	}
	animLength := walkFrames[action.ActionSkill]
	if animLength < f {
		animLength = f
	}
	return action.Info{
		Frames: func(next action.Action) int {
			if f < walkFrames[next] {
				return walkFrames[next]
			}
			return f
		},
		AnimationLength: animLength,
		CanQueueAfter:   f,
		State:           action.WalkState,
	}
}
