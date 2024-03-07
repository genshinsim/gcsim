package thoma

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/minazuki"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Thoma, NewChar)
}

type char struct {
	*tmpl.Character
	a1Stack      int
	c6buff       []float64
	burstWatcher *minazuki.Watcher
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	c.a1Stack = 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	if c.Base.Cons >= 6 {
		c.c6buff = make([]float64, attributes.EndStatType)
		c.c6buff[attributes.DmgP] = .15
	}

	w, err := minazuki.New(
		minazuki.WithMandatory(keys.Thoma, "thoma burst", burstKey, burstICDKey, 60, c.summonFieryCollapse, c.Core),
		minazuki.WithAnimationDelayCheck(model.AnimationXingqiuN0StartDelay, c.shouldDelay),
	)
	if err != nil {
		return err
	}
	c.burstWatcher = w

	return nil
}

func (c *char) maxShieldHP() float64 {
	return shieldppmax[c.TalentLvlSkill()]*c.MaxHP() + shieldflatmax[c.TalentLvlSkill()]
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 11
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) shouldDelay() bool {
	return c.Core.Player.ActiveChar().NormalCounter == 1
}
