package venti

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const normalHitNum = 6

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Venti, NewChar)
}

type char struct {
	*tmpl.Character
	qInfuse             attributes.Element
	infuseCheckLocation combat.AttackPattern
	aiAbsorb            combat.AttackInfo
	snapAbsorb          combat.Snapshot
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Anemo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassBow
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.infuseCheckLocation = combat.NewDefCircHit(0.1, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableObject)

	w.Character = &c

	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][1], 30)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 38)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 33)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 31)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 22)
	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5][0], 98)

	// high_plunge -> x
	// TODO: missing counts for plunge cancels?
	// using hitmark as placeholder for now
	highPlungeFrames = frames.InitAbilSlice(highPlungeHitmark)

	// aimed -> x
	// TODO: get separate counts for each cancel, currently using generic frames for all of them
	aimedFrames = frames.InitAbilSlice(94) // shouldn't this be 86?

	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(98)
	skillPressFrames[action.ActionAttack] = 22
	skillPressFrames[action.ActionAim] = 22   // assumed
	skillPressFrames[action.ActionSkill] = 22 // uses burst frames
	skillPressFrames[action.ActionBurst] = 22
	skillPressFrames[action.ActionDash] = 22
	skillPressFrames[action.ActionJump] = 22

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(289)
	skillHoldFrames[action.ActionHighPlunge] = 116

	// burst -> x
	burstFrames = frames.InitAbilSlice(96)
	burstFrames[action.ActionAttack] = 95
	burstFrames[action.ActionAim] = 95 // assumed
	burstFrames[action.ActionDash] = 95
	burstFrames[action.ActionJump] = 95
	burstFrames[action.ActionSwap] = 94
}

func (c *char) ReceiveParticle(p character.Particle, isActive bool, partyCount int) {
	c.Character.ReceiveParticle(p, isActive, partyCount)
	if c.Base.Cons >= 4 {
		//only pop this if active
		if !isActive {
			return
		}
		m := make([]float64, attributes.EndStatType)
		m[attributes.AnemoP] = 0.25
		c.AddStatMod("venti-c4", 600, attributes.AnemoP, func() ([]float64, bool) {
			return m, true
		})
	}
}
