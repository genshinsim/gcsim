package dehya

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

type char struct {
	*tmpl.Character
	// tracking skill information
	recastBefore    bool
	nextIsRecast    bool
	skillArea       combat.AttackPattern
	skillAttackInfo combat.AttackInfo
	skillSnapshot   combat.Snapshot
	sanctumICD      int
	burstCast       int
	burstCounter    int
	punchSrc        bool
	c1var           []float64
	c6count         int
}

func init() {
	core.RegisterCharFunc(keys.Dehya, NewChar)
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.EnergyMax = 70
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.skillHook()
	c.a4()
	c.burstCast = -241
	c.c1var = []float64{0.0, 0.0}
	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	return nil
}
func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.ActionFailure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.StatusIsActive(dehyaFieldKey) && !c.recastBefore {
		c.nextIsRecast = true
		return true, action.NoFailure
	}
	if a == action.ActionSkill && c.StatusIsActive(burstKey) {
		return true, action.NoFailure
	}
	if a == action.ActionAttack && c.StatusIsActive(burstKey) {
		return false, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.StatusIsActive(burstKey) {
			c.a1()
			c.DeleteStatus(burstKey)
			if remainingFieldDur > 0 { //place field
				c.QueueCharTask(func() {
					c.addField(remainingFieldDur)
				}, kickHitmark)
			}
		}
		return false
	}, "dehya-exit")
}

var burstIsJumpCancelled = false

func (c *char) Jump(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(burstKey) {
		burstIsJumpCancelled = true
		c.DeleteStatus(burstKey)
	}
	if remainingFieldDur > 0 { //place field
		c.QueueCharTask(func() {
			c.addField(remainingFieldDur)
		}, kickHitmark)
	}
	return c.Character.Jump(p)
}
