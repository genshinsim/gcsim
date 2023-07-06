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

func (c Char) Pretty() string {
	return charPrettyName[c]
}

const ChildePassive = "childe-talent-passive"

const (
	NoChar Char = iota
	AetherAnemo
	LumineAnemo
	AetherGeo
	LumineGeo
	AetherElectro
	LumineElectro
	AetherDendro
	LumineDendro
	AetherHydro
	LumineHydro
	AetherPyro
	LuminePyro
	AetherCryo
	LumineCryo
	Aether
	Lumine
	TravelerDelim // delim
	Albedo
	Aloy
	Amber
	Barbara
	Beidou
	Bennett
	Chongyun
	Cyno
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
	Kirara
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
	Yelan
	Kuki
	Heizou
	Tighnari
	Collei
	Dori
	Candace
	Nilou
	Nahida
	Alhaitham
	Layla
	Faruzan
	Wanderer
	Baizhu
	Dehya
	Yaoyao
	Mika
	Kaveh
	Freminet
	TestCharDoNotUse
	EndCharKeys
)

var charNames = []string{
	"",
	"aetheranemo",
	"lumineanemo",
	"aethergeo",
	"luminegeo",
	"aetherelectro",
	"lumineelectro",
	"aetherdendro",
	"luminedendro",
	"aetherhydro",
	"luminehydro",
	"aetherpyro",
	"luminepyro",
	"aethercryo",
	"luminecryo",
	"aether",
	"lumine",
	"", // delim for traveler
	"albedo",
	"aloy",
	"amber",
	"barbara",
	"beidou",
	"bennett",
	"chongyun",
	"cyno",
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
	"kirara",
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
	"yelan",
	"kuki",
	"heizou",
	"tighnari",
	"collei",
	"dori",
	"candace",
	"nilou",
	"nahida",
	"alhaitham",
	"layla",
	"faruzan",
	"wanderer",
	"baizhu",
	"dehya",
	"yaoyao",
	"mika",
	"kaveh",
	"freminet",
	"test_char_do_not_use",
}

var charPrettyName = []string{
	"Invalid",
	"Aether (Anemo)",
	"Lumine (Anemo)",
	"Aether (Geo)",
	"Lumine (Geo)",
	"Aether (Electro)",
	"Lumine (Electro)",
	"Aether (Dendro)",
	"Lumine (Dendro)",
	"Aether (Hydro)",
	"Lumine (Hydro)",
	"Aether (Pyro)",
	"Lumine (Pyro)",
	"Aether (Cryo)",
	"Lumine (Cryo)",
	"Aether",
	"Lumine",
	"Invalid",
	"Albedo",
	"Aloy",
	"Amber",
	"Barbara",
	"Beidou",
	"Bennett",
	"Chongyun",
	"Cyno",
	"Diluc",
	"Diona",
	"Eula",
	"Fischl",
	"Ganyu",
	"Hutao",
	"Jean",
	"Kazuha",
	"Kaeya",
	"Ayaka",
	"Ayato",
	"Keqing",
	"Kirara",
	"Klee",
	"Sara",
	"Lisa",
	"Mona",
	"Ningguang",
	"Noelle",
	"Qiqi",
	"Raiden",
	"Razor",
	"Rosaria",
	"Kokomi",
	"Sayu",
	"Sucrose",
	"Tartaglia",
	"Thoma",
	"Venti",
	"Xiangling",
	"Xiao",
	"Xingqiu",
	"Xinyan",
	"Yanfei",
	"Yoimiya",
	"Zhongli",
	"Gorou",
	"Itto",
	"Shenhe",
	"Yunjin",
	"Yae Miko",
	"Yelan",
	"Kuki",
	"Heizou",
	"Tighnari",
	"Collei",
	"Dori",
	"Candace",
	"Nilou",
	"Nahida",
	"Alhaitham",
	"Layla",
	"Faruzan",
	"Wanderer",
	"Baizhu",
	"Dehya",
	"Yaoyao",
	"Mika",
	"Kaveh",
	"Freminet",
	"!!!TEST CHAR DO NOT USE!!!",
}

var CharKeyToEle = map[Char]attributes.Element{
	AetherAnemo:      attributes.Anemo,
	LumineAnemo:      attributes.Anemo,
	AetherGeo:        attributes.Geo,
	LumineGeo:        attributes.Geo,
	AetherElectro:    attributes.Electro,
	LumineElectro:    attributes.Electro,
	AetherDendro:     attributes.Dendro,
	LumineDendro:     attributes.Dendro,
	AetherHydro:      attributes.Hydro,
	LumineHydro:      attributes.Hydro,
	AetherPyro:       attributes.Pyro,
	LuminePyro:       attributes.Pyro,
	AetherCryo:       attributes.Cryo,
	LumineCryo:       attributes.Cryo,
	Albedo:           attributes.Geo,
	Aloy:             attributes.Cryo,
	Amber:            attributes.Pyro,
	Barbara:          attributes.Hydro,
	Beidou:           attributes.Electro,
	Bennett:          attributes.Pyro,
	Chongyun:         attributes.Cryo,
	Cyno:             attributes.Electro,
	Diluc:            attributes.Pyro,
	Diona:            attributes.Cryo,
	Eula:             attributes.Cryo,
	Fischl:           attributes.Electro,
	Ganyu:            attributes.Cryo,
	Hutao:            attributes.Pyro,
	Jean:             attributes.Anemo,
	Kazuha:           attributes.Anemo,
	Kaeya:            attributes.Cryo,
	Ayaka:            attributes.Cryo,
	Ayato:            attributes.Hydro,
	Keqing:           attributes.Electro,
	Kirara:           attributes.Dendro,
	Klee:             attributes.Pyro,
	Sara:             attributes.Electro,
	Lisa:             attributes.Electro,
	Mona:             attributes.Hydro,
	Ningguang:        attributes.Geo,
	Noelle:           attributes.Geo,
	Qiqi:             attributes.Cryo,
	Raiden:           attributes.Electro,
	Razor:            attributes.Electro,
	Rosaria:          attributes.Cryo,
	Kokomi:           attributes.Hydro,
	Sayu:             attributes.Anemo,
	Sucrose:          attributes.Anemo,
	Tartaglia:        attributes.Hydro,
	Thoma:            attributes.Pyro,
	Venti:            attributes.Anemo,
	Xiangling:        attributes.Pyro,
	Xiao:             attributes.Anemo,
	Xingqiu:          attributes.Hydro,
	Xinyan:           attributes.Pyro,
	Yanfei:           attributes.Pyro,
	Yoimiya:          attributes.Pyro,
	Zhongli:          attributes.Geo,
	Gorou:            attributes.Geo,
	Itto:             attributes.Geo,
	Shenhe:           attributes.Cryo,
	Yunjin:           attributes.Geo,
	YaeMiko:          attributes.Electro,
	Yelan:            attributes.Hydro,
	Kuki:             attributes.Electro,
	Heizou:           attributes.Anemo,
	Tighnari:         attributes.Dendro,
	Collei:           attributes.Dendro,
	Dori:             attributes.Electro,
	Candace:          attributes.Hydro,
	Nilou:            attributes.Hydro,
	Nahida:           attributes.Dendro,
	Alhaitham:        attributes.Dendro,
	Layla:            attributes.Cryo,
	Faruzan:          attributes.Anemo,
	Wanderer:         attributes.Anemo,
	Baizhu:           attributes.Dendro,
	Dehya:            attributes.Pyro,
	Yaoyao:           attributes.Dendro,
	Mika:             attributes.Cryo,
	Kaveh:            attributes.Dendro,
	Freminet:         attributes.Cryo,
	TestCharDoNotUse: attributes.Geo,
}
