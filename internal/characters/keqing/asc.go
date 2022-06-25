package keqing

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) a1() {
	//account for it starting somewhere around hitmark
	dur := 300 + skillRecastHitmark
	c.Core.Status.Add("keqinginfuse", dur)
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		"keqing-a1",
		attributes.Electro,
		dur,
		true,
		combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
	)
}

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.15
	m[attributes.ER] = 0.15

	c.AddStatMod("keqing-a4", 480, attributes.NoStat, func() ([]float64, bool) {
		return m, true
	})
}
