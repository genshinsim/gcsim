package nahida

import (
	"fmt"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Nahida, NewChar)
}

type char struct {
	*tmpl.Character
	triKarmaInterval int
	burstSrc         int
	pyroCount        int
	electroCount     int
	hydroCount       int
	pyroBurstBuff    []float64
	a1Buff           []float64
	a4Buff           []float64
	c4Buff           []float64
	c6Count          int
	markCount        int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 50
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	// skill hooks
	c.Core.Events.Subscribe(event.OnEnemyDamage, c.triKarmaOnBloomDamage, "nahida-tri-karma")
	// considers shatter as an elemental reaction
	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Core.Events.Subscribe(i, c.triKarmaOnReaction, fmt.Sprintf("nahida-tri-karma-on-%v", i))
	}
	// skill cooldown
	c.updateTriKarmaInterval()

	// burst ele counts
	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Pyro:
			c.pyroCount++
		case attributes.Hydro:
			c.hydroCount++
		case attributes.Electro:
			c.electroCount++
		}
	}

	if c.Base.Cons >= 1 {
		c.c1()
	}

	// sanity check
	if c.pyroCount > 2 {
		c.pyroCount = 2
	}
	if c.hydroCount > 2 {
		c.hydroCount = 2
	}
	if c.electroCount > 2 {
		c.electroCount = 2
	}

	c.pyroBurstBuff = make([]float64, attributes.EndStatType)
	if c.pyroCount > 0 {
		c.pyroBurstBuff[attributes.DmgP] = burstTriKarmaDmgBonus[c.pyroCount-1][c.TalentLvlBurst()]
	}

	c.a1Buff = make([]float64, attributes.EndStatType)

	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4()
	c.a4Tick()

	if c.Base.Cons >= 4 {
		c.c4Buff = make([]float64, attributes.EndStatType)
		c.c4()
	}

	if c.Base.Cons >= 2 {
		c.c2()
	}

	return nil
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 9
	}
	return c.Character.AnimationStartDelay(k)
}
