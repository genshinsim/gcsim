package cryo

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	chargeFrames   [][]int
	chargeHitmarks = [][]int{{9, 20}, {14, 25}}
)

const (
	chargeTrueMoonICDKey = "travelercryo-special-ca-icd"
	trueMoonStackICDKey  = "travelercryo-truemoon-icd"
)

func init() {
	chargeFrames = make([][]int, 2)
	// Male
	chargeFrames[0] = frames.InitAbilSlice(55)                                       // CA -> N1
	chargeFrames[0][action.ActionSkill] = 37                                         // CA -> E
	chargeFrames[0][action.ActionBurst] = 36                                         // CA -> Q
	chargeFrames[0][action.ActionDash] = chargeHitmarks[0][len(chargeHitmarks[0])-1] // CA -> D
	chargeFrames[0][action.ActionJump] = chargeHitmarks[0][len(chargeHitmarks[0])-1] // CA -> J
	chargeFrames[0][action.ActionSwap] = 44                                          // CA -> Swap

	// Female
	chargeFrames[1] = frames.InitAbilSlice(58)                                       // CA -> N1
	chargeFrames[1][action.ActionSkill] = 34                                         // CA -> E
	chargeFrames[1][action.ActionBurst] = 35                                         // CA -> Q
	chargeFrames[1][action.ActionDash] = chargeHitmarks[1][len(chargeHitmarks[1])-1] // CA -> D
	chargeFrames[1][action.ActionJump] = chargeHitmarks[1][len(chargeHitmarks[1])-1] // CA -> J
	chargeFrames[1][action.ActionSwap] = chargeHitmarks[1][len(chargeHitmarks[1])-1] // CA -> Swap
}

func (c *Traveler) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}

	conversion := c.chargeAttackTruemoon()
	if conversion == nil {
		conversion = func(ai *info.AttackInfo, ap *info.AttackPattern) { c.a1Conversion(ai) }
	}

	for i, mult := range charge[c.gender] {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.2)
		conversion(&ai, &ap)
		c.Core.QueueAttack(
			ai,
			ap,
			chargeHitmarks[c.gender][i],
			chargeHitmarks[c.gender][i],
			c.naCaPlungeCB,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames[c.gender]),
		AnimationLength: chargeFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[c.gender][len(chargeHitmarks[c.gender])-1],
		State:           action.ChargeAttackState,
	}, nil
}

// Charged Attack: Freezing Ice
//   - Cryo DMG dealt is increased by 140% of the Traveler's ATK, and the Traveler also gains 2 stacks
//     of Frostglow.
//   - Charged Attack: Freezing Ice can be unleashed once every 15s. Additionally, the DMG Bonus from
//     the Ascension Talent "Ever-Keen Frost" does not apply to Charged Attack: Freezing Ice.
//
// Radiance: Stellar Glimmer: DMG from Charged Attack: Freezing Ice is instead considered Cryo DMG
// of the corresponding Stellar Glimmer reaction type.
func (c *Traveler) chargeAttackTruemoon() func(*info.AttackInfo, *info.AttackPattern) {
	if !c.trueMoonBuff {
		return nil
	}

	if c.trueMoonStacks < 3 {
		return nil
	}

	if !c.StatusIsActive(skillKey) {
		return nil
	}

	if c.StatusIsActive(chargeTrueMoonICDKey) {
		return nil
	}

	// the ICD seems to be added on hitmark (ICD lasts 914f for Lumine) but gets checked on CA start
	c.AddStatus(chargeTrueMoonICDKey, 15*60+chargeHitmarks[c.gender][0], true)
	c.trueMoonStacks = 0
	c.flostglowStacks = min(c.flostglowStacks+2, skillStacksMax)

	return func(ai *info.AttackInfo, ap *info.AttackPattern) {
		ai.Element = attributes.Cryo
		ai.IgnoreInfusion = true
		ai.ICDTag = attacks.ICDTagTravelerEnhancedCA
		ai.Mult += 1.4

		ai.HitlagHaltFrames = 0.05 * 60
		ai.HitlagFactor = 0.05
		ai.CanBeDefenseHalted = true

		*ap = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, 3.2)

		switch c.getRadiance() {
		case radianceStellarConduct:
			ai.Abil += stellarConductText
			ai.AttackTag = attacks.AttackTagDirectStellarConduct
			ai.ICDTag = attacks.ICDTagNone
			ai.Durability = 0
			ai.IgnoreDefPercent = 1
		case radianceStellarSwirl:
			ai.Abil += stellarSwirlText
			ai.AttackTag = attacks.AttackTagDirectStellarSwirl
			ai.ICDTag = attacks.ICDTagNone
			ai.Durability = 0
			ai.IgnoreDefPercent = 1
		}
	}
}

// Additionally, when a nearby party member deals Stellar-Conduct or Stellar Swirl DMG, the Traveler
// gains 1 stack of "Icepoint." 1 stack of Icepoint can be obtained in this way once every 2s, max 3
// stacks. When the Traveler unleashes a Charged Attack at 3 stacks of Icepoint, all stacks are
// consumed to change the attack into a special Charged Attack: Freezing Ice.
func (c *Traveler) trueMoonInit() {
	if !c.trueMoonBuff {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if !atk.Info.AttackTag.IsStellar() {
			return
		}

		if c.StatusIsActive(trueMoonStackICDKey) {
			return
		}

		c.AddStatus(trueMoonStackICDKey, 2*60, true)

		// TODO: the stacks only last 30s
		c.trueMoonStacks = min(c.trueMoonStacks+1, 3)
	}, "cryo-mc-truemoon")
}
