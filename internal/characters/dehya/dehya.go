package dehya

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

type char struct {
	*tmpl.Character
	// tracking skill information
	hasRecastSkill     bool
	hasC2DamageBuff    bool
	skillArea          combat.AttackPattern
	skillAttackInfo    combat.AttackInfo
	skillSnapshot      combat.Snapshot
	skillRedmanesBlood float64
	skillSelfDoTQueued bool
	sanctumSavedDur    int
	sanctumICD         int
	burstCounter       int
	burstHitSrc        int // I am using this value as a counter because if I use frame I can get duplicates
	c1FlatDmgRatioE    float64
	c1FlatDmgRatioQ    float64
	c6Count            int
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

	c.skillHurtHook()
	c.skillDmgHook()

	c.a4()

	c.c1()
	c.c2()
	c.c6()

	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.StatusIsActive(dehyaFieldKey) && !c.hasRecastSkill {
		return true, action.NoFailure
	}
	if a == action.ActionSkill && (c.StatusIsActive(burstKey) || c.StatusIsActive(kickKey)) {
		return true, action.NoFailure
	}
	if a == action.ActionAttack && (c.StatusIsActive(burstKey) || c.StatusIsActive(kickKey)) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if !c.StatusIsActive(burstKey) && !c.StatusIsActive(kickKey) {
			return false
		}
		c.DeleteStatus(burstKey)
		if dur := c.sanctumSavedDur; dur > 0 { // place field
			c.sanctumSavedDur = 0
			c.QueueCharTask(func() { c.addField(dur) }, burstKickHitmark)
		}

		return false
	}, "dehya-exit")
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 22
	case model.AnimationYelanN0StartDelay:
		return 22
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
