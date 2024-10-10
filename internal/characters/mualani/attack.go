package mualani

import (
	"fmt"
	"slices"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const normalHitNum = 3

var (
	attackFrames   [][]int
	attackHitmarks = []int{11, 9, 30}
	attackHitboxes = []float64{0.7, 0.7, 0.7}

	sharkBiteFrames      [][]int
	sharkBiteHitmarks    = []int{7, 7, 7, 42}
	sharkBiteHitboxes    = 1.0
	sharkMissileHitboxes = 5.0
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 33)
	attackFrames[0][action.ActionAttack] = 23
	attackFrames[0][action.ActionCharge] = 22

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 32)
	attackFrames[1][action.ActionAttack] = 22
	attackFrames[1][action.ActionCharge] = 24

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 67)
	attackFrames[1][action.ActionAttack] = 54
	attackFrames[2][action.ActionCharge] = 52

	sharkBiteFrames = make([][]int, 4)

	sharkBiteFrames[0] = frames.InitAbilSlice(215) // swap
	sharkBiteFrames[0][action.ActionAttack] = 109
	sharkBiteFrames[0][action.ActionSkill] = 39
	sharkBiteFrames[0][action.ActionBurst] = 40
	sharkBiteFrames[0][action.ActionDash] = 40
	sharkBiteFrames[0][action.ActionJump] = 40
	sharkBiteFrames[0][action.ActionWalk] = 39

	sharkBiteFrames[1] = sharkBiteFrames[0]
	sharkBiteFrames[2] = sharkBiteFrames[0]

	sharkBiteFrames[3] = frames.InitAbilSlice(258)
	sharkBiteFrames[3][action.ActionAttack] = 144
	sharkBiteFrames[3][action.ActionSkill] = 79
	sharkBiteFrames[3][action.ActionBurst] = 82
	sharkBiteFrames[3][action.ActionDash] = 78
	sharkBiteFrames[3][action.ActionJump] = 81
	sharkBiteFrames[3][action.ActionWalk] = 81
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.sharkBite(p), nil
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(
			c.Core.Combat.PrimaryTarget(),
			nil,
			attackHitboxes[c.NormalCounter],
		),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackFrames[c.NormalCounter][action.ActionSwap],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) sharkBite(p map[string]int) action.Info {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	c.NormalCounter = 0
	c.momentumSrc = c.Core.F
	momentumStacks := c.momentumStacks

	nextMomentumFrame := max(c.StatusDuration(momentumIcd), sharkBiteFrames[momentumStacks][action.ActionWalk])
	c.QueueCharTask(c.momentumStackGain(c.momentumSrc), nextMomentumFrame)
	c.QueueCharTask(func() {
		c.momentumStacks = 0
		mult := bite[c.TalentLvlSkill()] + momentumBonus[c.TalentLvlSkill()]*float64(momentumStacks) + c.c1()
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Sharky's Bite (%v momentum)", momentumStacks),
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
		}

		if momentumStacks >= 3 {
			ai.Abil = "Sharky's Surging Bite"
			mult += surgingBite[c.TalentLvlSkill()]
		}

		ap := combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key())

		enemiesBite := c.Core.Combat.EnemiesWithinArea(
			ap,
			nil,
		)

		markOfPrey := false
		for _, e := range enemiesBite {
			if e.StatusIsActive(markedAsPreyKey) {
				markOfPrey = true
			}
		}
		totalEnemies := len(enemiesBite)

		var enemiesMissile []combat.Enemy
		if markOfPrey {
			ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, sharkMissileHitboxes)
			enemiesMissile = c.Core.Combat.EnemiesWithinArea(
				ap,
				func(e combat.Enemy) bool { return !slices.Contains(enemiesBite, e) },
			)
			totalEnemies += len(enemiesMissile)
		}

		mult *= max(1.00-0.14*float64(totalEnemies-1), 0.72)
		ai.FlatDmg = mult * c.MaxHP()
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.PrimaryTarget(),
				nil,
				sharkBiteHitboxes,
			),
			0,
			0,
			c.particleCB,
			c.a1cb(),
		)

		for _, e := range enemiesMissile {
			c.Core.QueueAttack(
				ai,
				combat.NewSingleTargetHit(e.Key()),
				0,
				travel,
			)
		}

		c.SetCDWithDelay(action.ActionAttack, 1.8*60, 0)
	}, sharkBiteHitmarks[momentumStacks])

	minAction := action.ActionWalk
	if momentumStacks >= 3 {
		minAction = action.ActionDash
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(sharkBiteFrames[momentumStacks]),
		AnimationLength: sharkBiteFrames[momentumStacks][action.WalkState], // shorter animation state so that a single bite doesn't make 3 yelan/xq waves. In game it only does 1.
		CanQueueAfter:   sharkBiteFrames[momentumStacks][minAction],
		State:           action.NormalAttackState,
	}
}
