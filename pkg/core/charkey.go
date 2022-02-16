package core

type CharKey int

const (
	NoChar CharKey = iota
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

func (c CharKey) String() string {
	return charNames[c]
}

var CharNameToKey = map[string]CharKey{
	"travelerelectro":    TravelerElectro,
	"traveleranemo":      TravelerAnemo,
	"travelergeo":        TravelerGeo,
	"travelerhydro":      TravelerHydro,
	"travelercryo":       TravelerCryo,
	"travelerpyro":       TravelerPyro,
	"travelerdendro":     TravelerDendro,
	"traveler (electro)": TravelerElectro,
	"traveler (anemo)":   TravelerAnemo,
	"traveler (geo)":     TravelerGeo,
	"traveler (hydro)":   TravelerHydro,
	"traveler (cryo)":    TravelerCryo,
	"traveler (pyro)":    TravelerPyro,
	"traveler (dendro)":  TravelerDendro,
	"albedo":             Albedo,
	"aloy":               Aloy,
	"amber":              Amber,
	"barbara":            Barbara,
	"beidou":             Beidou,
	"bennett":            Bennett,
	"chongyun":           Chongyun,
	"diluc":              Diluc,
	"diona":              Diona,
	"eula":               Eula,
	"fischl":             Fischl,
	"ganyu":              Ganyu,
	"hutao":              Hutao,
	"jean":               Jean,
	"kaedeharakazuha":    Kazuha,
	"kazuha":             Kazuha,
	"kaeya":              Kaeya,
	"kamisatoayaka":      Ayaka,
	"ayaka":              Ayaka,
	"keqing":             Keqing,
	"klee":               Klee,
	"kujousara":          Sara,
	"kujosara":           Sara,
	"sara":               Sara,
	"lisa":               Lisa,
	"mona":               Mona,
	"ningguang":          Ningguang,
	"noelle":             Noelle,
	"qiqi":               Qiqi,
	"raidenshogun":       Raiden,
	"raiden":             Raiden,
	"razor":              Razor,
	"rosaria":            Rosaria,
	"sangonomiyakokomi":  Kokomi,
	"kokomi":             Kokomi,
	"sayu":               Sayu,
	"sucrose":            Sucrose,
	"tartaglia":          Tartaglia,
	"thoma":              Thoma,
	"venti":              Venti,
	"xiangling":          Xiangling,
	"xiao":               Xiao,
	"xingqiu":            Xingqiu,
	"xinyan":             Xinyan,
	"yanfei":             Yanfei,
	"yoimiya":            Yoimiya,
	"yunjin":             Yunjin,
	"zhongli":            Zhongli,
	"gorou":              Gorou,
	"aratakiitto":        Itto,
	"shenhe":             Shenhe,
	"yaemiko":            YaeMiko,
}

var charNames = []string{
	"",
	"traveler (electro)",
	"traveler (anemo)",
	"traveler (geo)",
	"traveler (hydro)",
	"traveler (cryo)",
	"traveler (pyro)",
	"traveler (dendro)",
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

var CharKeyToEle = map[CharKey]EleType{
	TravelerElectro: Electro,
	TravelerAnemo:   Anemo,
	TravelerGeo:     Geo,
	TravelerHydro:   Hydro,
	TravelerCryo:    Cryo,
	TravelerPyro:    Pyro,
	TravelerDendro:  Dendro,
	Albedo:          Geo,
	Aloy:            Cryo,
	Amber:           Pyro,
	Barbara:         Hydro,
	Beidou:          Electro,
	Bennett:         Pyro,
	Chongyun:        Cryo,
	Diluc:           Pyro,
	Diona:           Cryo,
	Eula:            Cryo,
	Fischl:          Electro,
	Ganyu:           Cryo,
	Hutao:           Pyro,
	Jean:            Anemo,
	Kazuha:          Anemo,
	Kaeya:           Cryo,
	Ayaka:           Cryo,
	Keqing:          Electro,
	Klee:            Pyro,
	Sara:            Electro,
	Lisa:            Electro,
	Mona:            Hydro,
	Ningguang:       Geo,
	Noelle:          Geo,
	Qiqi:            Cryo,
	Raiden:          Electro,
	Razor:           Electro,
	Rosaria:         Cryo,
	Kokomi:          Hydro,
	Sayu:            Anemo,
	Sucrose:         Anemo,
	Tartaglia:       Hydro,
	Thoma:           Pyro,
	Venti:           Anemo,
	Xiangling:       Pyro,
	Xiao:            Anemo,
	Xingqiu:         Hydro,
	Xinyan:          Pyro,
	Yanfei:          Pyro,
	Yoimiya:         Pyro,
	Zhongli:         Geo,
	Gorou:           Geo,
	Itto:            Geo,
	Shenhe:          Cryo,
	Yunjin:          Geo,
}
