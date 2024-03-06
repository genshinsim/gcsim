package qiqi

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

const (
	talismanKey    = "qiqi-talisman"
	talismanICDKey = "qiqi-talisman-icd"
)

func init() {
	core.RegisterCharFunc(keys.Qiqi, NewChar)
}

type char struct {
	*tmpl.Character
	skillLastUsed     int
	skillHealSnapshot combat.Snapshot // Required as both on hit procs and continuous healing need to use this
}

// TODO: Not implemented - C6 (revival mechanic, not suitable for sim)
// C4 - Enemy Atk reduction, not useful in this sim version
func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.skillLastUsed = 0

	w.Character = &c

	return nil
}

// Ensures the set of targets are initialized properly
func (c *char) Init() error {
	c.a1()
	c.talismanHealHook()
	c.onNACAHitHook()
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}

// Helper function to calculate healing amount dynamically using current character stats, which has all mods applied
func (c *char) healDynamic(healScalePer, healScaleFlat []float64, talentLevel int) float64 {
	atk := c.Base.Atk + c.Weapon.BaseAtk*(1+c.Stat(attributes.ATKP)) + c.Stat(attributes.ATK)
	heal := healScaleFlat[talentLevel] + atk*healScalePer[talentLevel]
	return heal
}

// Helper function to calculate healing amount from a snapshot instance
func (c *char) healSnapshot(d *combat.Snapshot, healScalePer, healScaleFlat []float64, talentLevel int) float64 {
	atk := d.BaseAtk*(1+d.Stats[attributes.ATKP]) + d.Stats[attributes.ATK]
	heal := healScaleFlat[talentLevel] + atk*healScalePer[talentLevel]
	return heal
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 7
	}
	return c.Character.AnimationStartDelay(k)
}
