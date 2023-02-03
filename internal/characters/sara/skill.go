package sara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

// c2 hitmark
const c2Hitmark = 103

const coverKey = "sara-e-cover"

func init() {
	skillFrames = frames.InitAbilSlice(52) // E -> D
	skillFrames[action.ActionAttack] = 29  // E -> N1
	skillFrames[action.ActionAim] = 30     // E -> CA
	skillFrames[action.ActionBurst] = 32   // E -> Q
	skillFrames[action.ActionJump] = 51    // E -> J
	skillFrames[action.ActionSwap] = 50    // E -> Swap
}

// Implements skill handling. Fairly barebones since most of the actual stuff happens elsewhere
// Gains Crowfeather Cover for 18s, and when Kujou Sara fires a fully-charged Aimed Shot, Crowfeather Cover will be consumed, and will leave a Crowfeather at the target location.
// Crowfeathers will trigger Tengu Juurai: Ambush after a short time, dealing Electro DMG and granting the active character within its AoE an ATK Bonus based on Kujou Sara's Base ATK.
// The ATK Bonuses from different Tengu Juurai will not stack, and their effects and duration will be determined by the last Tengu Juurai to take effect.
// Also implements C2: Unleashing Tengu Stormcall will leave a Weaker Crowfeather at Kujou Sara's original position that will deal 30% of its original DMG.
func (c *char) Skill(p map[string]int) action.ActionInfo {

	// Snapshot for all of the crowfeathers are taken upon cast
	c.Core.Status.Add(coverKey, 18*60)

	// C2 handling
	if c.Base.Cons >= 2 {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Tengu Juurai: Ambush C2",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.3 * skill[c.TalentLvlSkill()],
		}
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)

		c.Core.QueueAttack(ai, ap, 50, c2Hitmark, c.a4)
		c.attackBuff(ap, c2Hitmark)
	}

	c.SetCDWithDelay(action.ActionSkill, 600, 7)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack], // earliest cancel
		State:           action.SkillState,
	}
}

// Handles attack boost from Sara's skills
// Checks for the onfield character at the delay frame, then applies buff to that character
func (c *char) attackBuff(a combat.AttackPattern, delay int) {
	c.Core.Tasks.Add(func() {
		// TODO: this should be a 0 dmg attack
		if collision, _ := c.Core.Combat.Player().AttackWillLand(a); !collision {
			return
		}

		active := c.Core.Player.ActiveChar()
		buff := atkBuff[c.TalentLvlSkill()] * float64(c.Base.Atk+c.Weapon.Atk)

		c.Core.Log.NewEvent("sara attack buff applied", glog.LogCharacterEvent, c.Index).
			Write("char", active.Index).
			Write("buff", buff).
			Write("expiry", c.Core.F+360)

		m := make([]float64, attributes.EndStatType)
		m[attributes.ATK] = buff
		active.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("sara-attack-buff", 360),
			AffectedStat: attributes.ATK,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		if c.Base.Cons >= 1 {
			c.c1()
		}
		if c.Base.Cons >= 6 {
			c.c6(active)
		}
	}, delay)
}
