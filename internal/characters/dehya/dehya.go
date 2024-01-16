package dehya

import (
	"errors"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	// tracking skill information
	hasSkillRecast  bool
	hasC2DamageBuff bool
	skillArea       combat.AttackPattern
	skillAttackInfo combat.AttackInfo
	skillSnapshot   combat.Snapshot
	sanctumICD      int
	burstCast       int
	burstCounter    int
	burstHitSrc     int // I am using this value as a counter because if I use frame I can get duplicates
	c1var           []float64
	c6count         int
	sanctumSavedDur int
	burstJumpCancel bool
}

func init() {
	core.RegisterCharFunc(keys.Dehya, NewChar)
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
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
	if c.Base.Cons >= 2 {
		// TODO: figure out damage mitigation and move out of c2
		c.skillHurtHook()
		c.c2()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.StatusIsActive(dehyaFieldKey) && !c.hasSkillRecast {
		return true, action.NoFailure
	}
	if a == action.ActionSkill && (c.StatusIsActive(burstKey) || c.StatusIsActive(kickKey)) {
		return true, action.NoFailure
	}
	if a == action.ActionAttack && (c.StatusIsActive(burstKey) || c.StatusIsActive(kickKey)) {
		return false, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if !c.StatusIsActive(burstKey) && !c.StatusIsActive(kickKey) {
			return false
		}
		c.a1()
		c.DeleteStatus(burstKey)
		if dur := c.sanctumSavedDur; dur > 0 { // place field
			c.sanctumSavedDur = 0
			c.QueueCharTask(func() { c.addField(dur) }, kickHitmark)
		}

		return false
	}, "dehya-exit")
}

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if !c.StatusIsActive(burstKey) {
		if c.StatusIsActive(kickKey) {
			return action.Info{}, errors.New("can't jump cancel burst kick")
		}
		return c.Character.Jump(p)
	}

	c.burstJumpCancel = true
	c.DeleteStatus(burstKey)

	if dur := c.sanctumSavedDur; dur > 0 { // place field
		c.sanctumSavedDur = 0
		c.QueueCharTask(func() { c.addField(dur) }, kickHitmark)
	}
	return c.Character.Jump(p)
}
