package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type base struct {
	HP  float64 `json:"hp"`
	Atk float64 `json:"attack"`
	Def float64 `json:"defense"`
}

type curve struct {
	HP  string `json:"hp"`
	Atk string `json:"attack"`
	Def string `json:"defense"`
}

type promo struct {
	Max         int     `json:"maxlevel"`
	HP          float64 `json:"hp"`
	Atk         float64 `json:"attack"`
	Def         float64 `json:"defense"`
	Specialized float64 `json:"specialized"`
}

type data struct {
	profile
	Key           string
	Base          base    `json:"base"`
	Curve         curve   `json:"curve"`
	Specialized   string  `json:"specialized"`
	PromotionData []promo `json:"promotion"`
}

type profile struct {
	Body       string
	Element    string
	Rarity     string
	Region     string
	WeaponType string
}

func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	b, err := fetch("src/data/stats/characters.json")
	if err != nil {
		return err
	}

	var d map[string]data
	if err := json.Unmarshal([]byte(b), &d); err != nil {
		return err
	}

	// fix the specialized key
	for k, v := range d { //nolint:gocritic // map values are not addressable/can't edit without copying and writing back
		v.Specialized = SpecKeyToStat[v.Specialized]
		v.Key = CharNameToKey[k]

		if v.Key == "" {
			log.Printf("skipping '%v' no valid key\n", k)
			continue
		}

		// fetch char profile
		b, err := fetch(fmt.Sprintf("src/data/English/characters/%s.json", k))
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(b), &v.profile); err != nil {
			return err
		}

		v.Body = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(v.Body))
		// special case for traveler and aloy
		if v.Region == "" {
			v.Region = "Unknown"
		}
		// special case for traveler
		if v.Element == "None" {
			v.Element = "NoElement"
		}
		if v.WeaponType == "Polearm" {
			v.WeaponType = "Spear"
		}

		d[k] = v
		log.Println(v.Key)
	}

	// fmt.Println(d)
	of, err := os.Create("./_output.go")
	if err != nil {
		return err
	}
	defer of.Close()

	t, err := template.New("out").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(of, d)
}

func fetch(path string) (string, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/theBowja/genshin-db/main/" + path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%v: %v", resp.Status, path)
	}

	out, err := io.ReadAll(resp.Body)
	return string(out), err
}

var tmpl = `package curves

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"

	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)


var CharBaseMap = map[keys.Char]CharBase{
	{{- range $key, $value := . }}
	{{- if $value.Key }}
	keys.{{$value.Key}}: {
		Rarity: {{$value.Rarity}},
		Body: profile.Body {{- $value.Body}},
		Element: attributes. {{- $value.Element}},
		Region: profile.Zone {{- $value.Region}},
		WeaponType: info.WeaponClass {{- $value.WeaponType}},
		HPCurve: {{$value.Curve.HP}},
		AtkCurve: {{$value.Curve.Atk}},
		DefCurve: {{$value.Curve.Def}},
		BaseHP: {{$value.Base.HP}},
		BaseAtk: {{$value.Base.Atk}},
		BaseDef: {{$value.Base.Def}},
		Specialized: {{$value.Specialized}},
		PromotionBonus: []PromoData{
			{{- range $e := $value.PromotionData}}
			{
				MaxLevel: {{$e.Max}},
				HP:       {{$e.HP}},
				Atk:      {{$e.Atk}},
				Def:      {{$e.Def}},
				Special:  {{$e.Specialized}},
			},
			{{- end }}
		},
	},
	{{- end }}
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
	"FIGHT_PROP_ELEC_ADD_HURT":     "attributes.ElectroP",
	"FIGHT_PROP_ROCK_ADD_HURT":     "attributes.GeoP",
	"FIGHT_PROP_FIRE_ADD_HURT":     "attributes.PyroP",
	"FIGHT_PROP_WATER_ADD_HURT":    "attributes.HydroP",
	"FIGHT_PROP_DEFENSE_PERCENT":   "attributes.DEFP",
	"FIGHT_PROP_ICE_ADD_HURT":      "attributes.CryoP",
	"FIGHT_PROP_WIND_ADD_HURT":     "attributes.AnemoP",
	"FIGHT_PROP_GRASS_ADD_HURT":    "attributes.DendroP",
}

var CharNameToKey = map[string]string{
	"albedo":            "Albedo",
	"aloy":              "Aloy",
	"amber":             "Amber",
	"barbara":           "Barbara",
	"beidou":            "Beidou",
	"bennett":           "Bennett",
	"chongyun":          "Chongyun",
	"diluc":             "Diluc",
	"diona":             "Diona",
	"eula":              "Eula",
	"fischl":            "Fischl",
	"ganyu":             "Ganyu",
	"hutao":             "Hutao",
	"jean":              "Jean",
	"kaedeharakazuha":   "Kazuha",
	"kazuha":            "Kazuha",
	"kaeya":             "Kaeya",
	"kamisatoayaka":     "Ayaka",
	"ayaka":             "Ayaka",
	"kamisatoayato":     "Ayato",
	"ayato":             "Ayato",
	"keqing":            "Keqing",
	"klee":              "Klee",
	"kujousara":         "Sara",
	"lisa":              "Lisa",
	"mona":              "Mona",
	"ningguang":         "Ningguang",
	"noelle":            "Noelle",
	"qiqi":              "Qiqi",
	"raidenshogun":      "Raiden",
	"raiden":            "Raiden",
	"razor":             "Razor",
	"rosaria":           "Rosaria",
	"sangonomiyakokomi": "Kokomi",
	"kokomi":            "Kokomi",
	"sayu":              "Sayu",
	"sucrose":           "Sucrose",
	"tartaglia":         "Tartaglia",
	"thoma":             "Thoma",
	"venti":             "Venti",
	"xiangling":         "Xiangling",
	"xiao":              "Xiao",
	"xingqiu":           "Xingqiu",
	"xinyan":            "Xinyan",
	"yanfei":            "Yanfei",
	"yoimiya":           "Yoimiya",
	"zhongli":           "Zhongli",
	"gorou":             "Gorou",
	"aratakiitto":       "Itto",
	"aether":            "TravelerMale",
	"lumine":            "TravelerFemale",
	"shenhe":            "Shenhe",
	"yunjin":            "Yunjin",
	"yaemiko":           "YaeMiko",
	"yelan":             "Yelan",
	"kukishinobu":       "Kuki",
	"shikanoinheizou":   "Heizou",
	"tighnari":          "Tighnari",
	"candace":           "Candace",
	"nilou":             "Nilou",
	"alhaitham":         "Alhaitham",
	"layla":             "Layla",
	"wanderer":          "Wanderer",
	"baizhu":            "Baizhu",
	"dehya":             "Dehya",
	"yaoyao":            "Yaoyao",
	"mika":              "Mika",
}
