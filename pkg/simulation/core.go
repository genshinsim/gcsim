package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func NewCore(seed int64, debug bool, cfg *info.ActionList) (*core.Core, error) {
	return core.New(core.Opt{
		Seed:            seed,
		Debug:           debug,
		Delays:          cfg.Settings.Delays,
		DefHalt:         cfg.Settings.DefHalt,
		DamageMode:      cfg.Settings.DamageMode,
		EnableHitlag:    cfg.Settings.EnableHitlag,
		ErCalc:          cfg.Settings.ErCalc,
		ExpectedCritDmg: cfg.Settings.ExpectedCritDmg,
	})
}
