package ganyu

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("ganyu", NewChar)
}

type char struct {
	*character.Tmpl
	a2expiry int
}

func NewChar(s core.Sim, log *zap.SugaredLogger, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 6
	c.BurstCon = 3
	c.SkillCon = 5

	//add a2
	val := make([]float64, core.EndStatType)
	val[core.CR] = 0.2
	c.AddMod(core.CharStatMod{
		Key: "ganyu-a2",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, c.a2expiry > c.Sim.Frame() && a == core.AttackTagExtra
		},
		Expiry: -1,
	})

	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.Tags["last"] = -1
	}

	return &c, nil
}
