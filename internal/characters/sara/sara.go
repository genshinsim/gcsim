package sara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const normalHitNum = 5

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Sara, NewChar)
}

type char struct {
	*tmpl.Character
	a4LastProc int
	c1LastProc int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassBow
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	w.Character = &c

	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 19)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 25)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 38)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 41)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 58)

	// aimed -> x
	aimedFrames = frames.InitAbilSlice(78)

	// skill -> x
	skillFrames = frames.InitAbilSlice(50)

	// burst -> x
	burstFrames = frames.InitAbilSlice(60)
}
