package chiori

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Chiori, NewChar)
}

type char struct {
	*tmpl.Character

	// dolls
	skillSearchAoE   float64
	skillDoll        *ticker // 1st doll
	rockDoll         *ticker // 2nd doll from c1 / construct
	constructChecker *ticker

	// a1 tracking
	a1Triggered   bool
	a1AttackCount int

	a4Buff []float64

	// cons
	c1Active bool
	kinus    []*ticker
	c2Ticker *ticker

	c4AttackCount int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = base.SkillDetails.BurstEnergyCost
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.a1TapestrySetup()
	c.a4()

	c.skillSearchAoE = 12
	c.c1()
	c.c4()

	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if a1 window is active is on-field
	if a == action.ActionSkill && c.StatusIsActive(a1WindowKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 11
	case model.AnimationYelanN0StartDelay:
		return 3
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
