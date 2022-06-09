package yaemiko

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const (
	normalHitNum   = 3
	yaeTotemCount  = "totems"
	yaeTotemStatus = "yae_oldest_totem_expiry"
)

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.YaeMiko, NewChar)
}

type char struct {
	*tmpl.Character
	kitsunes         []*kitsune
	totemParticleICD int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 90
	}
	c.Energy = float64(e)
	c.EnergyMax = 90
	c.Weapon.Class = weapon.WeaponClassCatalyst
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	c.SetNumCharges(action.ActionSkill, 3)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 21)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 23)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 21)

	// charge -> x
	chargeFrames = frames.InitAbilSlice(114)
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark

	// skill -> x
	skillFrames = frames.InitAbilSlice(34)

	// burst -> x
	burstFrames = frames.InitAbilSlice(111)
}
