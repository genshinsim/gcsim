package dehya

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

type char struct {
	*tmpl.Character
	// tracking skill information
	sanctumActive          bool
	recastBefore           bool
	nextIsRecast           bool
	sanctumRetrieved       bool
	skillArea              combat.AttackPattern
	skillAttackInfo        combat.AttackInfo
	skillSnapshot          combat.Snapshot
	sanctumSource          int
	sanctumExpiry          int
	sanctumICD             int
	sanctumPickupExtension int
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
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum
	c.sanctumPickupExtension = 24 // On recast from Burst/Skill-2 the field duration is extended by 0.4s

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillHook()
	c.a4()

	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}

	return nil
}
func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.ActionFailure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.sanctumActive && !c.recastBefore {
		c.nextIsRecast = true
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}
