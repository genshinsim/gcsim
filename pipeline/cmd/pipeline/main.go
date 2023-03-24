package main

import (
	"flag"
	"log"

	"github.com/genshinsim/gcsim/pipeline/pkg/character"
)

type config struct {
	//input data
	charPath  string
	excelPath string

	//output paths
	uiOut string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.charPath, "char", "./internal/characters", "folder to look for character files")
	flag.StringVar(&cfg.excelPath, "excels", "./pipeline/data/ExcelBinOutput", "folder to look for excel data dump")
	flag.StringVar(&cfg.uiOut, "outui", "./ui/packages/ui/src/Data", "folder to output generated json for UI")
	flag.Parse()

	//generate character data
	log.Println("running pipeline for characters...")
	g, err := character.NewGenerator(character.GeneratorConfig{
		Root:   cfg.charPath,
		Excels: cfg.excelPath,
	})
	if err != nil {
		panic(err)
	}

	log.Println("generate character data for ui...")
	err = g.DumpUIJSON(cfg.uiOut)
	if err != nil {
		panic(err)
	}

}
