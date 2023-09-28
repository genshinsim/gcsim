package main

import (
	"flag"
	"log"

	"github.com/genshinsim/gcsim/pipeline/pkg/artifact"
	"github.com/genshinsim/gcsim/pipeline/pkg/character"
	"github.com/genshinsim/gcsim/pipeline/pkg/weapon"
)

type config struct {
	// input data
	charPath     string
	weapPath     string
	artifactPath string
	excelPath    string

	// output paths
	uiOut string
	dbOut string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.charPath, "char", "./internal/characters", "folder to look for character files")
	flag.StringVar(&cfg.weapPath, "weap", "./internal/weapons", "folder to look for weapon files")
	flag.StringVar(&cfg.artifactPath, "art", "./internal/artifacts", "folder to look for artifact files")
	flag.StringVar(&cfg.excelPath, "excels", "./pipeline/data/ExcelBinOutput", "folder to look for excel data dump")
	flag.StringVar(&cfg.uiOut, "outui", "./ui/packages/ui/src/Data", "folder to output generated json for UI")
	flag.StringVar(&cfg.dbOut, "outdb", "./ui/packages/db/src/Data", "folder to output generated json for DB")
	flag.Parse()

	// generate character data
	log.Println("running pipeline for characters...")
	g, err := character.NewGenerator(character.GeneratorConfig{
		Root:   cfg.charPath,
		Excels: cfg.excelPath,
	})
	if err != nil {
		panic(err)
	}

	log.Println("generate character data for ui...")
	err = g.DumpJSON(cfg.uiOut)
	if err != nil {
		panic(err)
	}

	log.Println("generate character data for db...")
	err = g.DumpJSON(cfg.dbOut)
	if err != nil {
		panic(err)
	}

	// generate weapon data
	log.Println("running pipeline for weapons...")
	gw, err := weapon.NewGenerator(weapon.GeneratorConfig{
		Root:   cfg.weapPath,
		Excels: cfg.excelPath,
	})
	if err != nil {
		panic(err)
	}

	log.Println("generate weapon data for ui...")
	err = gw.DumpUIJSON(cfg.uiOut)
	if err != nil {
		panic(err)
	}

	// generate artifact data
	log.Println("running pipeline for artifacts...")
	ga, err := artifact.NewGenerator(artifact.GeneratorConfig{
		Root: cfg.artifactPath,
	})
	if err != nil {
		panic(err)
	}

	log.Println("generate artifacts data for ui...")
	err = ga.DumpJSON(cfg.uiOut)
	if err != nil {
		panic(err)
	}
}
