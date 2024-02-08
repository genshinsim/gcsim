package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var elementMap = template.FuncMap{
	"Pyro":     func() attributes.Element { return attributes.Pyro },
	"Dendro":   func() attributes.Element { return attributes.Dendro },
	"Hydro":    func() attributes.Element { return attributes.Hydro },
	"Electro":  func() attributes.Element { return attributes.Electro },
	"Anemo":    func() attributes.Element { return attributes.Anemo },
	"Cryo":     func() attributes.Element { return attributes.Cryo },
	"Geo":      func() attributes.Element { return attributes.Geo },
	"Physical": func() attributes.Element { return attributes.Physical },
}

var keyRegex = regexp.MustCompile(`\W+`) // for removing spaces

func main() {
	dir := os.Getenv("EXCEL_BIN_OUTPUT")
	if len(dir) == 0 {
		dir = "./GenshinData/ExcelBinOutput"
	}

	var configs []monsterExcelConfig
	var describeConfigs []monsterDescribeExcelConfig
	var curveConfigs []monsterCurveExcelConfig
	var textMap map[string]string

	if err := openConfig(&configs, dir, "MonsterExcelConfigData.json"); err != nil {
		log.Fatalf("err: %v\n", err)
	}
	if err := openConfig(&describeConfigs, dir, "MonsterDescribeExcelConfigData.json"); err != nil {
		log.Fatalf("err: %v\n", err)
	}
	if err := openConfig(&curveConfigs, dir, "MonsterCurveExcelConfigData.json"); err != nil {
		log.Fatalf("err: %v\n", err)
	}
	if err := openConfig(&textMap, dir, "..", "TextMap", "TextMapEN.json"); err != nil {
		log.Fatalf("err: %v\n", err)
	}

	profiles := getEnemyProfiles(configs, describeConfigs, textMap)
	if err := runCodeGen(profiles); err != nil {
		log.Fatalf("err: %v\n", err)
	}

	curves := getGrowthMult(curveConfigs)
	if err := runGrowthMult(curves); err != nil {
		log.Fatalf("err: %v\n", err)
	}

	names := getEnemyNames(profiles)
	if err := runShortcuts(names); err != nil {
		log.Fatalf("err: %v\n", err)
	}
}

func getEnemyProfiles(configs []monsterExcelConfig, describeConfigs []monsterDescribeExcelConfig, textMap map[string]string) []info.EnemyProfile {
	visited := map[string]bool{}
	profiles := make([]info.EnemyProfile, 0, len(configs))
	for i := range configs {
		enemy := &configs[i]
		monsterName := getMonsterName(enemy.DescribeId, describeConfigs, textMap)
		if monsterName == "" { // not valid
			continue
		}
		if enemy.PropGrowCurves[0].GrowCurve == "" { // no hp grow?
			continue
		}

		// is already added
		if _, ok := visited[monsterName]; ok {
			continue
		}

		info := toEnemyProfile(enemy)
		info.MonsterName = monsterName
		profiles = append(profiles, info)
		visited[monsterName] = true
	}
	return profiles
}

func getEnemyNames(profiles []info.EnemyProfile) map[string]int {
	names := make(map[string]int, len(profiles))
	for i := range profiles {
		v := &profiles[i]
		names[v.MonsterName] = v.Id
	}
	return names
}

func runCodeGen(profiles []info.EnemyProfile) error {
	fdata, err := os.Create("enemycurves.go.txt")
	if err != nil {
		return err
	}
	defer fdata.Close()
	tdata, err := template.New("enemycurves").Funcs(elementMap).Funcs(template.FuncMap{
		"curveString": statCurveString,
	}).Parse(tmplEnemyStats)
	if err != nil {
		return err
	}
	return tdata.Execute(fdata, profiles)
}

func runGrowthMult(profiles []map[info.EnemyStatCurve]float64) error {
	fdata, err := os.Create("enemygrowth.go.txt")
	if err != nil {
		return err
	}
	defer fdata.Close()
	tdata, err := template.New("enemygrowth").Funcs(template.FuncMap{
		"curveString": statCurveString,
	}).Parse(tmplGrowth)
	if err != nil {
		return err
	}
	return tdata.Execute(fdata, profiles)
}

func runShortcuts(names map[string]int) error {
	fdata, err := os.Create("enemies.go.txt")
	if err != nil {
		return err
	}
	defer fdata.Close()
	tdata, err := template.New("outshortcuts").Parse(tmplShortcuts)
	if err != nil {
		return err
	}
	return tdata.Execute(fdata, names)
}

func openConfig(result interface{}, path ...string) error {
	jsonFile := filepath.Join(path...)
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	return nil
}

func toEnemyProfile(enemyInfo *monsterExcelConfig) info.EnemyProfile {
	hpGrowCurve := getStatCurve(enemyInfo.PropGrowCurves[0].GrowCurve)
	drops := []info.HpDrop{}
	for _, v := range enemyInfo.HpDrops {
		if v.DropId == 0 || v.HpPercent == 0 {
			continue
		}
		pd := v
		pd.HpPercent /= 100
		drops = append(drops, info.HpDrop{
			DropId:    pd.DropId,
			HpPercent: pd.HpPercent,
		})
	}
	if enemyInfo.KillDropId != 0 {
		// add killDropId as particle drop
		drops = append(drops, info.HpDrop{
			DropId:    enemyInfo.KillDropId,
			HpPercent: 0,
		})
	}
	freezeResist := 0.0 // TODO: dm?
	if enemyInfo.Typ == "MONSTER_BOSS" {
		freezeResist = 1.0
	}
	return info.EnemyProfile{
		Resist: map[attributes.Element]float64{
			attributes.Pyro:     enemyInfo.FireSubHurt,
			attributes.Dendro:   enemyInfo.GrassSubHurt,
			attributes.Hydro:    enemyInfo.WaterSubHurt,
			attributes.Electro:  enemyInfo.ElecSubHurt,
			attributes.Anemo:    enemyInfo.WindSubHurt,
			attributes.Cryo:     enemyInfo.IceSubHurt,
			attributes.Geo:      enemyInfo.RockSubHurt,
			attributes.Physical: enemyInfo.PhysicalSubHurt,
		},
		FreezeResist:  freezeResist,
		ParticleDrops: drops,
		HpBase:        enemyInfo.HpBase,
		HpGrowCurve:   hpGrowCurve,
		Id:            enemyInfo.Id,
		MonsterName:   "",
	}
}

func getMonsterName(describeId int, describeConfigs []monsterDescribeExcelConfig, textMap map[string]string) string {
	for i := range describeConfigs {
		// find the info
		describeInfo := &describeConfigs[i]
		if describeInfo.Id != describeId {
			continue
		}

		// get the key
		text := textMap[strconv.Itoa(describeInfo.NameTextMapHash)]
		if text == "" {
			return ""
		}
		text = keyRegex.ReplaceAllString(text, "")
		return strings.ToLower(text)
	}
	return ""
}

func getGrowthMult(curveConfigs []monsterCurveExcelConfig) []map[info.EnemyStatCurve]float64 {
	growthMult := make([]map[info.EnemyStatCurve]float64, 0)
	for i := range curveConfigs {
		curve := &curveConfigs[i]

		levelGrowth := make(map[info.EnemyStatCurve]float64, 0)
		for j := range curve.CurveInfos {
			info := &curve.CurveInfos[j]

			if stat := getStatCurve(info.Type); stat != 0 {
				levelGrowth[stat] = info.Value
			}
		}
		growthMult = append(growthMult, levelGrowth)
	}
	return growthMult
}

// TODO: protobuf?
func getStatCurve(typ string) info.EnemyStatCurve {
	switch typ {
	case "GROW_CURVE_HP":
		return info.GROW_CURVE_HP
	case "GROW_CURVE_HP_2":
		return info.GROW_CURVE_HP_2
	case "GROW_CURVE_HP_ENVIRONMENT":
		return info.GROW_CURVE_HP_ENVIRONMENT
	}
	return 0
}

func statCurveString(stat info.EnemyStatCurve) string {
	switch stat {
	case info.GROW_CURVE_HP:
		return "GROW_CURVE_HP"
	case info.GROW_CURVE_HP_2:
		return "GROW_CURVE_HP_2"
	case info.GROW_CURVE_HP_ENVIRONMENT:
		return "GROW_CURVE_HP_ENVIRONMENT"
	}
	return ""
}

var tmplEnemyStats = `// Code generated DO NOT EDIT.
package curves

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var EnemyMap = map[int]info.EnemyProfile{
{{- range .}}
	{{.Id}}: {
		Resist: map[attributes.Element]float64{
			attributes.Pyro: {{index .Resist Pyro}},
			attributes.Dendro: {{index .Resist Dendro}},
			attributes.Hydro: {{index .Resist Hydro}},
			attributes.Electro: {{index .Resist Electro}},
			attributes.Anemo: {{index .Resist Anemo}},
			attributes.Cryo: {{index .Resist Cryo}},
			attributes.Geo: {{index .Resist Geo}},
			attributes.Physical: {{index .Resist Physical}},
		},
		ParticleDrops: []info.HpDrop{
		{{- range .ParticleDrops}}
			{{- if .DropId}}
			{
				DropId: {{.DropId}},
				HpPercent: {{.HpPercent}},
			},
			{{- end}}
		{{- end}}
		},
		FreezeResist: {{.FreezeResist}},
		HpBase: {{.HpBase}},
		HpGrowCurve: info.{{curveString .HpGrowCurve}},
		Id: {{.Id}},
		MonsterName: "{{.MonsterName}}",
	},
{{- end}}
}

`

var tmplShortcuts = `package shortcut

var MonsterNameToID = map[string]int{
{{- range $key, $value := .}}
	"{{$key}}": {{$value}},
{{- end}}
}
`

var tmplGrowth = `package curves

import "github.com/genshinsim/gcsim/pkg/core/info"

// EnemyStatGrowthMult provide multiplier for each lvl with 0 being lvl 1
var EnemyStatGrowthMult = []map[info.EnemyStatCurve]float64{
{{- range $key, $value := .}}
	{
	{{- range $stat, $mult := $value}}
		info.{{curveString $stat}}: {{$mult}},
	{{- end}}
	},
{{- end}}
}
`
