package yaoyao

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Yaoyao, NewChar)
}

type char struct {
	*tmpl.Character
	skillRadishAI    combat.AttackInfo
	burstRadishAI    combat.AttackInfo
	numYueguiJumping int
	yueguiJumping    []*yuegui
	a4Srcs           []int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillRadishAI = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Radish (Skill)",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		ICDGroup:           attacks.ICDGroupYaoyaoRadishSkill,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               skillRadishDMG[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.burstRadishAI = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Radish (Burst)",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagElementalBurst,
		ICDGroup:           attacks.ICDGroupYaoyaoRadishBurst,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               burstRadishDMG[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.onExitField()
	c.yueguiJumping = make([]*yuegui, 3)
	c.a4Srcs = make([]int, 4)
	return nil
}
