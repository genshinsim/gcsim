package furina

import (
	"math"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const ()

func init() {
	core.RegisterCharFunc(keys.Furina, NewChar)
}

type char struct {
	*tmpl.Character
	burstStacks         float64
	maxBurstStacks      float64
	burstBuff           []float64
	a4Buff              []float64
	cancelPreviousSkill bool
	a1HealsStopFrameMap []int
	a1HealsFlagMap      []bool
	lastSkillUseFrame   int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.maxBurstStacks = 300

	if c.Base.Cons >= 1 {
		c.maxBurstStacks = 450
	}

	c.burstStacks = 0
	c.burstBuff = make([]float64, attributes.EndStatType)

	c.a1HealsStopFrameMap = make([]int, len(c.Core.Player.Chars()))
	c.a1HealsFlagMap = make([]bool, len(c.Core.Player.Chars()))
	c.a1()

	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4()
	c.a4Tick()

	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("furina-burst-damage-buff", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				c.burstBuff[attributes.DmgP] = float64(c.burstStacks) * 0.0023

				return c.burstBuff, c.StatusIsActive(burstKey)
			},
		})

		char.AddHealBonusMod(character.HealBonusMod{
			Base: modifier.NewBase("furina-burst-heal-buff", -1),
			Amount: func() (float64, bool) {
				return float64(c.burstStacks) * 0.0009, c.StatusIsActive(burstKey)
			},
		})
	}

	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}

		di := args[0].(player.DrainInfo)

		if di.Amount <= 0 {
			return false
		}

		char := c.Core.Player.ByIndex(di.ActorIndex)
		stacksAmount := di.Amount / char.MaxHP() * 100

		if c.Base.Cons >= 2 {
			stacksAmount *= 3
		}

		c.burstStacks = min(c.maxBurstStacks, c.burstStacks+stacksAmount)

		return false
	}, "furina-burst-stack-on-hp-drain")

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}

		target := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)

		if amount <= 0 {
			return false
		}

		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}

		char := c.Core.Player.ByIndex(target)
		stacksAmount := (amount - overheal) / char.MaxHP() * 100

		if c.Base.Cons >= 2 {
			stacksAmount *= 3
		}

		c.burstStacks = min(c.maxBurstStacks, c.burstStacks+stacksAmount)

		return false
	}, "furina-burst-stack-on-heal")

	return nil
}
