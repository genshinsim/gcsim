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

	"github.com/genshinsim/gcsim/pkg/enemy"
)

func main() {
	dir := os.Getenv("EXCEL_BIN_OUTPUT")
	if len(dir) == 0 {
		dir = "./GenshinData/ExcelBinOutput"
	}
	dat, err := os.ReadFile(filepath.Join(dir, "MonsterExcelConfigData.json"))
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	parsed := []enemy.MonsterExcelConfig{}
	err = json.Unmarshal(dat, &parsed)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	t, err := template.New("enemystat").Parse(tmplEnemyStats)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	err = t.Execute(writer, parsed)
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

var tmplEnemyStats = `// Code generated DO NOT EDIT.
package enemy

var monsterInfos = map[int]MonsterExcelConfig{
{{range .}} {{.Id}}: {
		MonsterName: "{{.MonsterName}}",
		Typ: "{{.Typ}}",
		HpDrops: []HpDrop{
		{{range .HpDrops}} {{if .DropId}}
			{
				DropId: {{.DropId}},
				HpPercent: {{.HpPercent}},
			},
		{{end}} {{end}}
		},
		KillDropId: {{.KillDropId}},
		HpBase: {{.HpBase}},
		PropGrowCurves: []PropGrowCurve{
		{{range .PropGrowCurves}} {{if .GrowCurve}}
			{
				GrowCurve: "{{.GrowCurve}}",
			},
		{{end}} {{end}}
		},
		FireSubHurt: {{.FireSubHurt}},
		GrassSubHurt: {{.GrassSubHurt}},
		WaterSubHurt: {{.WaterSubHurt}},
		ElecSubHurt: {{.ElecSubHurt}},
		WindSubHurt: {{.WindSubHurt}},
		IceSubHurt: {{.IceSubHurt}},
		RockSubHurt: {{.RockSubHurt}},
		PhysicalSubHurt: {{.PhysicalSubHurt}},
	}, {{end}}
}

`

var tmplNames = `var monsterNameIds = map[string]int{ {{range .}}
"{{.MonsterName}}": {{.Id}},
{{end}}
}
`
