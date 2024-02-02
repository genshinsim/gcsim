package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type propGrowCurve struct {
	GrowCurve string `json:"grow_curve"`
}

type monsterExcelConfig struct {
	MonsterName     string          `json:"monster_name"`
	Typ             string          `json:"type"`
	HpDrops         []info.HpDrop   `json:"hp_drops"`
	KillDropId      int             `json:"kill_drop_id"`
	HpBase          float64         `json:"hp_base"`
	PropGrowCurves  []propGrowCurve `json:"prop_grow_curves"`
	FireSubHurt     float64         `json:"fire_sub_hurt"`
	GrassSubHurt    float64         `json:"grass_sub_hurt"`
	WaterSubHurt    float64         `json:"water_sub_hurt"`
	ElecSubHurt     float64         `json:"elec_sub_hurt"`
	WindSubHurt     float64         `json:"wind_sub_hurt"`
	IceSubHurt      float64         `json:"ice_sub_hurt"`
	RockSubHurt     float64         `json:"rock_sub_hurt"`
	PhysicalSubHurt float64         `json:"physical_sub_hurt"`
	Id              int             `json:"id"`
}

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

func main() {
	if err := runCodeGen(); err != nil {
		log.Fatalf("err: %v\n", err)
	}
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
		drops = append(drops, pd)
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
		MonsterName:   enemyInfo.MonsterName,
	}
}

func runCodeGen() error {
	configs, err := readMonsterConfigs()
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	err = writeMonsterInfo(configs, writer)
	if err != nil {
		return err
	}
	err = writeNameIds(configs, writer)
	if err != nil {
		return err
	}
	writer.Flush()
	content, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	err = os.WriteFile("enemystat.go", content, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func writeMonsterInfo(configs []monsterExcelConfig, out *bufio.Writer) error {
	profiles := make([]info.EnemyProfile, len(configs))
	for i := range configs {
		profiles[i] = toEnemyProfile(&configs[i])
	}
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

func writeNameIds(configs []monsterExcelConfig, out *bufio.Writer) error {
	used := []*monsterExcelConfig{}
	visited := map[string]bool{}
	for i := range configs {
		v := &configs[i]
		if _, ok := visited[v.MonsterName]; !ok {
			visited[v.MonsterName] = true
			used = append(used, v)
		}
	}
	t, err := template.New("enemyname").Parse(tmplNames)
	if err != nil {
		return err
	}
	err = t.Execute(out, used)
	if err != nil {
		return err
	}
	return nil
}

func readMonsterConfigs() ([]monsterExcelConfig, error) {
	dir := os.Getenv("EXCEL_BIN_OUTPUT")
	if len(dir) == 0 {
		dir = "./GenshinData/ExcelBinOutput"
	}
	dat, err := os.ReadFile(filepath.Join(dir, "MonsterExcelConfigData.json"))
	if err != nil {
		return nil, err
	}
	parsed := []monsterExcelConfig{}
	err = json.Unmarshal(dat, &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

var tmplEnemyStats = `// Code generated DO NOT EDIT.
package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var monsterInfos = map[int]info.EnemyProfile{
{{range .}}
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
{{end}}
}

`

var tmplNames = `var monsterNameIds = map[string]int{ {{range .}}
"{{.MonsterName}}": {{.Id}}, {{end}}
}
`
