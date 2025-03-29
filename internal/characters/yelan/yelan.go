package yelan

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/minazuki"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

const (
	breakthroughStatus = "yelan_breakthrough"
	c6Status           = "yelan_c6"
	burstKey           = "yelanburst"
	burstICDKey        = "yelanburstICD"
)

func init() {
	core.RegisterCharFunc(keys.Yelan, NewChar)
}

type char struct {
	*tmpl.Character
	a4buff       []float64
	breakthrough bool // tracks breakthrough state
	c2icd        int
	c6count      int
	c4count      int // keep track of number of enemies tagged
	burstWatcher *minazuki.Watcher
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.c2icd = 0
	c.c6count = 0

	breakthrough, ok := p.Params["breakthrough"]
	if !ok || breakthrough > 0 {
		c.breakthrough = true
	}

	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4buff = make([]float64, attributes.EndStatType)
	c.a1()

	w, err := minazuki.New(
		minazuki.WithMandatory(keys.Yelan, "yelan burst", burstKey, burstICDKey, 60, c.burstWaveWrapper, c.Core),
		minazuki.WithAnimationDelayCheck(model.AnimationYelanN0StartDelay, c.shouldDelay),
	)
	if err != nil {
		return err
	}
	c.burstWatcher = w
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "breakthrough":
		return c.breakthrough, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 9
	}
	if k == model.AnimationYelanN0StartDelay {
		return 6
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) shouldDelay() bool {
	return c.Core.Player.ActiveChar().NormalCounter == 1
}
