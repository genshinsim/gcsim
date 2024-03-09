package xianyun

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Xianyun, NewChar)
}

const noSrcVal = -1

type char struct {
	*tmpl.Character
	skillCounter        int
	skillSrc            int
	skillWasC6          bool
	skillEnemiesHit     []targets.TargetKey
	a1Buffer            []int
	a4Atk               float64
	a4src               int
	a4Max               float64
	a4Ratio             float64
	adeptalAssistStacks int

	// leapFrames  []int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3

	c.skillSrc = noSrcVal

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1Buffer = make([]int, len(c.Core.Player.Chars()))
	c.a1()
	c.a4()

	c.c1()
	c.c2()

	c.burstPlungeDoTTrigger()
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.StatusIsActive(skillStateKey) {
		return true, action.NoFailure
	}
	if (a == action.ActionAttack || a == action.ActionCharge) && c.StatusIsActive(skillStateKey) {
		return false, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "adeptal-assistance":
		return c.adeptalAssistStacks, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	return 10
}

func (c *char) getTotalAtk() float64 {
	stats, _ := c.Stats()
	return c.Base.Atk*(1+stats[attributes.ATKP]) + stats[attributes.ATK]
}
