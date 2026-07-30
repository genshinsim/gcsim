// Package basicimport provides conservative, compile-safe actions for characters
// whose verified frame and runtime-state data is not yet available.
package basicimport

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type Scaling struct {
	Values []float64
	UseDef bool
	UseHP  bool
	UseEM  bool
}

type Profile struct {
	Name       string
	Element    attributes.Element
	Weapon     string
	Energy     float64
	SkillCD    int
	BurstCD    int
	Attack     []Scaling
	Charge     Scaling
	Collision  Scaling
	LowPlunge  Scaling
	HighPlunge Scaling
	Skill      Scaling
	Burst      Scaling
}

type Character struct {
	*tmpl.Character
	Data Profile
}

func New(s *core.Core, w *character.CharWrapper, data Profile) *Character {
	c := &Character{Character: tmpl.NewWithWrapper(s, w), Data: data}
	c.EnergyMax = data.Energy
	c.NormalHitNum = len(data.Attack)
	c.SkillCon = 3
	c.BurstCon = 5
	return c
}

func (c *Character) Init() error { return nil }

func (c *Character) attackInfo(abil string, tag attacks.AttackTag, element attributes.Element, scaling Scaling) info.AttackInfo {
	level := c.TalentLvlAttack()
	if tag == attacks.AttackTagElementalArt {
		level = c.TalentLvlSkill()
	}
	if tag == attacks.AttackTagElementalBurst {
		level = c.TalentLvlBurst()
	}
	mult := 0.0
	if level >= 0 && level < len(scaling.Values) {
		mult = scaling.Values[level]
	}
	return info.AttackInfo{
		ActorIndex: c.Index(), Abil: abil, AttackTag: tag,
		ICDTag: attacks.ICDTagNormalAttack, ICDGroup: attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault, Element: element, Durability: 25,
		Mult: mult, UseDef: scaling.UseDef, UseHP: scaling.UseHP, UseEM: scaling.UseEM,
	}
}

func (c *Character) Attack(map[string]int) (action.Info, error) {
	i := c.NormalCounter
	element := attributes.Physical
	if c.Data.Weapon == "Catalyst" {
		element = c.Data.Element
	}
	ai := c.attackInfo(fmt.Sprintf("Normal %d", i+1), attacks.AttackTagNormal, element, c.Data.Attack[i])
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1), 20, 20)
	c.AdvanceNormalIndex()
	f := frames.InitNormalCancelSlice(20, 40)
	return action.Info{Frames: frames.NewAbilFunc(f), AnimationLength: 40, CanQueueAfter: 20, State: action.NormalAttackState}, nil
}

func (c *Character) ChargeAttack(map[string]int) (action.Info, error) {
	element := attributes.Physical
	if c.Data.Weapon == "Catalyst" || c.Data.Weapon == "Bow" {
		element = c.Data.Element
	}
	ai := c.attackInfo("Charged Attack", attacks.AttackTagExtra, element, c.Data.Charge)
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 2), 30, 30)
	f := frames.InitAbilSlice(50)
	return action.Info{Frames: frames.NewAbilFunc(f), AnimationLength: 50, CanQueueAfter: 30, State: action.ChargeAttackState}, nil
}

func (c *Character) LowPlungeAttack(map[string]int) (action.Info, error) {
	return c.plunge(c.Data.LowPlunge, "Low Plunge")
}

func (c *Character) HighPlungeAttack(map[string]int) (action.Info, error) {
	return c.plunge(c.Data.HighPlunge, "High Plunge")
}

func (c *Character) plunge(s Scaling, name string) (action.Info, error) {
	c.Core.Player.SetAirborne(player.Grounded)
	ai := c.attackInfo(name, attacks.AttackTagPlunge, attributes.Physical, s)
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3), 40, 40)
	f := frames.InitAbilSlice(60)
	return action.Info{Frames: frames.NewAbilFunc(f), AnimationLength: 60, CanQueueAfter: 40, State: action.PlungeAttackState}, nil
}

func (c *Character) Skill(map[string]int) (action.Info, error) {
	ai := c.attackInfo("Elemental Skill", attacks.AttackTagElementalArt, c.Data.Element, c.Data.Skill)
	ai.ICDTag = attacks.ICDTagElementalArt
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3), 20, 20)
	c.SetCD(action.ActionSkill, c.Data.SkillCD)
	f := frames.InitAbilSlice(45)
	return action.Info{Frames: frames.NewAbilFunc(f), AnimationLength: 45, CanQueueAfter: 20, State: action.SkillState}, nil
}

func (c *Character) Burst(map[string]int) (action.Info, error) {
	ai := c.attackInfo("Elemental Burst", attacks.AttackTagElementalBurst, c.Data.Element, c.Data.Burst)
	ai.ICDTag = attacks.ICDTagElementalBurst
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5), 60, 60)
	c.SetCD(action.ActionBurst, c.Data.BurstCD)
	c.ConsumeEnergy(0)
	f := frames.InitAbilSlice(90)
	return action.Info{Frames: frames.NewAbilFunc(f), AnimationLength: 90, CanQueueAfter: 60, State: action.BurstState}, nil
}
