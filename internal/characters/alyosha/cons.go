package alyosha

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key    = "alyosha-c1"
	c1IcdKey = "alyosha-c1-icd"
	c6Key    = "alyosha-c6"
)

// Alyosha regenerates 15 Elemental Energy when nearby party members trigger an Electro-related
// reaction on an opponent. This effect can trigger once every 18s.
func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		if c.StatusIsActive(c1IcdKey) {
			return
		}
		c.AddStatus(c1IcdKey, 18*60, true)
		c.AddEnergy(c1Key, 15)
	}

	c.Core.Events.Subscribe(event.OnOverload, hook, c1Key)
	c.Core.Events.Subscribe(event.OnElectroCharged, hook, c1Key)
	c.Core.Events.Subscribe(event.OnLunarCharged, hook, c1Key)
	c.Core.Events.Subscribe(event.OnSuperconduct, hook, c1Key)
	c.Core.Events.Subscribe(event.OnStellarConduct, hook, c1Key)
	c.Core.Events.Subscribe(event.OnSwirlElectro, hook, c1Key)
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, hook, c1Key)
	c.Core.Events.Subscribe(event.OnHyperbloom, hook, c1Key)
	c.Core.Events.Subscribe(event.OnQuicken, hook, c1Key)
	c.Core.Events.Subscribe(event.OnAggravate, hook, c1Key)
}

// The duration of the Elemental Burst Hunter's Advance is extended by 6s.
// Additionally, each time Tugarin attacks the opponent, he will also apply the Hunter's Mark to a
// target. This effect does not activate existing Hunter's Marks.
func (c *char) c2BurstDur() int {
	if c.Base.Cons < 2 {
		return 0
	}

	return 6 * 60
}

func (c *char) c2MakeTugarinCB() info.AttackCBFunc {
	if c.Base.Cons < 2 {
		return c.triggerSkillMarkCB(false)
	}

	return c.triggerSkillMarkCB(true)
}

// When Tugarin attacks, he will also restore 60% of Alyosha's ATK as HP to the nearby party member
// with the lowest HP by percentage.
func (c *char) c4OnTugarin() {
	if c.Base.Cons < 4 {
		return
	}

	lowestHPInd := c.Index()
	lowestHP := c.CurrentHPRatio()
	for i, char := range c.Core.Player.Chars() {
		if char.CurrentHPRatio() < lowestHP {
			lowestHPInd = i
			lowestHP = char.CurrentHPRatio()
		}
	}

	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  lowestHPInd,
		Message: "Alyosha (C4)",
		Src:     0.6 * c.TotalAtk(),
		Bonus:   c.Stat(attributes.Heal),
	})
}

// The Hunter's Precision effect obtained from activating a Hunter's Mark can now stack, max 2
// stacks. When the stack count reaches 2, characters affected by Hunter's Precision will also have
// their Elemental Mastery increased by 100 if they are on the field.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 100
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c6Key+"-em", -1),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				if c.skillStacks < 2 {
					return nil
				}

				if c.Core.Player.Active() != char.Index() {
					return nil
				}

				if !c.StatusIsActive(skillBuffKey) {
					return nil
				}

				return m
			},
		})
	}
}

func (c *char) c6ExtraSkillStacks() int {
	if c.Base.Cons < 6 {
		return 0
	}
	return 1
}
