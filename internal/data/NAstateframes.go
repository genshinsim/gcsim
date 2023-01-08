package data

import (
	"github.com/genshinsim/gcsim/internal/characters/ayato"
	"github.com/genshinsim/gcsim/internal/characters/cyno"
	"github.com/genshinsim/gcsim/internal/characters/tartaglia"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

var PercentDelay5 = make([]int, keys.EndCharKeys)
var PercentDelay5AltForms = make([]int, keys.EndCharKeys)
var AltFormStatusKeys = make([]string, keys.EndCharKeys)

const Unused = -1

func init() {
	for i := range PercentDelay5AltForms {
		PercentDelay5AltForms[i] = Unused
	}

	PercentDelay5[keys.Nahida] = 9
	PercentDelay5[keys.Xingqiu] = 7
	PercentDelay5[keys.Yelan] = 9
	PercentDelay5[keys.Raiden] = 13
	PercentDelay5[keys.Bennett] = 7
	PercentDelay5[keys.Diluc] = 15
	PercentDelay5[keys.Kazuha] = 10
	PercentDelay5[keys.Keqing] = 8
	PercentDelay5[keys.Xiangling] = 7
	PercentDelay5[keys.Albedo] = 9
	PercentDelay5[keys.Ayaka] = 7

	PercentDelay5[keys.Tartaglia] = 9
	PercentDelay5AltForms[keys.Tartaglia] = 12
	AltFormStatusKeys[keys.Tartaglia] = tartaglia.MeleeKey

	PercentDelay5[keys.Fischl] = 9
	PercentDelay5[keys.Ganyu] = 10
	PercentDelay5[keys.Jean] = 6

	PercentDelay5[keys.Lumine] = 7
	PercentDelay5[keys.LumineAnemo] = 7
	PercentDelay5[keys.LumineCryo] = 7
	PercentDelay5[keys.LumineDendro] = 7
	PercentDelay5[keys.LumineElectro] = 7
	PercentDelay5[keys.LumineGeo] = 7
	PercentDelay5[keys.LumineHydro] = 7
	PercentDelay5[keys.LuminePyro] = 7

	PercentDelay5[keys.Nilou] = 11
	// I didn't test Nilou E stance, assuming it's the same values for now

	PercentDelay5[keys.Venti] = 9
	PercentDelay5[keys.Zhongli] = 9
	PercentDelay5[keys.Amber] = 8
	PercentDelay5[keys.Collei] = 11
	PercentDelay5[keys.Diona] = 9
	PercentDelay5[keys.Faruzan] = 9
	PercentDelay5[keys.Gorou] = 11
	PercentDelay5[keys.Heizou] = 10
	PercentDelay5[keys.Kaeya] = 6
	PercentDelay5[keys.Kuki] = 15
	PercentDelay5[keys.Qiqi] = 7
	PercentDelay5[keys.Rosaria] = 10
	PercentDelay5[keys.Sara] = 14
	PercentDelay5[keys.Thoma] = 11
	PercentDelay5[keys.Yanfei] = 4
	PercentDelay5[keys.Yunjin] = 12

	PercentDelay5[keys.Beidou] = 22
	PercentDelay5[keys.Chongyun] = 18
	PercentDelay5[keys.Dori] = 29
	PercentDelay5[keys.Itto] = 27
	PercentDelay5[keys.Noelle] = 23
	PercentDelay5[keys.Razor] = 18
	PercentDelay5[keys.Sayu] = 24
	PercentDelay5[keys.Xinyan] = 28

	PercentDelay5[keys.Aether] = 8
	PercentDelay5[keys.AetherAnemo] = 8
	PercentDelay5[keys.AetherCryo] = 8
	PercentDelay5[keys.AetherDendro] = 8
	PercentDelay5[keys.AetherElectro] = 8
	PercentDelay5[keys.AetherGeo] = 8
	PercentDelay5[keys.AetherHydro] = 8
	PercentDelay5[keys.AetherPyro] = 8

	PercentDelay5[keys.Ayato] = 15
	PercentDelay5AltForms[keys.Ayato] = 17
	AltFormStatusKeys[keys.Ayato] = ayato.SkillBuffKey

	PercentDelay5[keys.Candace] = 14
	PercentDelay5[keys.Eula] = 22
	PercentDelay5[keys.Hutao] = 10
	PercentDelay5[keys.Yoimiya] = 17

	PercentDelay5[keys.Cyno] = 10
	PercentDelay5AltForms[keys.Cyno] = 12
	AltFormStatusKeys[keys.Cyno] = cyno.BurstKey

	PercentDelay5[keys.Layla] = 12
	PercentDelay5[keys.Shenhe] = 12
	PercentDelay5[keys.YaeMiko] = 7

	// TODO: Uncomment when Wanderer Implementation is done
	// PercentDelay5[keys.Wanderer] = 0
	// PercentDelay5AltForms[keys.Wanderer] = 12
	// AltFormStatusKeys[keys.Wanderer] = wanderer.SkillKey

	// Technically it's 15 for Left, 5 for Right, and 13 for Twirl
	PercentDelay5[keys.Ningguang] = (15 + 5 + 13) / 3

	// jumping/dashing during the NA windup for some catalysts modifies their frames - said by koli
	// thus the current method of NA -> jump to test for N0 timing may not work on them
	PercentDelay5[keys.Kokomi] = 0
	PercentDelay5[keys.Sucrose] = 0
	PercentDelay5[keys.Barbara] = 0
	PercentDelay5[keys.Lisa] = 0
	PercentDelay5[keys.Mona] = 0
	PercentDelay5[keys.Klee] = 0
}
