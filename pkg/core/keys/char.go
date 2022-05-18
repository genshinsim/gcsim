package keys

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type Char int

func (c *Char) MarshalJSON() ([]byte, error) {
	return json.Marshal(charNames[*c])
}

func (c *Char) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range charNames {
		if v == s {
			*c = Char(i)
			return nil
		}
	}
	return errors.New("unrecognized character key")
}

func (c Char) String() string {
	return charNames[c]
}

const ChildePassive = "childe-talent-passive"

const (
	NoChar Char = iota
	TravelerElectro
	TravelerAnemo
	TravelerGeo
	TravelerHydro
	TravelerCryo
	TravelerPyro
	TravelerDendro
	TravelerMale
	TravelerFemale
	TravelerDelim // delim
	Albedo
	Aloy
	Amber
	Barbara
	Beidou
	Bennett
	Chongyun
	Diluc
	Diona
	Eula
	Fischl
	Ganyu
	Hutao
	Jean
	Kazuha
	Kaeya
	Ayaka
	Ayato
	Keqing
	Klee
	Sara
	Lisa
	Mona
	Ningguang
	Noelle
	Qiqi
	Raiden
	Razor
	Rosaria
	Kokomi
	Sayu
	Sucrose
	Tartaglia
	Thoma
	Venti
	Xiangling
	Xiao
	Xingqiu
	Xinyan
	Yanfei
	Yoimiya
	Zhongli
	Gorou
	Itto
	Shenhe
	Yunjin
	YaeMiko
)

var charNames = []string{
	"",
	"travelerelectro",
	"traveleranemo",
	"travelergeo",
	"travelerhydro",
	"travelercryo",
	"travelerpyro",
	"travelerdendro",
	"aether",
	"lumine",
	"", //delim for traveler
	"albedo",
	"aloy",
	"amber",
	"barbara",
	"beidou",
	"bennett",
	"chongyun",
	"diluc",
	"diona",
	"eula",
	"fischl",
	"ganyu",
	"hutao",
	"jean",
	"kazuha",
	"kaeya",
	"ayaka",
	"ayato",
	"keqing",
	"klee",
	"sara",
	"lisa",
	"mona",
	"ningguang",
	"noelle",
	"qiqi",
	"raiden",
	"razor",
	"rosaria",
	"kokomi",
	"sayu",
	"sucrose",
	"tartaglia",
	"thoma",
	"venti",
	"xiangling",
	"xiao",
	"xingqiu",
	"xinyan",
	"yanfei",
	"yoimiya",
	"zhongli",
	"gorou",
	"itto",
	"shenhe",
	"yunjin",
	"yaemiko",
}

var CharKeyToEle = map[Char]attributes.Element{
	TravelerElectro: attributes.Electro,
	TravelerAnemo:   attributes.Anemo,
	TravelerGeo:     attributes.Geo,
	TravelerHydro:   attributes.Hydro,
	TravelerCryo:    attributes.Cryo,
	TravelerPyro:    attributes.Pyro,
	TravelerDendro:  attributes.Dendro,
	Albedo:          attributes.Geo,
	Aloy:            attributes.Cryo,
	Amber:           attributes.Pyro,
	Barbara:         attributes.Hydro,
	Beidou:          attributes.Electro,
	Bennett:         attributes.Pyro,
	Chongyun:        attributes.Cryo,
	Diluc:           attributes.Pyro,
	Diona:           attributes.Cryo,
	Eula:            attributes.Cryo,
	Fischl:          attributes.Electro,
	Ganyu:           attributes.Cryo,
	Hutao:           attributes.Pyro,
	Jean:            attributes.Anemo,
	Kazuha:          attributes.Anemo,
	Kaeya:           attributes.Cryo,
	Ayaka:           attributes.Cryo,
	Ayato:           attributes.Hydro,
	Keqing:          attributes.Electro,
	Klee:            attributes.Pyro,
	Sara:            attributes.Electro,
	Lisa:            attributes.Electro,
	Mona:            attributes.Hydro,
	Ningguang:       attributes.Geo,
	Noelle:          attributes.Geo,
	Qiqi:            attributes.Cryo,
	Raiden:          attributes.Electro,
	Razor:           attributes.Electro,
	Rosaria:         attributes.Cryo,
	Kokomi:          attributes.Hydro,
	Sayu:            attributes.Anemo,
	Sucrose:         attributes.Anemo,
	Tartaglia:       attributes.Hydro,
	Thoma:           attributes.Pyro,
	Venti:           attributes.Anemo,
	Xiangling:       attributes.Pyro,
	Xiao:            attributes.Anemo,
	Xingqiu:         attributes.Hydro,
	Xinyan:          attributes.Pyro,
	Yanfei:          attributes.Pyro,
	Yoimiya:         attributes.Pyro,
	Zhongli:         attributes.Geo,
	Gorou:           attributes.Geo,
	Itto:            attributes.Geo,
	Shenhe:          attributes.Cryo,
	Yunjin:          attributes.Geo,
	YaeMiko:         attributes.Electro,
}

var CharNameToKey = map[string]Char{
	"travelerelectro":   TravelerElectro,
	"traveleranemo":     TravelerAnemo,
	"travelergeo":       TravelerGeo,
	"travelerhydro":     TravelerHydro,
	"travelercryo":      TravelerCryo,
	"travelerpyro":      TravelerPyro,
	"travelerdendro":    TravelerDendro,
	"albedo":            Albedo,
	"aloy":              Aloy,
	"amber":             Amber,
	"barbara":           Barbara,
	"beidou":            Beidou,
	"bennett":           Bennett,
	"chongyun":          Chongyun,
	"diluc":             Diluc,
	"diona":             Diona,
	"eula":              Eula,
	"fischl":            Fischl,
	"ganyu":             Ganyu,
	"hutao":             Hutao,
	"jean":              Jean,
	"kaedeharakazuha":   Kazuha,
	"kazuha":            Kazuha,
	"kaeya":             Kaeya,
	"kamisatoayaka":     Ayaka,
	"ayaka":             Ayaka,
	"kamisatoayato":     Ayato,
	"ayato":             Ayato,
	"keqing":            Keqing,
	"klee":              Klee,
	"kujousara":         Sara,
	"kujosara":          Sara,
	"sara":              Sara,
	"lisa":              Lisa,
	"mona":              Mona,
	"ningguang":         Ningguang,
	"noelle":            Noelle,
	"qiqi":              Qiqi,
	"raidenshogun":      Raiden,
	"raiden":            Raiden,
	"razor":             Razor,
	"rosaria":           Rosaria,
	"sangonomiyakokomi": Kokomi,
	"kokomi":            Kokomi,
	"sayu":              Sayu,
	"sucrose":           Sucrose,
	"tartaglia":         Tartaglia,
	"thoma":             Thoma,
	"venti":             Venti,
	"xiangling":         Xiangling,
	"xiao":              Xiao,
	"xingqiu":           Xingqiu,
	"xinyan":            Xinyan,
	"yanfei":            Yanfei,
	"yoimiya":           Yoimiya,
	"yunjin":            Yunjin,
	"zhongli":           Zhongli,
	"gorou":             Gorou,
	"aratakiitto":       Itto,
	"itto":              Itto,
	"shenhe":            Shenhe,
	"yae":               YaeMiko,
	"yaemiko":           YaeMiko,
}
