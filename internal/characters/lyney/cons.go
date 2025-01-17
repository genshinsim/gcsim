package lyney

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1ICDKey = "lyney-c1-icd"
	c1ICD    = 15 * 60
)

// Lyney can have 2 Grin-Malkin Hats present at once.
// Additionally, Prop Arrows will summon 2 Grin-Malkin Hats and grant Lyney 1 extra stack of Prop Surplus.
// This effect can occur once every 15s.
func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	c.maxHatCount = 2
}

// C1 Prop stack is granted regardless of HP drain
func (c *char) addC1PropStack() func() {
	return func() {
		if c.Base.Cons < 1 || c.StatusIsActive(c1ICDKey) {
			return
		}
		c.increasePropSurplusStacks(1)
	}
}

func (c *char) c1HatIncrease() int {
	addCount := 0
	if c.Base.Cons >= 1 && !c.StatusIsActive(c1ICDKey) {
		addCount = 1
		c.AddStatus(c1ICDKey, c1ICD, true)
	}
	return addCount
}

// TODO: proper frames?
const c2Interval = 2 * 60

// When Lyney is on the field, he will gain a stack of Crisp Focus every 2s.
// This will increase his CRIT DMG by 20%. Max 3 stacks.
// This effect will be canceled when Lyney leaves the field.
func (c *char) c2Setup() {
	if c.Base.Cons < 2 {
		return
	}

	// listen for swap to clear/apply C2
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		// swapping off lyney means clearing C2
		prev := args[0].(int)
		if prev == c.Index {
			c.c2Src = -1
			c.c2Stacks = 0
			return false
		}
		// swapping to lyney means applying C2
		next := args[1].(int)
		if next == c.Index {
			c.c2Src = c.Core.F
			c.QueueCharTask(c.c2StackCheck(c.Core.F), c2Interval)
			c.Core.Log.NewEvent("Lyney C2 started", glog.LogCharacterEvent, c.Index).Write("c2_stacks", c.c2Stacks)
			return false
		}
		return false
	}, "lyney-c2-swap")

	// add buff
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("lyney-c2", -1),
		AffectedStat: attributes.CD,
		Amount: func() ([]float64, bool) {
			m[attributes.CD] = float64(c.c2Stacks) * 0.2
			return m, true
		},
	})
}

func (c *char) c2StackCheck(src int) func() {
	return func() {
		// don't add stack if src changed via swapping off
		if src != c.c2Src {
			return
		}
		// don't add stack if no longer on-field
		// sanity check, should be guaranteed via previous check + event subscription
		if c.Index != c.Core.Player.Active() {
			return
		}
		// don't add stack if already at max
		// no way of losing stacks except swap so no need to queue up stack check if already at max
		if c.c2Stacks == 3 {
			return
		}
		// add stack
		c.c2Stacks++
		c.Core.Log.NewEvent("Lyney C2 stack added", glog.LogCharacterEvent, c.Index).Write("c2_stacks", c.c2Stacks)
		// queue up stack check
		c.QueueCharTask(c.c2StackCheck(src), c2Interval)
	}
}

// After an opponent is hit by Lyney's Pyro Charged Attack, this opponent's Pyro RES will be decreased by 20% for 6s.
func (c *char) makeC4CB() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("lyney-c4", 6*60),
			Ele:   attributes.Pyro,
			Value: -0.20,
		})
	}
}

// When Lyney fires a Prop Arrow, he will fire a Pyrotechnic Strike: Reprised that will deal 80% of a Pyrotechnic Strike's DMG.
// This DMG is considered Charged Attack DMG.
func (c *char) c6(c6Travel int) {
	if c.Base.Cons < 6 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Pyrotechnic Strike: Reprised",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagLyneyEndBoom,
		ICDGroup:   attacks.ICDGroupLyneyExtra,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       propPyrotechnic[c.TalentLvlAttack()] * 0.8,
	}
	// TODO: snapshot
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			1,
		),
		0,
		c6Travel,
		c.makeC4CB(),
	)
}
