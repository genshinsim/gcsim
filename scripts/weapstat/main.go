package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type stats struct {
	Max         int     `json:"maxlevel"`
	HP          float64 `json:"hp"`
	Atk         float64 `json:"attack"`
	Def         float64 `json:"defense"`
	Specialized float64 `json:"specialized"`
}

type curve struct {
	Atk         string `json:"attack"`
	Specialized string `json:"specialized"`
}

type data struct {
	Base          stats   `json:"base"`
	Curve         curve   `json:"curve"`
	Specialized   string  `json:"specialized"`
	PromotionData []stats `json:"promotion"`
	TitleCase     string
}

func main() {
	//https://raw.githubusercontent.com/theBowja/genshin-db/main/src/data/stats/weapons.json -> weapons.json
	//https://raw.githubusercontent.com/theBowja/genshin-db/main/src/data/index/English/weapons.json -> weapons_names.json
	os.Remove("./weapons.json")
	err := download("./weapons.json", "https://raw.githubusercontent.com/theBowja/genshin-db/main/src/data/stats/weapons.json")
	if err != nil {
		panic(err)
	}
	os.Remove("./weapons_names.json")
	err = download("./weapons_names.json", "https://raw.githubusercontent.com/theBowja/genshin-db/main/src/data/index/English/weapons.json")
	if err != nil {
		panic(err)
	}

	names := readNameMap()

	f, err := os.ReadFile("./weapons.json")
	if err != nil {
		log.Panic(err)
	}
	var d map[string]data
	err = json.Unmarshal(f, &d)
	if err != nil {
		log.Panic(err)
	}
	// fix the specialized key
	for k, v := range d { //nolint:gocritic
		v.Specialized = SpecKeyToStat[v.Specialized]
		if v.Specialized == "" {
			v.Specialized = "attributes.NoStat"
		}
		v.TitleCase = names[k]
		d[k] = v
		// fmt.Println(k)
	}
	// fmt.Println(d)
	writeTmpl(tmpl, "./weap.txt", d)
	writeTmpl(tmplKeys, "./keys.txt", d)
	writeTmpl(tmplShortcuts, "./shortcuts.txt", d)
}

func download(filepath, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func writeTmpl(tmplStr, outFile string, d map[string]data) {
	t, err := template.New("out").Parse(tmplStr)
	if err != nil {
		log.Panic(err)
	}
	os.Remove(outFile)
	of, err := os.Create(outFile)
	if err != nil {
		log.Panic(err)
	}
	err = t.Execute(of, d)
	if err != nil {
		log.Panic(err)
	}
}

type namemap struct {
	Names map[string]string `json:"namemap"`
}

var re = regexp.MustCompile(`(?i)[^0-9a-z]`)

func readNameMap() map[string]string {
	f, err := os.ReadFile("./weapons_names.json")
	if err != nil {
		log.Panic(err)
	}
	var m namemap
	err = json.Unmarshal(f, &m)
	if err != nil {
		log.Panic(err)
	}

	for k, v := range m.Names {
		// strip out any none word characters
		v = strings.ReplaceAll(v, "'", "")
		m.Names[k] = re.ReplaceAllString(cases.Title(language.Und, cases.NoLower).String(v), "")
	}
	return m.Names
}

var tmpl = `package curves

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)


var WeaponBaseMap = map[keys.Weapon]WeaponBase{
	{{- range $key, $value := . }}
	keys.{{$value.TitleCase}}: {
		AtkCurve: {{$value.Curve.Atk}},
		SpecializedCurve: {{$value.Curve.Specialized}},
		BaseAtk: {{$value.Base.Atk}},
		BaseSpecialized: {{$value.Base.Specialized}},
		Specialized: {{$value.Specialized}},
		PromotionBonus: []PromoData{
			{{- range $e := $value.PromotionData}}
			{
				MaxLevel: {{$e.Max}},
				Atk:      {{$e.Atk}},
			},
			{{- end }}
		},
	},
	{{- end }}
}

`

var tmplKeys = `package keys

import (
	"encoding/json"
	"errors"
	"strings"
)

type Weapon int

func (c *Weapon) MarshalJSON() ([]byte, error) {
	return json.Marshal(weaponNames[*c])
}

func (c *Weapon) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range weaponNames {
		if v == s {
			*c = Weapon(i)
			return nil
		}
	}
	return errors.New("unrecognized weapon key")
}

func (c Weapon) String() string {
	return weaponNames[c]
}


var weaponNames = []string{
	"",
	{{- range $key, $value := . }}
	"{{$key}}",
	{{- end }}
}

const (
	NoWeapon Weapon = iota
	{{- range $key, $value := . }}
	{{$value.TitleCase}}
	{{- end }}
)
`

var tmplShortcuts = `package shortcut

import "github.com/genshinsim/gcsim/pkg/core/keys"

var WeaponNameToKey = map[string]keys.Weapon{
	{{- range $key, $value := . }}
	"{{$key}}": keys.{{$value.TitleCase}},
	{{- end }}
}
`

var SpecKeyToStat = map[string]string{
	"FIGHT_PROP_CRITICAL_HURT":     "attributes.CD",
	"FIGHT_PROP_HEAL_ADD":          "attributes.Heal",
	"FIGHT_PROP_ATTACK_PERCENT":    "attributes.ATKP",
	"FIGHT_PROP_ELEMENT_MASTERY":   "attributes.EM",
	"FIGHT_PROP_HP_PERCENT":        "attributes.HPP",
	"FIGHT_PROP_CHARGE_EFFICIENCY": "attributes.ER",
	"FIGHT_PROP_CRITICAL":          "attributes.CR",
	"FIGHT_PROP_PHYSICAL_ADD_HURT": "attributes.PhyP",
	"FIGHT_PROP_ELEC_ADD_HURT":     "attributes.EleP",
	"FIGHT_PROP_ROCK_ADD_HURT":     "attributes.GeoP",
	"FIGHT_PROP_FIRE_ADD_HURT":     "attributes.PyroP",
	"FIGHT_PROP_WATER_ADD_HURT":    "attributes.HydroP",
	"FIGHT_PROP_DEFENSE_PERCENT":   "attributes.DEFP",
	"FIGHT_PROP_ICE_ADD_HURT":      "attributes.CryoP",
	"FIGHT_PROP_WIND_ADD_HURT":     "attributes.AnemoP",
}
