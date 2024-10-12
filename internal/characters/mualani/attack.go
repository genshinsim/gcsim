package mualani

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const normalHitNum = 3

var (
	attackFrames   [][]int
	attackHitmarks = []int{11, 9, 31}

	sharkBiteFrames      [][]int
	sharkBiteHitmarks    = []int{7, 7, 7, 42}
	sharkMissileHitboxes = 5.0
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 33) // walk
	attackFrames[0][action.ActionAttack] = 23
	attackFrames[0][action.ActionCharge] = 22
	attackFrames[0][action.ActionSkill] = 13
	attackFrames[0][action.ActionBurst] = 11
	attackFrames[0][action.ActionDash] = 12
	attackFrames[0][action.ActionJump] = 13
	attackFrames[0][action.ActionSwap] = 11

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 32) // walk
	attackFrames[1][action.ActionAttack] = 22
	attackFrames[1][action.ActionCharge] = 24
	attackFrames[1][action.ActionSkill] = 11
	attackFrames[1][action.ActionBurst] = 11
	attackFrames[1][action.ActionDash] = 10
	attackFrames[1][action.ActionJump] = 10
	attackFrames[1][action.ActionSwap] = 8

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 67) // walk
	attackFrames[2][action.ActionAttack] = 54
	attackFrames[2][action.ActionCharge] = 52
	attackFrames[2][action.ActionSkill] = 32
	attackFrames[2][action.ActionBurst] = 33
	attackFrames[2][action.ActionDash] = 32
	attackFrames[2][action.ActionJump] = 34
	attackFrames[2][action.ActionSwap] = 32

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

	sharkBiteFrames[3] = frames.InitAbilSlice(258) // swap
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
			0.7,
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

	nextMomentumFrame := sharkBiteFrames[momentumStacks][action.ActionWalk]
	c.QueueCharTask(c.momentumStackGain(c.momentumSrc), nextMomentumFrame)
	c.QueueCharTask(func() {
		c.momentumStacks = 0
		mult := bite[c.TalentLvlSkill()] + momentumBonus[c.TalentLvlSkill()]*float64(momentumStacks) + c.c1()
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           fmt.Sprintf("Sharky's Bite (%v momentum)", momentumStacks),
			AttackTag:      attacks.AttackTagNormal,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        attributes.Hydro,
			Durability:     25,
		}

		if momentumStacks >= 3 {
			ai.Abil = "Sharky's Surging Bite"
			mult += surgingBite[c.TalentLvlSkill()]
		}

		primaryEnemy, ok := c.Core.Combat.PrimaryTarget().(combat.Enemy)
		if !ok {
			return
		}
		var enemiesMissile []combat.Enemy
		if primaryEnemy.StatusIsActive(markedAsPreyKey) {
			ap := combat.NewCircleHitOnTarget(primaryEnemy, nil, sharkMissileHitboxes)
			enemiesMissile = c.Core.Combat.EnemiesWithinArea(
				ap,
				func(e combat.Enemy) bool { return e.StatusIsActive(markedAsPreyKey) && e != primaryEnemy },
			)
			neighbours := len(enemiesMissile)
			mult *= max(1.00-0.14*float64(neighbours), 0.72)
		}

		ai.FlatDmg = mult * c.MaxHP()
		c.Core.QueueAttack(
			ai,
			combat.NewSingleTargetHit(primaryEnemy.Key()),
			0,
			0,
			c.particleCB,
			c.removeEnemyMarkCB,
			c.a1cb(),
		)

		ai.Abil = "Shark Missile"
		for _, e := range enemiesMissile {
			c.Core.QueueAttack(
				ai,
				combat.NewSingleTargetHit(e.Key()),
				0,
				travel,
				c.removeEnemyMarkCB,
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

func (c *char) removeEnemyMarkCB(a combat.AttackCB) {
	enemy, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	enemy.DeleteStatus(markedAsPreyKey)
}
