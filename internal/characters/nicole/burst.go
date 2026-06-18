package nicole

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	burstHitmark      = 108
	projectionHitmark = 30
	burstKey          = "silent-contemplation"
	burstICDKey       = "nicole-burst-projection-icd"
)

type APFunc = func(*enemy.Enemy) info.AttackPattern

var (
	projectionHitbox map[keys.Char]func(*enemy.Enemy) info.AttackPattern
	defaultHitBox    func(*enemy.Enemy) info.AttackPattern
	burstFrames      []int
)

func init() {
	burstFrames = frames.InitAbilSlice(128)
	burstFrames[action.ActionAttack] = 117
	burstFrames[action.ActionCharge] = 120
	burstFrames[action.ActionSkill] = 116
	burstFrames[action.ActionDash] = 118
	burstFrames[action.ActionJump] = 117
	burstFrames[action.ActionSwap] = 113

	meleeOffset := info.Point{X: -1.5, Y: 2.5}
	rangeOffset := info.Point{X: -2.5, Y: 2.5}

	projectionHitbox = make(map[keys.Char]func(*enemy.Enemy) info.AttackPattern)

	norm := meleeOffset.Normalize()

	// is this correct? The additional 3.2 offset uses direction from the melee position to the target
	// varkaOffset := meleeOffset.Sub(norm.Scale(3.2))
	// projectionHitbox[keys.Varka] = func(t *enemy.Enemy) info.AttackPattern {
	// 	return combat.NewBoxHitOnTarget(t, varkaOffset, 4.5)
	// }

	// projectionHitbox[keys.Durin] = func(t *enemy.Enemy) info.AttackPattern {
	// 	return combat.NewCircleHitOnTarget(t, meleeOffset, 4.5)
	// }

	projectionHitbox[keys.Albedo] = func(t *enemy.Enemy) info.AttackPattern {
		return combat.NewBoxHitOnTarget(t, meleeOffset, 3, 8)
	}

	projectionHitbox[keys.Razor] = func(t *enemy.Enemy) info.AttackPattern {
		return combat.NewCircleHitOnTarget(t, meleeOffset, 4.5)
	}

	// is this correct? The additional X=-0.5, Y=-2 offset uses direction from the melee position to the target
	// lohenOffset := meleeOffset.Sub(norm.Scale(-2)).Add(norm.Perp().Scale(-0.5))
	// projectionHitbox[keys.Lohen] = func(t *enemy.Enemy) info.AttackPattern {
	// 	return combat.NewBoxHitOnTarget(t, lohenOffset, 4.5, 10)
	// }

	// projectionHitbox[keys.Prune] = func(t *enemy.Enemy) info.AttackPattern {
	// 	return combat.NewCircleHitOnTarget(t, meleeOffset, 4.5)
	// }

	projectionHitbox[keys.Sucrose] = func(t *enemy.Enemy) info.AttackPattern {
		return combat.NewCircleHitOnTarget(t, meleeOffset, 4.5)
	}

	// is this correct? The additional X=0.5, Y=1.5 offset uses direction from the melee position to the target
	kleeOffset := meleeOffset.Sub(norm.Scale(1.5)).Add(norm.Perp().Scale(0.5))
	fmt.Println("Klee offset", kleeOffset)
	projectionHitbox[keys.Klee] = func(t *enemy.Enemy) info.AttackPattern {
		return combat.NewCircleHitOnTarget(t, kleeOffset, 4.5)
	}

	// Venti, Fischl, Mona, use the same attack pattern as the default
	defaultHitBox = func(t *enemy.Enemy) info.AttackPattern {
		return combat.NewCircleHitOnTarget(t, rangeOffset, 4.5)
	}
}

func (c *char) Burst(_ map[string]int) (action.Info, error) {
	c.DeleteStatus(burstKey)

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Revelation: Ladder of Divine Ascent",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstInitial[c.TalentLvlBurst()],
	}

	c.ConsumeEnergy(12)
	c.SetCD(action.ActionBurst, 15*60)

	// initial hit
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.5), burstHitmark, burstHitmark)

	c.QueueCharTask(func() {
		c.projections = 0
		c.AddStatus(burstKey, 20*60, true)
	}, burstHitmark+1)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}

func (c *char) burstInit() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		ae := args[1].(*info.AttackEvent)

		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		if c.projections >= 4 {
			return
		}

		if !c.StatusIsActive(burstKey) {
			return
		}

		if c.StatusIsActive(burstICDKey) {
			return
		}

		c.AddStatus(burstICDKey, 3*60, true)
		c.projections += 1

		char := c.Core.Player.Chars()[ae.Info.ActorIndex]

		c.Core.Tasks.Add(func() {
			ai := info.AttackInfo{
				ActorIndex: char.Index(),
				Abil:       "Arcane Projection",
				AttackTag:  attacks.AttackTagNone,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    char.Base.Element,
				Mult:       burstProjection[c.TalentLvlBurst()],
				FlatDmg:    c.hexereiOnProjection(char),
			}

			apFn, ok := projectionHitbox[char.Base.Key]
			if !ok {
				apFn = defaultHitBox
			}
			ap := apFn(t)

			c.Core.QueueAttack(ai, ap, 0, 0)
		}, projectionHitmark)
	}, "nicole-burst-hook")
}
