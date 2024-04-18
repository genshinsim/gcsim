package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/genshinsim/gcsim/pipeline/pkg/artifact"
	"github.com/genshinsim/gcsim/pipeline/pkg/character"
	"github.com/genshinsim/gcsim/pipeline/pkg/enemy"
	"github.com/genshinsim/gcsim/pipeline/pkg/translation"
	"github.com/genshinsim/gcsim/pipeline/pkg/weapon"
	"github.com/genshinsim/gcsim/pkg/model"
)

type config struct {
	// input data
	charPath     string
	weapPath     string
	artifactPath string
	excelPath    string

	// what to run
	runChar      bool
	runWeap      bool
	runArtifacts bool
	runEnemies   bool

	// output paths
	pkgOut      string
	uiOut       string
	dbOut       string
	transOut    string
	icdPath     string
	keyPath     string
	importsPath string
	docRoot     string
	assetsRoot  string
}

func main() {
	var cfg config
	flag.BoolVar(&cfg.runChar, "gen-char", false, "generate char data")
	flag.BoolVar(&cfg.runWeap, "gen-weap", false, "generate weapon data")
	flag.BoolVar(&cfg.runArtifacts, "gen-art", false, "generate artifacts data")
	flag.BoolVar(&cfg.runEnemies, "gen-enemies", false, "generate enemies data")
	flag.StringVar(&cfg.charPath, "char", "./internal/characters", "folder to look for character files")
	flag.StringVar(&cfg.weapPath, "weap", "./internal/weapons", "folder to look for weapon files")
	flag.StringVar(&cfg.artifactPath, "art", "./internal/artifacts", "folder to look for artifact files")
	flag.StringVar(&cfg.excelPath, "excels", "", "folder to look for excel data dump")
	flag.StringVar(&cfg.pkgOut, "outpkg", "./pkg", "for to output generated go files to pkg")
	flag.StringVar(&cfg.uiOut, "outui", "./ui/packages/ui/src/Data", "folder to output generated json for UI")
	flag.StringVar(&cfg.dbOut, "outdb", "./ui/packages/db/src/Data", "folder to output generated json for DB")
	flag.StringVar(&cfg.transOut, "outtrans", "./ui/packages/localization/src/locales", "folder to output generated json for DB")
	flag.StringVar(&cfg.icdPath, "icd", "./pkg/core/attacks", "file to store generated icd data")
	flag.StringVar(&cfg.keyPath, "keys", "./pkg/core/keys", "path to store generated keys data")
	flag.StringVar(&cfg.importsPath, "imports", "./pkg/simulation", "path to store generated imports data")
	flag.StringVar(&cfg.docRoot, "outdocs", "./ui/packages/docs/src/components", "file to store generated icd data")
	flag.StringVar(&cfg.assetsRoot, "assets", "./internal/services/assets", "path to store generate asset data")
	flag.Parse()

	// try env first
	if cfg.excelPath == "" {
		cfg.excelPath = os.Getenv("GENSHIN_DATA_REPO")
	}
	// if not, use default
	if cfg.excelPath == "" {
		cfg.excelPath = "./pipeline/data"
	}

	excels := filepath.Join(cfg.excelPath, "ExcelBinOutput")

	var (
		charData      []*model.AvatarData
		weapData      []*model.WeaponData
		artifactsData []*model.ArtifactData
		monsterData   []*model.MonsterData
	)

	// generate character data
	if cfg.runChar {
		charData = genChar(cfg, excels)
	}

	// generate weapon data
	if cfg.runWeap {
		weapData = genWeap(cfg, excels)
	}

	// generate artifact data
	if cfg.runArtifacts {
		artifactsData = genArtifacts(cfg, excels)
	}

	// generate enemy data
	if cfg.runEnemies {
		monsterData = genEnemies(cfg, excels)
	}

	// generate translation data

	if cfg.runChar && cfg.runArtifacts && cfg.runEnemies && cfg.runWeap {
		transCfg := translation.GeneratorConfig{
			Characters: charData,
			Weapons:    weapData,
			Artifacts:  artifactsData,
			Enemies:    monsterData,
			Languages: map[string]string{
				"English":  filepath.Join(cfg.excelPath, "TextMap", "TextMapEN.json"),
				"Chinese":  filepath.Join(cfg.excelPath, "TextMap", "TextMapCHS.json"),
				"Japanese": filepath.Join(cfg.excelPath, "TextMap", "TextMapJP.json"),
				"Korean":   filepath.Join(cfg.excelPath, "TextMap", "TextMapKR.json"),
				"Spanish":  filepath.Join(cfg.excelPath, "TextMap", "TextMapES.json"),
				"Russian":  filepath.Join(cfg.excelPath, "TextMap", "TextMapRU.json"),
				"German":   filepath.Join(cfg.excelPath, "TextMap", "TextMapDE.json"),
			},
		}
		ts, err := translation.NewGenerator(transCfg)
		if err != nil {
			panic(err)
		}
		err = ts.DumpUIJSON(cfg.transOut)
		if err != nil {
			panic(err)
		}
	}
}

func genChar(cfg config, excels string) []*model.AvatarData {
	log.Println("running pipeline for characters...")
	g, err := character.NewGenerator(character.GeneratorConfig{
		Root:   cfg.charPath,
		Excels: excels,
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

	log.Println("generate character template data...")
	err = g.GenerateCharTemplate()
	if err != nil {
		panic(err)
	}

	log.Println("generate character icd data...")
	err = g.GenerateICDData(cfg.icdPath)
	if err != nil {
		panic(err)
	}

	log.Println("generate char keys data...")
	err = g.GenerateKeys(cfg.keyPath)
	if err != nil {
		panic(err)
	}

	log.Println("generate char imports data...")
	err = g.GenerateImports(cfg.importsPath)
	if err != nil {
		panic(err)
	}

	log.Println("generate character assets mapping...")
	err = g.GenerateAssetsKey(cfg.assetsRoot)
	if err != nil {
		panic(err)
	}

	// documentation
	log.Println("generate character documentation data...")
	err = g.WriteFieldDocs(filepath.Join(cfg.docRoot, "/Fields/character_data.json"))
	if err != nil {
		panic(err)
	}

	return g.Data()
}

func genWeap(cfg config, excels string) []*model.WeaponData {
	log.Println("running pipeline for weapons...")
	g, err := weapon.NewGenerator(weapon.GeneratorConfig{
		Root:   cfg.weapPath,
		Excels: excels,
	})
	if err != nil {
		panic(err)
	}

	log.Println("generate weapon data for ui...")
	err = g.DumpUIJSON(cfg.uiOut)
	if err != nil {
		panic(err)
	}

	log.Println("generate weapon template data...")
	err = g.GenerateTemplate()
	if err != nil {
		panic(err)
	}

	log.Println("generate character assets mapping...")
	err = g.GenerateAssetsKey(cfg.assetsRoot)
	if err != nil {
		panic(err)
	}

	return g.Data()
}

func genArtifacts(cfg config, excels string) []*model.ArtifactData {
	log.Println("running pipeline for artifacts...")
	g, err := artifact.NewGenerator(artifact.GeneratorConfig{
		Root:   cfg.artifactPath,
		Excels: excels,
	})
	if err != nil {
		panic(err)
	}

	log.Println("generate artifact assets mapping...")
	err = g.GenerateAssetsKey(cfg.assetsRoot)
	if err != nil {
		panic(err)
	}

	log.Println("generate artifacts data for ui...")
	err = g.DumpJSON(cfg.uiOut)
	if err != nil {
		panic(err)
	}
	return g.Data()
}

func genEnemies(cfg config, excels string) []*model.MonsterData {
	log.Println("running pipeline for enemies...")
	g, err := enemy.NewGenerator(enemy.GeneratorConfig{
		Root:   cfg.charPath,
		Excels: excels,
	})
	if err != nil {
		panic(err)
	}
	err = g.GenerateEnemyStats(cfg.pkgOut)
	if err != nil {
		panic(err)
	}
	err = g.GenerateEnemyShortcuts(cfg.pkgOut)
	if err != nil {
		panic(err)
	}
	return g.Data()
}
