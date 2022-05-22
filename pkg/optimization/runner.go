package optimization

import (
	"github.com/genshinsim/gcsim/internal/simulator"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/result"
	"log"
	"os"
)

type simRunner func(config core.SimulationConfig) result.Summary

func generateSimRunner(cfg string, simopt simulator.Options) simRunner {
	return func(simcfg core.SimulationConfig) result.Summary {
		return runSimWithConfig(cfg, simcfg, simopt)
	}
}

// Just runs the sim with specified settings
func runSimWithConfig(cfg string, simcfg core.SimulationConfig, simopt simulator.Options) result.Summary {
	result, err := simulator.RunWithConfig(cfg, simcfg, simopt)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return result
}
