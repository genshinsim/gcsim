package xinyan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const normalHitNum = 4

type char struct {
	*tmpl.Character
	shieldLevel int
	c1buff      []float64
	c2buff      []float64
}

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Xinyan, NewChar)
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Pyro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneLiyue

	w.Character = &c

	c.shieldLevel = 1

	return nil
}

func (c *char) Init() error {
	c.a4()

	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}

	return nil
}

// need to update frames
func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 25)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 44)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 70)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 65)

	skillFrames = frames.InitAbilSlice(65)
	burstFrames = frames.InitAbilSlice(98)
}
