package keys

type Char int

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
)

func (c Char) String() string {
	return charNames[c]
}

var CharNameToKey = map[string]Char{
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
	"zhongli":            Zhongli,
	"gorou":              Gorou,
	"aratakiitto":        Itto,
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
}
