package enemy

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

var levelMultiplier = [][]float64{
	{
		5.367859,
		6.818905,
		8.421231,
		10.17484,
		12.07987,
		14.13796,
		16.34828,
		17.50736,
		19.26782,
		21.09865,
		24.06383,
		27.14657,
		30.34872,
		33.90172,
		37.58709,
		41.4068,
		46.00103,
		50.05359,
		54.18154,
		65.1649,
		68.61109,
		72.10564,
		75.35814,
		79.78938,
		84.30028,
		88.9361,
		93.66512,
		98.48837,
		103.4069,
		108.4217,
		114.4719,
		120.6543,
		125.744,
		148.7045,
		153.7546,
		158.886,
		165.88,
		173.5809,
		181.4102,
		198.3841,
		207.9781,
		217.7367,
		227.0249,
		236.4653,
		264.514,
		273.7571,
		289.6475,
		305.8628,
		322.4075,
		367.8213,
		384.5919,
		401.5548,
		418.7124,
		436.0667,
		452.5111,
		466.241,
		483.8348,
		503.1843,
		522.7799,
		616.9946,
		637.3271,
		659.6903,
		682.2833,
		711.7287,
		734.9008,
		753.9569,
		829.3079,
		855.3966,
		879.7074,
		960.8343,
		987.2563,
		1016.308,
		1037.415,
		1067.75,
		1098.412,
		1123.775,
		1153.119,
		1182.76,
		1210.233,
		1366.735,
		1394.867,
		1423.247,
		1440.909,
		1485.468,
		1503.56,
		1532.91,
		1563.946,
		1598.809,
		1634.08,
		1792.851,
		1835.423,
		1882.428,
		1930.047,
		1996.661,
		2042.759,
		2055.588,
		2069.809,
		2256.937,
		2272.524,
		2706.502,
	},
	{
		5.404888,
		6.867565,
		8.483289,
		10.252208,
		12.174544,
		14.251988,
		16.483805,
		17.65636,
		19.436102,
		21.287472,
		24.350191,
		27.538502,
		30.854609,
		34.542477,
		38.37272,
		42.347908,
		47.130226,
		51.37316,
		57.933083,
		62.656433,
		67.4997,
		72.46275,
		77.45362,
		82.53305,
		87.70267,
		92.40437,
		97.388695,
		102.556755,
		109.56383,
		116.679184,
		123.90555,
		131.39304,
		139.20303,
		147.34232,
		156.63199,
		166.4199,
		176.60277,
		187.12581,
		198.11426,
		209.49992,
		221.35829,
		234.31367,
		248.21857,
		263.0379,
		281.4478,
		300.1646,
		319.93707,
		340.79654,
		362.84512,
		386.00797,
		405.81985,
		426.7509,
		448.33606,
		470.60852,
		500.87494,
		529.12994,
		558.18567,
		588.0836,
		619.623,
		652.1601,
		685.4282,
		725.286,
		766.45514,
		808.56494,
		857.3944,
		899.7023,
		943.45435,
		987.4856,
		1032.8478,
		1078.629,
		1125.0859,
		1176.6754,
		1225.4973,
		1276.6016,
		1349.6738,
		1419.612,
		1492.4435,
		1566.4338,
		1645.9409,
		1740.0337,
		1832.4424,
		1926.5065,
		2018.1666,
		2109.6982,
		2223.0679,
		2317.169,
		2412.8042,
		2510.078,
		2613.8433,
		2731.1755,
		2831.9614,
		2934.1345,
		3063.6953,
		3191.605,
		3344.277,
		3495.401,
		3644.6707,
		3791.2686,
		3934.9932,
		4082.7458,
	},
}

var monsterInfos = map[string]MonsterInfo{
	"hilichurl": {
		baseHP: 13.584, levelType: 1, particleDropCount: 1, particleThresholdCount: 2,
	},
	"mitachurl": {
		baseHP: 40.752, levelType: 1, particleDropCount: 1, resist: map[attributes.Element]float64{
			attributes.Physical: 0.3,
		},
	},
	"ruinguard": {
		baseHP:            95.088,
		levelType:         1,
		particleDropCount: 4,
		resist:            map[attributes.Element]float64{attributes.Physical: 0.7},
	},
	"ruingrader": {
		baseHP:            122.256,
		levelType:         1,
		particleDropCount: 4,
		resist:            map[attributes.Element]float64{attributes.Physical: 0.7},
	},
	"ruindrakeearthguard": {
		baseHP:                 95.088,
		levelType:              2,
		particleThresholdCount: 3,
		resist:                 map[attributes.Element]float64{attributes.Physical: 0.5},
	},
	"primogeovishap": {
		baseHP:                 407.52,
		levelType:              1,
		particleThresholdCount: 3,
		particleDropCount:      1, // BUGGED: primo geovishap drops geo particles
		resist:                 map[attributes.Element]float64{attributes.Physical: 0.3, attributes.Geo: 0.5},
	},
	"maguukenki": {baseHP: 271.68, levelType: 1},
	"ruinserpent": {
		baseHP:                 217.344,
		levelType:              2,
		hpMultiplier:           2,
		particleThresholdCount: 3,
		resist:                 map[attributes.Element]float64{attributes.Physical: 0.7},
	},
	"aeonblightdrake": {
		baseHP:                 230.928,
		levelType:              2,
		particleThresholdCount: 3,
		resist:                 map[attributes.Element]float64{attributes.Physical: 0.7},
	},
	"asimon": {baseHP: 217.344, levelType: 2, particleThresholdCount: 3},
}

type MonsterInfo struct {
	baseHP                 float64
	hpMultiplier           float64
	levelType              int
	particleDropCount      int
	particleThresholdCount int
	resist                 map[attributes.Element]float64
}

func ConfigureTarget(profile *EnemyProfile, name string, params map[string]int) error {
	info, ok := monsterInfos[name]
	if !ok {
		return fmt.Errorf("invalid target name `%v`", name)
	}
	if !(1 <= profile.Level && profile.Level <= 100) {
		return fmt.Errorf("invalid target level: must be between 1 and 100")
	}
	profile.HP = info.baseHP * levelMultiplier[info.levelType-1][profile.Level-1]
	if mult, ok := params["mult"]; ok {
		profile.HP *= float64(mult)
	} else {
		mult := 2.5
		if info.hpMultiplier > 0 {
			mult = info.hpMultiplier
		}
		profile.HP *= mult
	}
	if info.resist != nil {
		for k, v := range info.resist {
			profile.Resist[k] = v
		}
	}
	if part, ok := params["particles"]; !ok || part != 0 {
		thresholdCount := 4
		if info.particleThresholdCount > 0 {
			thresholdCount = info.particleThresholdCount
		}
		profile.ParticleDropThreshold = profile.HP / float64(thresholdCount)
		profile.ParticleDropCount = 3
		if info.particleDropCount > 0 {
			profile.ParticleDropCount = float64(info.particleDropCount)
		}
	}
	return nil
}
