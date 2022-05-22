package optimization

import (
	"log"
	"os"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/pkg/simulator"
)

type simRunner func(config *ast.ActionList) result.Summary

func generateSimRunner(cfg string, simopt simulator.Options) simRunner {
	return func(simcfg *ast.ActionList) result.Summary {
		return runSimWithConfig(cfg, simcfg, simopt)
	}
}

// Just runs the sim with specified settings
func runSimWithConfig(cfg string, simcfg *ast.ActionList, simopt simulator.Options) result.Summary {
	result, err := simulator.RunWithConfig(cfg, simcfg, simopt)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return result
}
