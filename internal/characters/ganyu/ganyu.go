package ganyu

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("ganyu", NewChar)
}

type char struct {
	*character.Tmpl
	a2expiry int
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = def.WeaponClassBow
	c.NormalHitNum = 6
	c.BurstCon = 3
	c.SkillCon = 5

	//add a2
	val := make([]float64, def.EndStatType)
	val[def.CR] = 0.2
	c.AddMod(def.CharStatMod{
		Key: "ganyu-a2",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return val, c.a2expiry > c.Sim.Frame() && a == def.AttackTagExtra
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
