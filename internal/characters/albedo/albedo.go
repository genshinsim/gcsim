package albedo

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Albedo, NewChar)
}

type char struct {
	*tmpl.Character
	lastConstruct int
	bloomSnapshot combat.Snapshot
	// tracking skill information
	skillActive     bool
	skillArea       combat.AttackPattern
	skillAttackInfo combat.AttackInfo
	skillSnapshot   combat.Snapshot
	// c2 tracking
	c2stacks int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillHook()
	c.a1()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "elevator":
		return c.skillActive, nil
	case "c2stacks":
		return c.c2stacks, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 9
	}
	return c.Character.AnimationStartDelay(k)
}
