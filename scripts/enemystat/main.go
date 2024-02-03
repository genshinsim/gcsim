package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"go/format"
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
	var textMap map[string]string

	if err := openConfig(&configs, dir, "MonsterExcelConfigData.json"); err != nil {
		log.Fatalf("err: %v\n", err)
	}
	if err := openConfig(&describeConfigs, dir, "MonsterDescribeExcelConfigData.json"); err != nil {
		log.Fatalf("err: %v\n", err)
	}
	if err := openConfig(&textMap, dir, "..", "TextMap", "TextMapEN.json"); err != nil {
		log.Fatalf("err: %v\n", err)
	}

	profiles := getEnemyProfiles(configs, describeConfigs, textMap)
	if err := runCodeGen(profiles); err != nil {
		log.Fatalf("err: %v\n", err)
	}

	names := getEnemyNames(profiles)
	if err := runShortcuts(names); err != nil {
		log.Fatalf("err: %v\n", err)
	}
}

func getEnemyProfiles(configs []monsterExcelConfig, describeConfigs []monsterDescribeExcelConfig, textMap map[string]string) []info.EnemyProfile {
	profiles := make([]info.EnemyProfile, 0, len(configs))
	for i := range configs {
		monsterName := getMonsterName(configs[i].DescribeId, describeConfigs, textMap)
		if monsterName == "" { // not valid
			continue
		}

		info := toEnemyProfile(&configs[i])
		info.MonsterName = monsterName
		profiles = append(profiles, info)
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
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	if err := writeMonsterInfo(profiles, writer); err != nil {
		return err
	}
	writer.Flush()
	content, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	err = os.WriteFile("enemystat.go.txt", content, 0o644)
	if err != nil {
		return err
	}
	return nil
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
	if err := tdata.Execute(fdata, names); err != nil {
		return err
	}
	return nil
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
	hpGrowCurve := 1
	if enemyInfo.PropGrowCurves[0].GrowCurve == "GROW_CURVE_HP_2" {
		hpGrowCurve = 2
	}
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

func writeMonsterInfo(profiles []info.EnemyProfile, out *bufio.Writer) error {
	t, err := template.New("enemystat").Funcs(elementMap).Parse(tmplEnemyStats)
	if err != nil {
		return err
	}
	err = t.Execute(out, profiles)
	if err != nil {
		return err
	}
	return nil
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

var tmplEnemyStats = `// Code generated DO NOT EDIT.
package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var monsterInfos = map[int]info.EnemyProfile{
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
		{{range .ParticleDrops}} {{if .DropId}} {
				DropId: {{.DropId}},
				HpPercent: {{.HpPercent}},
			},
		{{end}} {{end}}
		},
		FreezeResist: {{.FreezeResist}},
		HpBase: {{.HpBase}},
		HpGrowCurve: {{.HpGrowCurve}},
		Id: {{.Id}},
		MonsterName : "{{.MonsterName}}",
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
