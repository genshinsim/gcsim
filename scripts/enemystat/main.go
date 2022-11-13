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
	"github.com/genshinsim/gcsim/pkg/enemy"
)

type propGrowCurve struct {
	GrowCurve string `json:"growCurve"`
}

type monsterExcelConfig struct {
	MonsterName     string          `json:"monsterName"`
	Typ             string          `json:"type"`
	HpDrops         []enemy.HpDrop  `json:"hpDrops"`
	KillDropId      int             `json:"killDropId"`
	HpBase          float64         `json:"hpBase"`
	PropGrowCurves  []propGrowCurve `json:"propGrowCurves"`
	FireSubHurt     float64         `json:"fireSubHurt"`
	GrassSubHurt    float64         `json:"grassSubHurt"`
	WaterSubHurt    float64         `json:"waterSubHurt"`
	ElecSubHurt     float64         `json:"elecSubHurt"`
	WindSubHurt     float64         `json:"windSubHurt"`
	IceSubHurt      float64         `json:"iceSubHurt"`
	RockSubHurt     float64         `json:"rockSubHurt"`
	PhysicalSubHurt float64         `json:"physicalSubHurt"`
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
	dir := os.Getenv("EXCEL_BIN_OUTPUT")
	if len(dir) == 0 {
		dir = "./GenshinData/ExcelBinOutput"
	}
	dat, err := os.ReadFile(filepath.Join(dir, "MonsterExcelConfigData.json"))
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	parsed := []monsterExcelConfig{}
	err = json.Unmarshal(dat, &parsed)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	profiles := make([]enemy.EnemyProfile, len(parsed))
	for i, v := range parsed {
		profiles[i] = toEnemyProfile(v)
	}
	t, err := template.New("enemystat").Funcs(elementMap).Parse(tmplEnemyStats)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	err = t.Execute(writer, profiles)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	writer.Flush()
	visited := map[string]bool{}
	k := 0
	for _, v := range parsed {
		if _, ok := visited[v.MonsterName]; !ok {
			visited[v.MonsterName] = true
			parsed[k] = v
			k++
		}
	}
	parsed = parsed[:k]
	t, err = template.New("enemyname").Parse(tmplNames)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	err = t.Execute(writer, parsed)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	writer.Flush()
	content, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	err = os.WriteFile("enemystat.go", content, 0644)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
}

func toEnemyProfile(info monsterExcelConfig) enemy.EnemyProfile {
	hpGrowCurve := 1
	if info.PropGrowCurves[0].GrowCurve == "GROW_CURVE_HP_2" {
		hpGrowCurve = 2
	}
	drops := []enemy.HpDrop{}
	for _, v := range info.HpDrops {
		if v.DropId == 0 || v.HpPercent == 0 {
			continue
		}
		pd := v
		pd.HpPercent /= 100
		drops = append(drops, pd)
	}
	if info.KillDropId != 0 {
		// add killDropId as particle drop
		drops = append(drops, enemy.HpDrop{
			DropId:    info.KillDropId,
			HpPercent: 0,
		})
	}
	return enemy.EnemyProfile{
		Resist: map[attributes.Element]float64{
			attributes.Pyro:     info.FireSubHurt,
			attributes.Dendro:   info.GrassSubHurt,
			attributes.Hydro:    info.WaterSubHurt,
			attributes.Electro:  info.ElecSubHurt,
			attributes.Anemo:    info.WindSubHurt,
			attributes.Cryo:     info.IceSubHurt,
			attributes.Geo:      info.RockSubHurt,
			attributes.Physical: info.PhysicalSubHurt,
		},
		ParticleDrops: drops,
		ResistFrozen:  info.Typ == "MONSTER_BOSS",
		HpBase:        info.HpBase,
		HpGrowCurve:   hpGrowCurve,
		Id:            info.Id,
		MonsterName:   info.MonsterName,
	}
}

var tmplEnemyStats = `// Code generated DO NOT EDIT.
package enemy

import "github.com/genshinsim/gcsim/pkg/core/attributes"

var monsterInfos = map[int]EnemyProfile{
{{range .}} {{.Id}}: {
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
		ParticleDrops: []HpDrop{
		{{range .ParticleDrops}} {{if .DropId}} {
				DropId: {{.DropId}},
				HpPercent: {{.HpPercent}},
			},
		{{end}} {{end}}
		},
		ResistFrozen: {{.ResistFrozen}},
		HpBase: {{.HpBase}},
		HpGrowCurve: {{.HpGrowCurve}},
		Id: {{.Id}},
		MonsterName : "{{.MonsterName}}",
	}, {{end}}
}

`

var tmplNames = `var monsterNameIds = map[string]int{ {{range .}}
"{{.MonsterName}}": {{.Id}}, {{end}}
}
`
