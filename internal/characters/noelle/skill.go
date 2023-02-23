package noelle

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var skillFrames []int

const skillHitmark = 14

func init() {
	skillFrames = frames.InitAbilSlice(78)
	skillFrames[action.ActionAttack] = 12
	skillFrames[action.ActionSkill] = 14 // uses burst frames
	skillFrames[action.ActionBurst] = 14
	skillFrames[action.ActionDash] = 11
	skillFrames[action.ActionJump] = 11
	skillFrames[action.ActionWalk] = 43
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Breastplate",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Geo,
		Durability:         50,
		Mult:               shieldDmg[c.TalentLvlSkill()],
		UseDef:             true,
		CanBeDefenseHalted: true,
	}
	snap := c.Snapshot(&ai)

	//add shield first
	defFactor := snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
	shieldhp := shieldFlat[c.TalentLvlSkill()] + shieldDef[c.TalentLvlSkill()]*defFactor
	c.Core.Player.Shields.Add(c.newShield(shieldhp, shield.ShieldNoelleSkill, 720))

	//activate shield timer, on expiry explode
	c.shieldTimer = c.Core.F + 720 //12 seconds

	c.a4Counter = 0

	// initial E hit can proc her heal
	cb := c.skillHealCB()

	// center on player
	// use char queue for this just to be safe in case of C4
	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2),
			0,
			0,
			cb,
		)
	}, skillHitmark)

	// handle C4
	if c.Base.Cons >= 4 {
		c.Core.Tasks.Add(func() {
			if c.shieldTimer == c.Core.F {
				//deal damage
				c.explodeShield()
			}
		}, 720)
	}

	c.SetCDWithDelay(action.ActionSkill, 24*60, 6)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHealCB() combat.AttackCBFunc {
	done := false
	return func(atk combat.AttackCB) {
		if atk.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		// check for healing
		if c.Core.Player.Shields.Get(shield.ShieldNoelleSkill) != nil {
			var prob float64
			if c.Base.Cons >= 1 && c.StatModIsActive(burstBuffKey) {
				prob = 1
			} else {
				prob = healChance[c.TalentLvlSkill()]
			}
			if c.Core.Rand.Float64() < prob {
				// heal target
				x := atk.AttackEvent.Snapshot.BaseDef*(1+atk.AttackEvent.Snapshot.Stats[attributes.DEFP]) + atk.AttackEvent.Snapshot.Stats[attributes.DEF]
				heal := shieldHeal[c.TalentLvlSkill()]*x + shieldHealFlat[c.TalentLvlSkill()]
				c.Core.Player.Heal(player.HealInfo{
					Caller:  c.Index,
					Target:  -1,
					Message: "Breastplate (Attack)",
					Src:     heal,
					Bonus:   atk.AttackEvent.Snapshot.Stats[attributes.Heal],
				})
				done = true
			}
		}
	}
}

// C4:
// When Breastplate's duration expires or it is destroyed by DMG, it will deal 400% ATK of Geo DMG to surrounding opponents.
func (c *char) explodeShield() {
	c.shieldTimer = 0
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Breastplate (C4)",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Geo,
		Durability:         50,
		Mult:               4,
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.15 * 60,
		CanBeDefenseHalted: true,
	}

	//center on player
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4), 0, 0)
}
