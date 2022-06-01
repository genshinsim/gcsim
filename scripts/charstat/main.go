package main

import (
	"encoding/json"
	"log"
	"os"
	"text/template"
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
	Key           string
	Base          base    `json:"base"`
	Curve         curve   `json:"curve"`
	Specialized   string  `json:"specialized"`
	PromotionData []promo `json:"promotion"`
}

func main() {

	f, err := os.ReadFile("./characters.json")
	if err != nil {
		log.Panic(err)
	}
	var d map[string]data
	err = json.Unmarshal(f, &d)
	if err != nil {
		log.Panic(err)
	}
	//fix the specialized key
	for k, v := range d {
		v.Specialized = SpecKeyToStat[v.Specialized]
		v.Key = CharNameToKey[k]
		d[k] = v
		// fmt.Println(k)
	}
	// fmt.Println(d)
	t, err := template.New("out").Parse(tmpl)
	if err != nil {
		log.Panic(err)
	}
	os.Remove("./char.txt")
	of, err := os.Create("./char.txt")
	if err != nil {
		log.Panic(err)
	}
	err = t.Execute(of, d)
	if err != nil {
		log.Panic(err)
	}
}

var tmpl = `package curves

import (
	"github.com/genshinsim/gcsim/pkg/core"
)


var CharBaseMap = map[core.CharKey]CharBase{
	{{- range $key, $value := . }}
	core.{{$value.Key}}: {
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
}

`

var SpecKeyToStat = map[string]string{
	"FIGHT_PROP_CRITICAL_HURT":     "core.CD",
	"FIGHT_PROP_HEAL_ADD":          "core.Heal",
	"FIGHT_PROP_ATTACK_PERCENT":    "core.ATKP",
	"FIGHT_PROP_ELEMENT_MASTERY":   "core.EM",
	"FIGHT_PROP_HP_PERCENT":        "core.HPP",
	"FIGHT_PROP_CHARGE_EFFICIENCY": "core.ER",
	"FIGHT_PROP_CRITICAL":          "core.CR",
	"FIGHT_PROP_PHYSICAL_ADD_HURT": "core.PhyP",
	"FIGHT_PROP_ELEC_ADD_HURT":     "core.ElectroP",
	"FIGHT_PROP_ROCK_ADD_HURT":     "core.GeoP",
	"FIGHT_PROP_FIRE_ADD_HURT":     "core.PyroP",
	"FIGHT_PROP_WATER_ADD_HURT":    "core.HydroP",
	"FIGHT_PROP_DEFENSE_PERCENT":   "core.DEFP",
	"FIGHT_PROP_ICE_ADD_HURT":      "core.CryoP",
	"FIGHT_PROP_WIND_ADD_HURT":     "core.AnemoP",
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
}
